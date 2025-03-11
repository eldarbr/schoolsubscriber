package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eldarbr/schoolauth"
	"github.com/eldarbr/schoolsubscriber/internal/domain"
)

const (
	slotsCheckPerion = 10 * time.Second
)

var (
	ErrFileFormatRanges = errors.New("ranges file has wrong format")
)

func main() {
	flgUsername := flag.String("u", "", "username")
	flgPassword := flag.String("p", "", "password")
	flgPathRanges := flag.String("r", "", "path to the file with time ranges")

	flag.Parse()

	if *flgUsername == "" || *flgPassword == "" || *flgPathRanges == "" {
		log.Println("Error - please provide username and password and path to a file with valid time ranges")

		return
	}

	timeRanges, err := parseTimeRanges(*flgPathRanges)
	if err != nil {
		log.Println("Err Parse time ranges:", err)

		return
	}

	managedToken := schoolauth.NewManagedToken(*flgUsername, *flgPassword, nil)

	dom, err := domain.NewDomain(context.Background(), managedToken, *flgUsername)
	if err != nil {
		log.Println("Err New domain:", err)

		return
	}

	goals, err := dom.GetCurrentGoals(context.Background())
	if err != nil {
		log.Println("Err Get current goals:", err)

		return
	}

	goals = domain.GoalsFilterEvaluated(goals)
	if len(goals) < 1 {
		return
	}

	goalsIDx := interactiveGoalDecision(goals)

	for {
		var (
			succ  bool
			start time.Time
		)

		start, succ, err = dom.AttemptSubscribe(context.Background(), goals[goalsIDx].GoalID, timeRanges, true)
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

func parseTimeRanges(path string) ([][2]time.Time, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer file.Close()

	var (
		reader = bufio.NewReader(file)
		ranges = make([][2]time.Time, 0)

		line  []byte
		start time.Time
		end   time.Time
	)

	for {
		line, err = reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("read line with buffered: %w", err)
		}

		spl := strings.Split(strings.Trim(string(line), " \n"), "_")
		if len(spl) != 2 && len(spl) != 3 {
			return nil, ErrFileFormatRanges
		}

		start, err = time.ParseInLocation(time.DateTime, spl[0], time.Local)
		if err != nil {
			return nil, fmt.Errorf("parse start time: %w", err)
		}

		end, err = time.ParseInLocation(time.DateTime, spl[1], time.Local)
		if err != nil {
			return nil, fmt.Errorf("parse end time: %w", err)
		}

		ranges = append(ranges, [2]time.Time{start, end})
	}

	return ranges, nil
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
