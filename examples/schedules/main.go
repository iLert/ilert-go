package main

import (
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	maxResultsProPage := 10

	// Fetch first page
	resultFirstPage0, err := client.GetSchedules(&ilert.GetSchedulesInput{
		StartIndex: ilert.Int(0),
		MaxResults: ilert.Int(maxResultsProPage),
	})
	if err != nil {
		log.Println(resultFirstPage0)
		log.Fatalln("ERROR:", err)
		return
	}
	log.Printf("Found %d schedules on first page\n\n ", len(resultFirstPage0.Schedules))
	for _, schedule := range resultFirstPage0.Schedules {
		log.Printf("%+v\n", *schedule)
	}

	if len(resultFirstPage0.Schedules) == 10 {
		// Fetch second page
		resultSecondPage, err := client.GetSchedules(&ilert.GetSchedulesInput{
			StartIndex: ilert.Int(1),
			MaxResults: ilert.Int(maxResultsProPage),
		})
		if err != nil {
			log.Println(resultSecondPage)
			log.Fatalln("ERROR:", err)
			return
		}
		log.Printf("Found %d schedules on second page\n\n ", len(resultSecondPage.Schedules))
		for _, schedule := range resultSecondPage.Schedules {
			log.Printf("%+v\n", *schedule)
		}
	}
}
