package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	result, err := client.GetAlerts(&ilert.GetAlertsInput{
		States: []*string{ilert.String(ilert.AlertStatuses.Accepted)},
	})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d pending alerts\n\n ", len(result.Alerts)))
	for _, alert := range result.Alerts {
		log.Println(fmt.Sprintf("%+v\n", *alert))
	}

	// id := int64(your alert id)
	// input := &ilert.GetAlertInput{Include: []*string{ilert.String("escalationRules")}, AlertID: &id}
	// result, err := client.GetAlert(input)
	// if err != nil {
	// 	log.Println(result)
	// 	log.Fatalln("ERROR:", err)
	// }
	// log.Println(fmt.Sprintf("%+v\n", *result.Alert))
}
