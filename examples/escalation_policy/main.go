package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetEscalationPolicies(&ilert.GetEscalationPoliciesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d escalation policies\n\n ", len(result.EscalationPolicies))
	for _, escalationPolicy := range result.EscalationPolicies {
		log.Printf("%+v\n", *escalationPolicy)
	}
}
