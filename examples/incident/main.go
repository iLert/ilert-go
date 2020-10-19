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
	result, err := client.GetIncidents(&ilert.GetIncidentsInput{
		States: []*string{ilert.String(ilert.IncidentStatuses.Pending)},
	})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d pending incidents\n\n ", len(result.Incidents)))
	for _, incident := range result.Incidents {
		log.Println(fmt.Sprintf("%+v\n", *incident))
	}

}
