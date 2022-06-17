package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	// set your environment variables:
	// ILERT_ORGANIZATION="your organization"
	// ILERT_USERNAME="your username"
	// ILERT_PASSWORD="your password"
	client := ilert.NewClient()
	result, err := client.GetAlertSources(&ilert.GetAlertSourcesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d alert sources\n\n ", len(result.AlertSources)))
	for _, alertSource := range result.AlertSources {
		log.Println(fmt.Sprintf("%+v\n", *alertSource))
	}

}
