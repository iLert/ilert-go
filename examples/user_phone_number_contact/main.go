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
	user := result.User

	contactInput := ilert.CreateUserPhoneNumberContactInput{UserID: &user.ID, UserPhoneNumberContact: &ilert.UserPhoneNumberContact{Target: "your phone number", RegionCode: "your region code"}}
	contactResult, err := client.CreateUserPhoneNumberContact(&contactInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}
	log.Printf("User phone number contact successfully created!\n\n %+v\n", contactResult.UserPhoneNumberContact)
}
