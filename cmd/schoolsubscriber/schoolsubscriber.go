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
	slotsCheckPerion = 10 * time.Second
	aliveProbePeriod = 15 * time.Minute
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

		aliveProbe time.Time
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

	goalsIDx := interactiveGoalDecision(goals)

	taskID, answerID, err := client.GetTaskIDAnswerID(context.Background(), goals[goalsIDx].GoalID)
	if err != nil {
		log.Println("Err Get task and answer ids: ", err)

		return
	}

	for {
		var (
			succ  bool
			start time.Time
		)

		if time.Since(aliveProbe) >= aliveProbePeriod {
			log.Println("alive")

			aliveProbe = time.Now()
		}

		start, succ, err = client.AttemptSubscribe(context.Background(), taskID, answerID, timeRanges, true)
		if err != nil {
			log.Println("Err Attempt:", err)

			continue
		}

		if succ {
			log.Printf("Subscribed for the slot: %s\n", start.Local().Format(time.DateTime))
		} else {
			time.Sleep(slotsCheckPerion)
		}
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

func interactiveGoalDecision(goals []domain.Goal) int {
	if len(goals) < 2 {
		log.Println("a goal has been chosen automatically:")

		for i := range goals {
			fmt.Printf("%7s - %25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
		}

		return 0
	}

	var (
		scanner = bufio.NewReader(os.Stdin)

		input  []byte
		err    error
		goalID string
	)

	for {
		fmt.Println("Choose a goal:")

		for i := range goals {
			fmt.Printf("%7s - %25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
		}

		input, err = scanner.ReadBytes('\n')
		if err != nil {
			log.Println("Err Read bytes stdin:", err)

			continue
		}

		goalID = strings.Trim(string(input), " \n")

		for i := range goals {
			if goals[i].GoalID == goalID {
				return i
			}
		}

		log.Println("Err given GoalID is not valid")
	}
}

func (t *ConfTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string

	err := unmarshal(&str)
	if err != nil {
		return fmt.Errorf("config.UnmarshalYAML unmarshal failed: %w", err)
	}

	vt, err := time.ParseInLocation(time.DateTime, str, time.Local)
	if err != nil {
		return fmt.Errorf("config.UnmarshalYAML url.Parse failed: %w", err)
	}

	*t = ConfTime(vt)

	return nil
}
