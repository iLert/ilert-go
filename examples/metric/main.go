package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createMetricInput := &ilert.CreateMetricInput{
		Metric: &ilert.Metric{
			Name:            "example",
			AggregationType: ilert.MetricAggregationType.Average,
			DisplayType:     ilert.MetricDisplayType.Graph,
			Metadata: &ilert.MetricProviderMetadata{
				Query: "your prometheus query",
			},
			DataSource: &ilert.MetricDataSource{
				ID: 0, // your metric data source id
			},
		},
	}

	result, err := client.CreateMetric(createMetricInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Metric:\n\n %+v\n", result.Metric)
}
