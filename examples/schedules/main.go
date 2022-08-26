package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
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
