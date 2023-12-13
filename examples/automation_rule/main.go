package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions/post
func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	serviceID := int(0)
	input := &ilert.GetAutomationRulesInput{Service: &serviceID}
	result, err := client.GetAutomationRules(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d automation rules\n\n ", len(result.AutomationRules))
	for _, automationRule := range result.AutomationRules {
		log.Printf("Incident template:\n\n %+v\n", *automationRule)
	}
}
