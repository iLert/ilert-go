package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	var org = "your organization"
	var username = "your username"
	var password = "your password"
	client := ilert.NewClient(ilert.WithBasicAuth(org, username, password))
	result, err := client.GetUptimeMonitors(&ilert.GetUptimeMonitorsInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d uptime monitors\n\n ", len(result.UptimeMonitors)))
	for _, uptimeMonitor := range result.UptimeMonitors {
		log.Println(fmt.Sprintf("%+v\n", *uptimeMonitor))
	}

}
