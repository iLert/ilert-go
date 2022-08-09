package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJpbGVydCIsImlsX3QiOiJBUEkiLCJpbF92IjoxLCJpbF9rIjoiZDAxYmY5Njg4ZjRlNGMzYmIxNDY5YWNhNjllNTgxY2IifQ.eDSLCoArmP_Qbh60jly9bza1gmVYx1ROdLsvTMF0bvk"
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
