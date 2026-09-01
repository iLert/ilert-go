package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetCurrentUser()
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}
	log.Printf("User:\n\n %+v\n", result.User)

	input := ilert.CreateUserInput{
		User: &ilert.User{
			FirstName: "your first name",
			LastName:  "your last name",
			Email:     "your email",
			Role:      ilert.UserRole.User,
		},
		SendNoInvitation: ilert.Bool(true),

		// set to true to buy a license for this user instead of checking the account
		// quota, this always buys a seat and charges the account for it
		PurchaseSeat: ilert.Bool(false),
	}
	createResult, err := client.CreateUser(&input)
	if err != nil {
		log.Println(createResult)
		log.Fatalln("ERROR:", err.Error())
	}
	log.Printf("User successfully created!\n\n %+v\n", createResult.User)
}
