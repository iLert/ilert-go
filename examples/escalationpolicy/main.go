package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	// set your environment variables:
	// ILERT_ORGANIZATION="your organization"
	// ILERT_USERNAME="your username"
	// ILERT_PASSWORD="your password"
	client := ilert.NewClient()
	result, err := client.GetEscalationPolicies(&ilert.GetEscalationPoliciesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d escalation policies\n\n ", len(result.EscalationPolicies)))
	for _, escalationPolicy := range result.EscalationPolicies {
		log.Println(fmt.Sprintf("%+v\n", *escalationPolicy))
	}

}
