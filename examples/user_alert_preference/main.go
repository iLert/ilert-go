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

	preference := &ilert.UserAlertPreference{Method: ilert.UserPreferenceMethod.Voice, Contact: &ilert.UserContactShort{ID: contact.ID}, DelayMin: 0, Type: ilert.UserAlertPreferenceType.HighPriority}
	preferenceInput := ilert.CreateUserAlertPreferenceInput{UserID: &user.ID, UserAlertPreference: preference}
	preferenceResult, err := client.CreateUserAlertPreference(&preferenceInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}

	log.Printf("User alert notification preference successfully created!\n\n %+v\n", preferenceResult.UserAlertPreference)
}
