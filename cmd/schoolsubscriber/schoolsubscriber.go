package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eldarbr/go-auth/pkg/config"
	"github.com/eldarbr/schoolauth"
	"github.com/eldarbr/schoolsubscriber/internal/domain"
)

type ConfTime time.Time

type confTimeRanges struct {
	Start *ConfTime `yaml:"start"`
	End   *ConfTime `yaml:"end"`
}

type appConf struct {
	TimeRanges []confTimeRanges `yaml:"ranges"`
}

const (
	slotsCheckPerion  = 10 * time.Second
	aliveProbePeriod  = 30 * time.Minute
	appDateTimeLocale = time.DateTime
)

var (
	ErrFileFormatRanges = errors.New("ranges file has wrong format")
)

func main() {
	var (
		conf appConf

		flgUsername = flag.String("u", "", "username")
		flgPassword = flag.String("p", "", "password")
		flgConf     = flag.String("c", "", "path to the config")
	)

	flag.Parse()

	if *flgUsername == "" || *flgPassword == "" || *flgConf == "" {
		log.Println("Error - please provide username and password and path to a yaml with valid time ranges")

		return
	}

	err := config.ParseConfig(*flgConf, &conf)
	if err != nil {
		log.Println("Err Reading config:", err)

		return
	}

	timeRanges := convConfTimeRanges(conf.TimeRanges)
	if len(timeRanges) < 1 {
		log.Println("Err Parse time ranges: no ranges")

		return
	}

	managedToken := schoolauth.NewManagedToken(*flgUsername, *flgPassword, nil)

	client, err := domain.NewDomain(context.Background(), managedToken, *flgUsername)
	if err != nil {
		log.Println("Err New domain:", err)

		return
	}

	goals, err := client.GetCurrentGoals(context.Background())
	if err != nil {
		log.Println("Err Get current goals:", err)

		return
	}

	goals = domain.GoalsFilterEvaluated(goals)
	if len(goals) < 1 {
		log.Println("No goals to review :)")

		return
	}

	chosenGoals := interactiveGoalDecision(goals)

	PrintRanges(timeRanges)

	group := sync.WaitGroup{}
	for _, goal := range chosenGoals {
		group.Add(1)

		go attemptWorker(client, timeRanges, goal, &group)
	}

	group.Wait()
}

func attemptWorker(client *domain.Domain, timeRanges [][2]time.Time, goal domain.Goal, group *sync.WaitGroup) {
	if group != nil {
		defer group.Done()
	}

	taskID, answerID, err := client.GetTaskIDAnswerID(context.Background(), goal.GoalID)
	if err != nil {
		log.Println("-", goal.GoalID, "Err Get task and answer ids: ", err)

		return
	}

	aliveProbe := time.Time{}

	ticker := time.NewTicker(slotsCheckPerion)
	defer ticker.Stop()

	for {
		var (
			succ  bool
			start time.Time
		)

		if time.Since(aliveProbe) >= aliveProbePeriod {
			log.Println("-", goal.GoalID, "alive")

			aliveProbe = time.Now()
		}

		start, succ, err = client.AttemptSubscribe(context.Background(), taskID, answerID, timeRanges, true)
		if err != nil {
			log.Println("-", goal.GoalID, "Err Attempt:", err)

			continue
		}

		if succ {
			log.Println("-", goal.GoalID, "Subscribed for the slot:", start.Local().Format(appDateTimeLocale))

			continue
		}

		<-ticker.C
	}
}

func convConfTimeRanges(ranges []confTimeRanges) [][2]time.Time {
	result := make([][2]time.Time, 0, len(ranges))

	for _, r := range ranges {
		if r.Start == nil || r.End == nil {
			continue
		}

		result = append(result, [2]time.Time{time.Time(*r.Start), time.Time(*r.End)})
	}

	return result
}

func interactiveGoalDecision(goals []domain.Goal) []domain.Goal {
	if len(goals) < 1 {
		return []domain.Goal{}
	}

	if len(goals) < 2 {
		log.Println("a goal has been chosen automatically:")

		for i := range goals {
			fmt.Printf("%7s - %-25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
		}

		return []domain.Goal{goals[0]}
	}

	scanner := bufio.NewReader(os.Stdin)

	availableGoalIDs := make(map[string]int, len(goals))
	for i, g := range goals {
		availableGoalIDs[g.GoalID] = i
	}

inpLoop:
	for {
		for i := range goals {
			fmt.Printf("%7s - %-25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
		}

		fmt.Print("Choose goals, comma separated or \"all\": ")

		input, err := scanner.ReadBytes('\n')
		if err != nil {
			log.Println("Err Read bytes stdin:", err)

			continue
		}

		if input[len(input)-1] == '\n' {
			input = input[:len(input)-1]
		}

		if string(input) == "all" {
			return goals
		}

		goalIDs := strings.Split(strings.ReplaceAll(string(input), " ", ""), ",")
		result := make([]domain.Goal, 0, len(goalIDs))

		if len(goalIDs) < 1 {
			continue
		}

		for _, inp := range goalIDs {
			if srcID, ok := availableGoalIDs[inp]; ok {
				result = append(result, goals[srcID])
			} else {
				log.Println("Err given GoalID is not valid")

				continue inpLoop
			}
		}

		return result
	}
}

func (t *ConfTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string

	err := unmarshal(&str)
	if err != nil {
		return fmt.Errorf("config.UnmarshalYAML unmarshal failed: %w", err)
	}

	vt, err := time.ParseInLocation(appDateTimeLocale, str, time.Local)
	if err != nil {
		return fmt.Errorf("config.UnmarshalYAML url.Parse failed: %w", err)
	}

	*t = ConfTime(vt)

	return nil
}

func PrintRanges(ranges [][2]time.Time) {
	fmt.Println("Working with this set of time ranges:")

	for _, r := range ranges {
		fmt.Println(
			" - from:",
			r[0].Format(appDateTimeLocale),
			"\tto:",
			r[1].Format(appDateTimeLocale),
		)
	}
}
