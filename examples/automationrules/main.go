package main

import (
	"fmt"
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
	log.Println(fmt.Sprintf("Found %d AutomationRules\n\n ", len(result.AutomationRules)))
	for _, automationRule := range result.AutomationRules {
		log.Println(fmt.Sprintf("IncidentTemplate:\n\n %+v\n", *automationRule))
	}

	// id := "your automation rule id"
	// result, err := client.DeleteAutomationRule(&ilert.DeleteAutomationRuleInput{AutomationRuleID: &id})
	// if err != nil {
	// 	log.Println(result)
	// 	log.Fatalln("ERROR:", err)
	// }
	// log.Println("AutomationRule deleted")
}
