package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	var org = "your organization"
	var username = "your username"
	var password = "your password"
	client := ilert.NewClient(ilert.WithBasicAuth(org, username, password))
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
