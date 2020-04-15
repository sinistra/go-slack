package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("could not load env file.")
	}

	token := os.Getenv("BOT_TOKEN")
	userId := os.Getenv("USER_ID")
	channelId := os.Getenv("CHANNEL_ID")
	log.Println("token", token, "userId", userId)

	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello
			fmt.Println("RTM incoming => Hello event")

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", channelId))

		case *slack.MessageEvent:
			fmt.Println("MessageEvent")
			spew.Dump(ev)
			fmt.Println(ev.Type, ev.SubType, ev.ClientMsgID, ev.User, ev.BotID, ev.Username, ev.Timestamp, ev.ThreadTimestamp)
			if ev.SubType == "" && ev.BotID == "" && len(ev.ThreadTimestamp) > 0 {
				log.Println("Do we have a reply to a message??")
				resp := rtm.NewOutgoingMessage(fmt.Sprintf("reply to %s", ev.Text), ev.Channel)

				// Respond in thread if not a direct message.
				if !strings.HasPrefix(ev.Channel, "D") {
					resp.ThreadTimestamp = ev.Timestamp
				}

				// Respond in same thread if message came from a thread.
				if ev.ThreadTimestamp != "" {
					resp.ThreadTimestamp = ev.ThreadTimestamp
				}

				rtm.SendMessage(resp)
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.DesktopNotificationEvent:
			fmt.Printf("Desktop Notification: %v\n", ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials\n")
			return

		case *slack.UserTypingEvent:
			fmt.Printf("User typing - don't care\n")

		default:
			// Ignore other events..
			spew.Dump(ev)
			fmt.Printf("Unexpected: %v\n", ev)
		}
	}
}
