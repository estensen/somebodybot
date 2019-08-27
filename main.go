package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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
				// info := rtm.GetInfo()

				text := ev.Text
				text = strings.TrimSpace(text)
				text = strings.ToLower(text)

				matched, _ := regexp.MatchString("someone", text)

				if matched {
					rtm.SendMessage(rtm.NewOutgoingMessage("Someone do it!", ev.Channel))
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
