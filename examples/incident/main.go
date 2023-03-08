package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	incidentId := int64(0) //your specific incident id
	input := &ilert.GetIncidentInput{
		IncidentID: &incidentId,
		Include:    []*string{&ilert.IncidentInclude.AffectedTeams},
	}
	result, err := client.GetIncident(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Incident:\n\n %+v\n", result.Incident)
}
