package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	// set your environment variables:
	// ILERT_ORGANIZATION="your organization"
	// ILERT_USERNAME="your username"
	// ILERT_PASSWORD="your password"
	client := ilert.NewClient()
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
