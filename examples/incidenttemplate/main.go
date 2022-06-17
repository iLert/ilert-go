package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
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
	log.Println(fmt.Sprintf("Found %d IncidentTemplates\n\n ", len(result.IncidentTemplates)))
	for _, incidentTemplate := range result.IncidentTemplates {
		log.Println(fmt.Sprintf("IncidentTemplate:\n\n %+v\n", *incidentTemplate))
	}

	// id := int64(your incidenttemplate id)
	// input := &ilert.DeleteIncidentTemplateInput{IncidentTemplateID: &id}
	// result, err := client.DeleteIncidentTemplate(input)
	// if err != nil {
	// 	log.Println(result)
	// 	log.Fatalln("ERROR:", err)
	// }
	// log.Println("IncidentTemplates deleted")
}
