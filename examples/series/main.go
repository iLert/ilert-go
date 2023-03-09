package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createSingleSeriesInput := &ilert.CreateSingleSeriesInput{
		Series: &ilert.SingleSeries{
			Value: 500,
		},
		MetricKey: ilert.String("your metric integration key"),
	}

	err := client.CreateSingleSeries(createSingleSeriesInput)
	if err != nil {
		log.Fatalln("ERROR:", err)
	} else {
		log.Println("Series submitted successfully")
	}
}
