package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	var (
		apiToken string
		debug    bool
	)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("could not load env file.")
	}

	token := os.Getenv("OAUTH_TOKEN")
	// userId := os.Getenv("USER_ID")
	// log.Println("token", token, "userId", userId)

	flag.StringVar(&apiToken, "token", token, "Your Slack API Token")
	flag.BoolVar(&debug, "debug", false, "Show JSON output")
	flag.Parse()

	api := slack.New(apiToken, slack.OptionDebug(debug))

	// Get all stars for the usr.
	params := slack.NewStarsParameters()
	starredItems, _, err := api.GetStarred(params)
	if err != nil {
		fmt.Printf("Error getting stars: %s\n", err)
		return
	}
	for _, s := range starredItems {
		var desc string
		switch s.Type {
		case slack.TYPE_MESSAGE:
			desc = s.Message.Text
		case slack.TYPE_FILE:
			desc = s.File.Name
		case slack.TYPE_FILE_COMMENT:
			desc = s.File.Name + " - " + s.Comment.Comment
		case slack.TYPE_CHANNEL, slack.TYPE_IM, slack.TYPE_GROUP:
			desc = s.Channel
		}
		fmt.Printf("Starred %s: %s\n", s.Type, desc)
	}
}
