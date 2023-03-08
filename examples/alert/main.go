package main

import (
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
	log.Printf("Found %d pending alerts\n\n ", len(result.Alerts))
	for _, alert := range result.Alerts {
		log.Printf("%+v\n", *alert)
	}
}
