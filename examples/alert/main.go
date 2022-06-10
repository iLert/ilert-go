package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
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

}
