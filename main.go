package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var slackToken string
var concurrency = 1

func init() {
	flag.StringVar(&slackToken, "slack", "", "Slack Auth Token, See: https://api.slack.com/web#authentication")
	flag.IntVar(&concurrency, "threads", 1, "Number of simultaneous GETS against the bungie API")
}

func main() {

	var waitForResponses sync.WaitGroup

	flag.Parse()
	if slackToken == "" {
		fmt.Println("Please supply a slack authentication token")
		os.Exit(1)
	}

	var concurrencyLimiter = make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		concurrencyLimiter <- struct{}{}
	}

	users, err := getSlackUsers()
	if err != nil {
		log.Fatalf("Error retrieving user list from slack: %s", err.Error())
	}

	if len(users) == 0 {
		log.Fatalf("Apparently your slack instance has no users")
	}

	var responses = make(chan string, len(users))

	for _, user := range users {
		if user.Deleted == true {
			continue
		}
		if user.Bot == true {
			continue
		}
		waitForResponses.Add(1)
		<-concurrencyLimiter
		go func(username string) {
			id, err := resolveDestinyId(username)
			concurrencyLimiter <- struct{}{}
			if err != nil {
				log.Printf("Error resolving destiny ID: %s: %s", username, err.Error())
			} else {
				responses <- fmt.Sprintf("%s,%s", username, id)
			}
			waitForResponses.Done()
		}(user.Profile.FirstName)
	}
	waitForResponses.Wait()
	close(responses)
	for resp := range responses {
		fmt.Println(resp)
	}
}
