package main

import (
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("could not load env file.")
	}

	token := os.Getenv("OAUTH_TOKEN")
	userId := os.Getenv("USER_ID")
	log.Println("token", token, "userId", userId)
	api := slack.New(token)
	//Example for single user
	billingActive, err := api.GetBillableInfo(userId)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	log.Printf("ID: %s, BillingActive: %v\n\n\n", userId, billingActive[userId])

	//Example for team
	billingActiveForTeam, _ := api.GetBillableInfoForTeam()
	// spew.Dump(billingActiveForTeam)
	for id, value := range billingActiveForTeam {
		log.Printf("ID: %v, BillingActive: %v\n", id, value)
		user, err := api.GetUserInfo(id)
		if err != nil {
			log.Printf("%s\n", err)
			return
		}
		log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
	}
}
