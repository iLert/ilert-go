package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	client := ilert.NewClient()
	result, err := client.GetNumbers(&ilert.GetNumbersInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d numbers\n\n ", len(result.Numbers)))
	for _, number := range result.Numbers {
		log.Println(fmt.Sprintf("%+v\n", *number))
	}
}
