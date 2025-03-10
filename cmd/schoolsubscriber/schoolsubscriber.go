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

	dom, err := domain.NewDomain(managedToken, *flgUsername)
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

	for i := range goals {
		fmt.Printf("%7s - %25s - %s\n", goals[i].GoalID, goals[i].Name, goals[i].Status)
	}

	succ := false

	for !succ {
		succ, err = dom.AttemptSubscribe(context.Background(), goals[0].GoalID, timeRanges, true)
		if err != nil {
			log.Println(err)

			return
		}

		time.Sleep(time.Second * 10)
	}
}

func parseTimeRanges(path string) ([][2]time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	ranges := make([][2]time.Time, 0)
	b := []byte(nil)
	start := time.Time{}
	end := time.Time{}

	for {
		b, err = reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("read line with buffered: %w", err)
		}

		spl := strings.Split(strings.Trim(string(b), " \n"), "_")
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
