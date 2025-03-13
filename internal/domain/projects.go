package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/eldarbr/schoolsubscriber/internal/schoolgql"
	"github.com/eldarbr/schoolsubscriber/internal/schoolgql/queries"
)

type Goal struct {
	GoalID string
	Name   string
	Status string
}

type Domain struct {
	userID    string
	studentID string
	tokener   Tokener
}

type Tokener interface {
	Get(ctx context.Context) (string, error)
}

func NewDomain(ctx context.Context, tokener Tokener, username string) (*Domain, error) {
	userID, studentID, err := GetUserIDStudentID(ctx, tokener, username)
	if err != nil {
		return nil, fmt.Errorf("get current user id: %w", err)
	}

	return &Domain{
		tokener:   tokener,
		userID:    userID,
		studentID: studentID,
	}, nil
}

func (dom *Domain) GetCurrentGoals(ctx context.Context) ([]Goal, error) {
	token, err := dom.tokener.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.GetStudentCurrentProjects)
	if err != nil {
		return nil, fmt.Errorf("new req - get projects: %w", err)
	}

	req.Variables = queries.VarsGetStudentCurrentProjects{UserID: dom.userID}
	respProjects := queries.ResponseGetStudentCurrentProjects{}

	err = req.MakeRequest(ctx, token, &respProjects)
	if err != nil {
		return nil, fmt.Errorf("make request - get projects: %w", err)
	}

	result := []Goal{}

	for _, project := range respProjects.Data.Student.GetStudentCurrentProjects {
		if project.GoalID != nil && project.GoalStatus != nil && *project.GoalStatus != ProjectStatusUnavailable {
			result = append(result, Goal{GoalID: *project.GoalID, Status: *project.GoalStatus})
		}

		if project.LocalCourseID != nil {
			var goals []Goal

			goals, err = dom.GetCourseCurrentGoals(ctx, *project.LocalCourseID)
			if err != nil {
				return nil, fmt.Errorf("get course goals: %w", err)
			}

			result = append(result, goals...)
		}
	}

	return result, nil
}

func GoalsFilterEvaluated(g []Goal) []Goal {
	res := make([]Goal, 0, len(g))

	for _, goal := range g {
		if goal.Status == ProjectStatusEvaluation {
			res = append(res, goal)
		}
	}

	return slices.Clone(res)
}

func (dom *Domain) GetCourseCurrentGoals(ctx context.Context, courseID int) ([]Goal, error) {
	token, err := dom.tokener.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.GetLocalCourseGoals)
	if err != nil {
		return nil, fmt.Errorf("new req - get projects: %w", err)
	}

	req.Variables = queries.VarsGetLocalCourseGoals{LocalCourseID: courseID}
	respProjects := queries.ResponseGetLocalCourseGoals{}

	err = req.MakeRequest(ctx, token, &respProjects)
	if err != nil {
		return nil, fmt.Errorf("make request - get projects: %w", err)
	}

	result := make([]Goal, 0, len(respProjects.Data.Course.GetLocalCourseGoals.LocalCourseGoals))

	for _, goal := range respProjects.Data.Course.GetLocalCourseGoals.LocalCourseGoals {
		if goal.Status == ProjectStatusUnavailable {
			continue
		}

		result = append(result, Goal{
			GoalID: goal.GoalID,
			Name:   goal.GoalName,
			Status: goal.Status,
		})
	}

	return result, nil
}

func (dom *Domain) GetTaskIDAnswerID(ctx context.Context, goalID string) (string, string, error) {
	taskID, err := GetTaskIDByGoalID(ctx, dom.tokener, goalID, dom.studentID)
	if err != nil {
		return "", "", fmt.Errorf("get task id: %w", err)
	}

	answerID, err := GetAnswerIDByGoalID(ctx, dom.tokener, goalID, dom.studentID)
	if err != nil {
		return "", "", fmt.Errorf("get task id: %w", err)
	}

	return taskID, answerID, nil
}

func (dom *Domain) AttemptSubscribe(ctx context.Context, taskID, answerID string, ranges [][2]time.Time, online bool,
) (time.Time, bool, error) {
	slots, err := GetSlotsRanges(ctx, dom.tokener, taskID, ranges)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("get slots from the ranges: %w", err)
	}

	if len(slots) == 0 {
		return time.Time{}, false, nil
	}

	log.Printf("Found %d slots\n", len(slots))

	for _, start := range slots {
		_, err = OccupySlot(ctx, dom.tokener, answerID, start, online)
		if err == nil {
			return start, true, nil
		}

		log.Println("Occupy:", err.Error())
	}

	return time.Time{}, false, nil
}

func GetSlotsRanges(ctx context.Context, tokener Tokener, taskID string, ranges [][2]time.Time) ([]time.Time, error) {
	numWorkers := len(ranges)

	rangesChan := make(chan [2]time.Time)
	slotsChan := make(chan time.Time)
	errChan := make(chan error, numWorkers)
	group := sync.WaitGroup{}

	for range numWorkers {
		group.Add(1)

		go func() {
			defer group.Done()

			for timeRange := range rangesChan {
				slots, err := GetSlots(ctx, tokener, taskID, timeRange[0], timeRange[1])
				if err != nil {
					errChan <- fmt.Errorf("get slots: %w", err)

					return
				}

				for _, slot := range slots {
					slotsChan <- slot
				}
			}
		}()
	}

	go func() {
		for _, r := range ranges {
			rangesChan <- r
		}

		close(rangesChan)
	}()

	slots := make([]time.Time, 0)
	collected := make(chan struct{})
	// collector
	go func() {
		for slot := range slotsChan {
			slots = append(slots, slot)
		}

		close(collected)
	}()

	group.Wait()
	close(slotsChan)
	close(errChan)

	<-collected

	errs := make([]error, 0, len(errChan))

	for err := range errChan {
		errs = append(errs, err)
	}

	err := errors.Join(errs...)
	if err != nil {
		return nil, fmt.Errorf("collect slots: %w", err)
	}

	slices.SortFunc(slots, func(a, b time.Time) int { return int(a.Unix() - b.Unix()) })

	return slots, nil
}

func GetUserIDStudentID(ctx context.Context, tokener Tokener, username string) (string, string, error) {
	token, err := tokener.Get(ctx)
	if err != nil {
		return "", "", fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.GetCredentialsByLogin)
	if err != nil {
		return "", "", fmt.Errorf("new req get credentials: %w", err)
	}

	req.Variables = queries.VarsGetCredentialsByLogin{Login: username}
	respCreds := queries.ResponseGetCredentialsByLogin{}

	err = req.MakeRequest(ctx, token, &respCreds)
	if err != nil {
		return "", "", fmt.Errorf("new req get credentials: %w", err)
	}

	return respCreds.Data.School21.GetStudentByLogin.UserID, respCreds.Data.School21.GetStudentByLogin.StudentID, nil
}

func GetAnswerIDByGoalID(ctx context.Context, tokener Tokener, goalID, studentID string) (string, error) {
	token, err := tokener.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.GetProjectAttemptEvaluationsInfoByStudent)
	if err != nil {
		return "", fmt.Errorf("new req get attempts: %w", err)
	}

	req.Variables = queries.VarsGetProjectAttemptEvaluationsInfoByStudent{GoalID: goalID, StudentID: studentID}
	resp := queries.ResponseGetProjectAttemptEvaluationsInfoByStudent{}

	err = req.MakeRequest(ctx, token, &resp)
	if err != nil {
		return "", fmt.Errorf("make req get attempts: %w", err)
	}

	answerID := ""

	for _, attempt := range resp.Data.School21.GetProjectAttemptEvaluationsInfo {
		if attempt.Auto.Status == AnswerStatusNotScheduled {
			answerID = attempt.StudentAnswerID
		}
	}

	if answerID == "" {
		return "", errors.New("no unscheduled answers found")
	}

	return answerID, nil
}

func GetTaskIDByGoalID(ctx context.Context, tokener Tokener, goalID, studentID string) (string, error) {
	token, err := tokener.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.GetProjectInfoByStudent)
	if err != nil {
		return "", fmt.Errorf("new req get project info: %w", err)
	}

	req.Variables = queries.VarsGetProjectInfoByStudent{GoalID: goalID, StudentID: studentID}
	resp := queries.ResponseGetProjectInfoByStudent{}

	err = req.MakeRequest(ctx, token, &resp)
	if err != nil {
		return "", fmt.Errorf("make req get project info: %w", err)
	}

	return resp.Data.School21.GetModuleByID.CurrentTask.TaskID, nil
}

func GetSlots(ctx context.Context, tokener Tokener, taskID string, from, to time.Time) ([]time.Time, error) {
	token, err := tokener.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.CalendarGetNameLessStudentTimeslotsForReview)
	if err != nil {
		return nil, fmt.Errorf("new req get timeslots: %w", err)
	}

	req.Variables = queries.VarsCalendarGetNameLessStudentTimeslotsForReview{
		TaskID: taskID,
		From:   schoolgql.FormatTimeToStr(from),
		To:     schoolgql.FormatTimeToStr(to),
	}
	resp := queries.ResponseCalendarGetNameLessStudentTimeslotsForReview{}

	err = req.MakeRequest(ctx, token, &resp)
	if err != nil {
		return nil, fmt.Errorf("make req get timeslots: %w", err)
	}

	var startTime time.Time

	result := make([]time.Time, 0, len(resp.Data.Student.GetNameLessStudentTimeslotsForReview.TimeSlots))

	for _, slotSpan := range resp.Data.Student.GetNameLessStudentTimeslotsForReview.TimeSlots {
		for i := range slotSpan.ValidStartTimes {
			startTime, err = schoolgql.FormatStrToTime(slotSpan.ValidStartTimes[i])
			if err != nil {
				return nil, fmt.Errorf("parse time: %w", err)
			}

			result = append(result, startTime)
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func OccupySlot(ctx context.Context, tokener Tokener, answerID string, slotStart time.Time, isOnline bool,
) (string, error) {
	token, err := tokener.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("tokener get token: %w", err)
	}

	req, err := schoolgql.NewRequest(queries.CalendarAddBookingToEventSlot)
	if err != nil {
		return "", fmt.Errorf("new req add booking: %w", err)
	}

	req.Variables = queries.VarsCalendarAddBookingToEventSlot{
		StartTime:          schoolgql.FormatTimeToStr(slotStart),
		AnswerID:           answerID,
		IsOnline:           isOnline,
		WasStaffSlotChosen: false,
	}
	resp := queries.ResponseCalendarAddBookingToEventSlot{}

	err = req.MakeRequest(ctx, token, &resp)
	if err != nil {
		return "", fmt.Errorf("make req add booking: %w", err)
	}

	return resp.Data.Student.AddBookingP2PToEventSlot.ID, nil
}
