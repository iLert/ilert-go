package main

import (
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	var apiKey = "alert source API Key"
	event := &ilert.Event{
		APIKey:      apiKey,
		EventType:   ilert.EventTypes.Alert,
		Summary:     "My test incident summary",
		IncidentKey: "123456",
	}
	client := ilert.NewClient()
	resp, err := client.CreateEvent(event)
	if err != nil {
		log.Println(resp)
		log.Fatalln("ERROR:", err)
	}
	log.Println("Incident key:", resp.IncidentKey)
}
