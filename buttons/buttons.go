package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("could not load env file.")
	}

	token := os.Getenv("OAUTH_TOKEN")
	userId := os.Getenv("USER_ID")
	channelId := os.Getenv("CHANNEL_ID")

	log.Println("token", token, "userId", userId)
	api := slack.New(token)

	attachment := slack.Attachment{
		Pretext:    "pretext",
		Fallback:   "We don't currently support your client",
		CallbackID: "accept_or_reject",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "accept",
				Text:  "Accept",
				Type:  "button",
				Value: "accept",
			},
			slack.AttachmentAction{
				Name:  "reject",
				Text:  "Reject",
				Type:  "button",
				Value: "reject",
				Style: "danger",
			},
		},
	}

	message := slack.MsgOptionAttachments(attachment)
	channelID, timestamp, err := api.PostMessage(channelId, slack.MsgOptionText("TheMessageText", false), message)
	if err != nil {
		fmt.Printf("Could not send message: %v\n", err)
	}
	fmt.Printf("Message with buttons sucessfully sent to channel %s at %s\n", channelID, timestamp)
	http.HandleFunc("/actions", actionHandler)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v\n", err)
	}
	fmt.Printf("Message button pressed by user %s with value %s\n", payload.User.Name, payload.Value)
}
