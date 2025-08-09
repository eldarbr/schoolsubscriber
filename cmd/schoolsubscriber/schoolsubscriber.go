package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eldarbr/go-auth/pkg/config"
	"github.com/eldarbr/schoolauth"
	"github.com/eldarbr/schoolsubscriber/internal/client/tgbot"
	"github.com/eldarbr/schoolsubscriber/internal/domain"
)

type ConfTime time.Time

type confTimeRanges struct {
	Start *ConfTime `yaml:"start"`
	End   *ConfTime `yaml:"end"`
}

type BotSetting struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
}

type appConf struct {
	TimeRanges []confTimeRanges `yaml:"ranges"`
	Bot        *BotSetting      `yaml:"bot"`
}

const (
	slotsCheckPeriod  = 10 * time.Second
	aliveProbePeriod  = 30 * time.Minute
	appDateTimeLocale = time.DateTime
)

var ErrFileFormatRanges = errors.New("ranges file has wrong format")

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

	var bot domain.Notificator

	if conf.Bot != nil {
		bot = tgbot.NewBot(conf.Bot.Token, conf.Bot.ChatID)
		err = bot.SendMessage(context.Background(), "Hi! Searching slots")
		if err != nil {
			log.Println("Err bot initialization message:", err)
			bot = nil
		}
	}

	managedToken := schoolauth.NewManagedToken(*flgUsername, *flgPassword, nil)

	client, err := domain.NewDomain(context.Background(), managedToken, *flgUsername, bot)
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

	log.Println("-", goal.GoalID, "alive")

	aliveTicker := time.NewTicker(aliveProbePeriod)
	defer aliveTicker.Stop()

	attemptTicker := time.NewTicker(slotsCheckPeriod)
	defer attemptTicker.Stop()

	var (
		succ   bool
		succCh = make(chan bool, 1)
		start  time.Time
	)

	succCh <- true // initial tick

	attempt := func() {
		start, succ, err = client.AttemptSubscribe(context.Background(), taskID, answerID, timeRanges, true)
		if err != nil {
			log.Println("-", goal.GoalID, "Err Attempt:", err)

			return
		}

		if succ {
			succCh <- succ // try again immediately
			log.Println("-", goal.GoalID, "Subscribed for the slot:", start.Local().Format(appDateTimeLocale))
		}
	}

	for { // loop
		select {
		case <-aliveTicker.C:
			log.Println("-", goal.GoalID, "alive")
		case <-succCh:
			attempt()
		case <-attemptTicker.C:
			attempt()
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

func interactiveGoalDecision(goals []domain.Goal) []domain.Goal {
	if len(goals) < 1 {
		return []domain.Goal{}
	}

	if len(goals) == 1 {
		log.Println("a goal has been chosen automatically:")

		for i := range goals {
			fmt.Printf("%7v - %-25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
		}

		return []domain.Goal{goals[0]}
	}

	scanner := bufio.NewReader(os.Stdin)

	availableGoalIDs := make(map[int]int, len(goals))
	for i, g := range goals {
		availableGoalIDs[g.GoalID] = i
	}

inpLoop:
	for {
		for i := range goals {
			fmt.Printf("%7v - %-25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
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
			inpGoalID, err := strconv.Atoi(inp)
			if err != nil {
				log.Println("Err given GoalID is not valid")

				continue inpLoop
			}

			if srcID, ok := availableGoalIDs[inpGoalID]; ok {
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
