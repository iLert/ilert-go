package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetCurrentAccount(&ilert.GetCurrentAccountInput{})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	account := result.Account
	log.Printf("Account %s (%s), timezone %s, region %s\n", account.OrganizationName, account.ID, account.Timezone, account.Region)
	log.Printf("AI mode: %s\n", account.AiMode)
	log.Printf("Unlocked features: %v\n", account.ApplicationFeatures)
	if account.Subscription != nil {
		log.Printf("Plan: %s (%s)\n", account.Subscription.Name, account.Subscription.Status)
	}
}
