package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	// var apiToken = "your API token"
	// client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	// input := &ilert.GetIncidentsInput{}
	// result, err := client.GetIncidents(input)
	// if err != nil {
	// 	log.Println(result)
	// 	log.Fatalln("ERROR:", err)
	// }
	// log.Println(fmt.Sprintf("Found %d Incidents\n\n ", len(result.Incidents)))
	// for _, incident := range result.Incidents {
	// 	log.Println(fmt.Sprintf("Incidents:\n\n %+v\n", *incident))
	// }

	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	incidentId := int64(44504) //your specific incident id
	input := &ilert.GetIncidentInput{
		IncidentID: &incidentId,
		Include:    []*string{&ilert.Include.AffectedTeams},
	}
	result, err := client.GetIncident(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Incident:\n\n %+v\n", result.Incident))
}
