package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetUptimeMonitors(&ilert.GetUptimeMonitorsInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d uptime monitors\n\n ", len(result.UptimeMonitors))
	for _, uptimeMonitor := range result.UptimeMonitors {
		log.Printf("%+v\n", *uptimeMonitor)
	}
}
