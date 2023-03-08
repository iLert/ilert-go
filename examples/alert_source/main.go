package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetAlertSources(&ilert.GetAlertSourcesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d alert sources\n\n ", len(result.AlertSources))
	for _, alertSource := range result.AlertSources {
		log.Printf("%+v\n", *alertSource)
	}
}
