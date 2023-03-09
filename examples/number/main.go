package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetNumbers(&ilert.GetNumbersInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d numbers\n\n ", len(result.Numbers))
	for _, number := range result.Numbers {
		log.Printf("%+v\n", *number)
	}
}
