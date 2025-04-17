package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbGVydCIsImlsX2siOiJhZWYwOTJlNDE5NTI0Zjg0YjY3NmNiODg2OTM2YTRjZSIsImlsX3QiOiJBUEkiLCJpbF92IjoxfQ.Kxrazq0ZbtQJX1_j9NZX91PhbUDhx9GiAWV9rtJyRLk"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createHeartbeatMonitorInput := &ilert.CreateHeartbeatMonitorInput{
		HeartbeatMonitor: &ilert.HeartbeatMonitor{
			Name:        "example",
			IntervalSec: 60,
			// add an alert source for this heartbeat monitor to create alert
		},
	}

	result, err := client.CreateHeartbeatMonitor(createHeartbeatMonitorInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Heartbeat monitor:\n\n %+v\n", result.HeartbeatMonitor)
}
