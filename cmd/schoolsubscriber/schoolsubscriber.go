package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/eldarbr/schoolauth"
	"github.com/eldarbr/schoolsubscriber/internal/domain"
)

func main() {
	flgUsername := flag.String("u", "", "username")
	flgPassword := flag.String("p", "", "password")

	flag.Parse()

	if *flgUsername == "" || *flgPassword == "" {
		log.Println("Error - please provide username and password and target username")

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
		succ, err = dom.AttemptSubscribe(context.Background(), goals[0].GoalID, time.Now(), time.Now().Add(time.Hour*3), true)
		if err != nil {
			log.Println(err)

			return
		}

		time.Sleep(time.Second * 10)
	}
}
