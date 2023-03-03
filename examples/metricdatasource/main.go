package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createDataSourceInput := &ilert.CreateMetricDataSourceInput{
		MetricDataSource: &ilert.MetricDataSource{
			Name: "example",
			Type: ilert.MetricDataSourceType.Prometheus,
			Metadata: &ilert.MetricDataSourceMetadata{
				AuthType:  ilert.MetricDataSourceAuthType.Basic,
				BasicUser: "your prometheus username",
				BasicPass: "your prometheus password",
				Url:       "your prometheus url",
			},
		},
	}

	result, err := client.CreateMetricDataSource(createDataSourceInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Metric Data Source:\n\n %+v\n", result.MetricDataSource)
}
