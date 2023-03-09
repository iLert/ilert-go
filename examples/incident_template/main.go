package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	input := &ilert.GetIncidentTemplatesInput{}
	result, err := client.GetIncidentTemplates(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d incident templates\n\n ", len(result.IncidentTemplates))
	for _, incidentTemplate := range result.IncidentTemplates {
		log.Printf("Incident template:\n\n %+v\n", *incidentTemplate)
	}
}
