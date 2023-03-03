package main

import (
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	result, err := client.GetCurrentUser()
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}
	user := result.User

	contactInput := ilert.CreateUserEmailContactInput{UserID: &user.ID, UserEmailContact: &ilert.UserEmailContact{Target: "your email"}}
	contactResult, err := client.CreateUserEmailContact(&contactInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}
	log.Printf("User email contact successfully created!\n\n %+v\n", contactResult.UserEmailContact)
}
