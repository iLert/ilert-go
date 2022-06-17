package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	client := ilert.NewClient()
	result, err := client.GetConnectors(&ilert.GetConnectorsInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d connectors\n\n ", len(result.Connectors)))
	for _, connector := range result.Connectors {
		log.Println(fmt.Sprintf("%+v\n", *connector))
	}
}
