package main

import (
	"log"
	"os"
	"time"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiKey = "alert source API Key"
	event := &ilert.Event{
		APIKey:    apiKey,
		EventType: ilert.EventTypes.Alert,
		Summary:   "My test alert summary",
		AlertKey:  "123456",
	}
	input := &ilert.CreateEventInput{Event: event}
	client := ilert.NewClient(ilert.WithRetry(10, 5*time.Second, 20*time.Second))
	result, err := client.CreateEvent(input)
	if err != nil {
		if apiErr, ok := err.(*ilert.GenericAPIError); ok {
			if apiErr.Code == "NO_OPEN_INCIDENT_WITH_KEY" {
				log.Println("WARN:", apiErr.Error())
				os.Exit(0)
			} else {
				log.Fatalln("ERROR:", apiErr.Error())
			}
		} else {
			log.Println(result)
			log.Fatalln("ERROR:", err)
		}
	}

	log.Println("Incident key:", result.EventResponse.AlertKey)
}
