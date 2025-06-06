package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
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

	searchResult, err := client.SearchHeartbeatMonitor(&ilert.SearchHeartbeatMonitorInput{
		HeartbeatMonitorName: &result.HeartbeatMonitor.Name,
	})
	if err != nil {
		log.Println(searchResult)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Search heartbeat monitors:\n\n %+v\n", searchResult.HeartbeatMonitor)
}
