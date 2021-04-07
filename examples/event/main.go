package main

import (
	"log"
	"time"

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
	input := &ilert.CreateEventInput{Event: event}
	client := ilert.NewClient(ilert.WithRetry(10, 5*time.Second, 20*time.Second))
	result, err := client.CreateEvent(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println("Incident key:", result.EventResponse.IncidentKey)
}
