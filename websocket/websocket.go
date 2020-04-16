package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"

	"github.com/slack-go/slack"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("could not load env file.")
	}

	token := os.Getenv("BOT_TOKEN")
	userID := os.Getenv("USER_ID")
	channelID := os.Getenv("CHANNEL_ID")
	log.Println("token", token, "userId", userID)

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
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world, we're connected", channelID))

		case *slack.MessageEvent:
			fmt.Println("MessageEvent")
			spew.Dump(ev)
			fmt.Println(ev.Type, ev.SubType, ev.ClientMsgID, ev.User, ev.BotID, ev.Username, ev.Timestamp, ev.ThreadTimestamp)
			if ev.SubType == "" && ev.BotID == "" && len(ev.ThreadTimestamp) > 0 {
				log.Println("Do we have a reply to a message??")
				msgText := fmt.Sprintf("bot update to %s", ev.Text)
				// resp := rtm.NewOutgoingMessage(msgText, ev.Channel)
				msgResponse := slack.MsgOptionText(msgText, false)
				msgThumbNail := slack.MsgOptionIconURL("nms.com.au/images/thumbs-up20.jpg")
				alert := slack.Attachment{
					Fallback: "fallback text",
					Color:    "#ff0000",
					Title:    "attachment title",
					Text:     "attachment text",
					// ImageURL: "http://nms.com.au/favicon.png",
					ThumbURL: "http://nms.com.au/images/thumbs-up20.jpg",
					MarkdownIn: []string{
						"text",
						"pretext",
						"fields",
					},
					Footer: "attachment footer",
				}
				msgAttachment := slack.MsgOptionAttachments(alert)

				// Respond in thread if not a direct message.
				// if !strings.HasPrefix(ev.Channel, "D") {
				// resp.ThreadTimestamp = ev.Timestamp
				// }

				// Respond in same thread if message came from a thread.
				// if ev.ThreadTimestamp != "" {
				// resp.ThreadTimestamp = ev.ThreadTimestamp
				// }

				// rtm.PostMessage(ev.Channel, )

				rtm.UpdateMessage(ev.Channel, ev.ThreadTimestamp, msgResponse, msgThumbNail, msgAttachment)
				// rtm.SendMessage(resp)
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

		case *slack.AckMessage:
			fmt.Printf("AckMessage: Reply to %d, Timestamp:%s, Text:%s\n", ev.ReplyTo, ev.Timestamp, ev.Text)

		default:
			// Ignore other events..
			fmt.Println("Unexpected event:")
			spew.Dump(ev)
		}
	}
}
