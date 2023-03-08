package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

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
