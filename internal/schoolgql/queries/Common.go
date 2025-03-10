package queries

import (
	"bufio"
	"errors"
	"strings"
	"time"
)

type TQuery string

type TOperationName string

type varsEmpty struct{}

var mapOpToQuery = map[TOperationName]TQuery{
	CalendarAddBookingToEventSlot:                calendarAddBookingToEventSlotQuery,
	CalendarGetNameLessStudentTimeslotsForReview: calendarGetNameLessStudentTimeslotsForReviewQuery,
	GetCredentialsByLogin:                        getCredentialsByLoginQuery,
	GetCourseInfoByStudent:                       getCourseInfoByStudentQuery,
	GetProjectAttemptEvaluationsInfoByStudent:    getProjectAttemptEvaluationsInfoByStudentQuery,
	GetProjectInfoByStudent:                      getProjectInfoByStudentQuery,
	GetStudentCurrentProjects:                    getStudentCurrentProjectsQuery,
	PublicProfileGetPersonalInfo:                 publicProfileGetPersonalInfoQuery,
	PublicProfileGetProjects:                     publicProfileGetProjectsQuery,
	SendInvitation:                               sendInvitationQuery,
	GetLocalCourseGoals:                          getLocalCourseGoalsQuery,
}

var (
	ErrNoQuery    = errors.New("no query is mapped to this operation name")
	ErrValidation = errors.New("query validation")
)

func GetQueryByOperationName(opname TOperationName) (*TQuery, error) {
	query, ok := mapOpToQuery[opname]
	if !ok {
		return nil, ErrNoQuery
	}

	firstLine, err := bufio.NewReader(strings.NewReader(string(query))).ReadBytes('\n')
	if err != nil || !strings.Contains(string(firstLine), string(opname)) {
		return nil, ErrValidation
	}

	return &query, nil
}

func ConvertTimeFormat(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
