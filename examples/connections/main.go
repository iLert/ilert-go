package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	client := ilert.NewClient()
	result, err := client.GetConnections(&ilert.GetConnectionsInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d connections\n\n ", len(result.Connections)))
	for _, connection := range result.Connections {
		log.Println(fmt.Sprintf("%+v\n", *connection))
	}
}
