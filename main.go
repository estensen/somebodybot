package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	token := os.Getenv("SLACKTOKEN")
	if token == "" {
		log.Fatal("missing Slack token env var")
	}

	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received\n")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				text := ev.Text
				text = strings.TrimSpace(text)
				text = strings.ToLower(text)

				someone, _ := regexp.MatchString("someone", text)
				somebody, _ := regexp.MatchString("somebody", text)

				if someone || somebody {
					users, _ := api.GetUsers()
					var userIDs []string
					for _, user := range users {
						// Filter Apps
						if user.Profile.DisplayName != ""  &&  user.Profile.DisplayName != "Slackbot" {
							userIDs = append(userIDs, user.ID)
						}
					}

					rand.Seed(time.Now().Unix())
					randomUserID := userIDs[rand.Intn(len(userIDs))]
					message := fmt.Sprintf("<@%s> do it!", randomUserID)
					rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
				}

			case *slack.RTMError:
				fmt.Printf("error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("invalid credentials")
				break Loop

			default:
				// do nothing
			}
		}
	}
}
