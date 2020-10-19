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
	result, err := client.GetSchedules(&ilert.GetSchedulesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d schedules\n\n ", len(result.Schedules)))
	for _, schedule := range result.Schedules {
		log.Println(fmt.Sprintf("%+v\n", *schedule))
	}

}
