package main

import (
	"github.com/davecgh/go-spew/spew"
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
	log.Println("token", token, "userId", userId)

	api := slack.New(token)
	user, err := api.GetUserInfo(userId)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	spew.Dump(user)
	log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}
