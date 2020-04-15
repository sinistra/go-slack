package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

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

	api := slack.New(token)
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "a",
				Value: "no",
			},
		},
	}

	channelID, timestamp, err := api.PostMessage(channelId, slack.MsgOptionText("A new Message", false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	log.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}
