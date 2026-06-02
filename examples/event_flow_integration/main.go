package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	rf, err := client.GetEventFlows(&ilert.GetEventFlowsInput{})
	if err != nil {
		log.Println(rf)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d event flows\n", len(rf.EventFlows))

	if len(rf.EventFlows) == 0 {
		log.Fatalln("An event flow is required for this test")
	}
	eventFlowID := rf.EventFlows[0].ID

	rc, err := client.CreateEventFlowIntegration(&ilert.CreateEventFlowIntegrationInput{
		EventFlowIntegration: &ilert.EventFlowIntegration{
			IntegrationType: ilert.AlertSourceIntegrationTypes.API,
			EventFlowID:     eventFlowID,
		},
	})
	if err != nil {
		log.Println(rc)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("New event flow integration is created:\n%+v\n", *rc.EventFlowIntegration)

	rg, err := client.GetEventFlowIntegration(&ilert.GetEventFlowIntegrationInput{
		EventFlowIntegrationID: ilert.Int64(rc.EventFlowIntegration.ID),
	})
	if err != nil {
		log.Println(rg)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Fetched event flow integration:\n%+v\n", *rg.EventFlowIntegration)

	rl, err := client.GetEventFlowIntegrations(&ilert.GetEventFlowIntegrationsInput{
		EventFlowID: ilert.Int64(eventFlowID),
	})
	if err != nil {
		log.Println(rl)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d event flow integrations for event flow %d\n", len(rl.EventFlowIntegrations), eventFlowID)

	ru, err := client.UpdateEventFlowIntegration(&ilert.UpdateEventFlowIntegrationInput{
		EventFlowIntegrationID: ilert.Int64(rc.EventFlowIntegration.ID),
		EventFlowIntegration: &ilert.EventFlowIntegration{
			IntegrationType: ilert.AlertSourceIntegrationTypes.API,
			EventFlowID:     eventFlowID,
		},
	})
	if err != nil {
		log.Println(ru)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Updated event flow integration:\n%+v\n", *ru.EventFlowIntegration)

	_, err = client.DeleteEventFlowIntegration(&ilert.DeleteEventFlowIntegrationInput{
		EventFlowIntegrationID: ilert.Int64(rc.EventFlowIntegration.ID),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Deleted event flow integration %d\n", rc.EventFlowIntegration.ID)
}
