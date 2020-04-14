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
	// channelId := os.Getenv("CHANNEL_ID")
	log.Println("token", token, "userId", userId)

	api := slack.New(token)
	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		fmt.Println(channel.Name)
		// channel is of type conversation & groupConversation
		// see all available methods in `conversation.go`
	}
}
