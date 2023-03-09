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

	contact := &ilert.UserPhoneNumberContact{Target: "your phone number", RegionCode: "your region code"}
	contactInput := ilert.CreateUserPhoneNumberContactInput{UserID: &user.ID, UserPhoneNumberContact: contact}
	contactResult, err := client.CreateUserPhoneNumberContact(&contactInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}
	contact = contactResult.UserPhoneNumberContact

	preference := &ilert.UserUpdatePreference{Method: ilert.UserPreferenceMethod.Voice, Contact: &ilert.UserContactShort{ID: contact.ID}, Type: ilert.UserUpdatePreferenceType.AlertAccepted}
	preferenceInput := ilert.CreateUserUpdatePreferenceInput{UserID: &user.ID, UserUpdatePreference: preference}
	preferenceResult, err := client.CreateUserUpdatePreference(&preferenceInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}

	log.Printf("User update notification preference successfully created!\n\n %+v\n", preferenceResult.UserUpdatePreference)
}
