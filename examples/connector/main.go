package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetConnectors(&ilert.GetConnectorsInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d connectors\n\n ", len(result.Connectors))
	for _, connector := range result.Connectors {
		log.Printf("%+v\n", *connector)
	}
}
