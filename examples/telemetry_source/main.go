package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	labels := map[string]string{"env": "production"}
	createResult, err := client.CreateTelemetrySource(&ilert.CreateTelemetrySourceInput{
		TelemetrySource: &ilert.TelemetrySource{
			Name:              "OpenTelemetry collector",
			Type:              ilert.TelemetrySourceType.Otel,
			ServiceNamePrefix: "prod-",
			Labels:            &labels,
		},
		// teams are only part of the response when they are requested
		Include: []*string{&ilert.TelemetrySourceInclude.Teams},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Telemetry source:\n\n %+v\n", createResult.TelemetrySource)

	// the integration key is what the collector authenticates with, it is omitted from
	// list responses and for users without update permission on the telemetry source
	log.Printf("Integration key: %s\n", createResult.TelemetrySource.IntegrationKey)

	statsResult, err := client.GetTelemetrySourceStats(&ilert.GetTelemetrySourceStatsInput{
		TelemetrySourceID: &createResult.TelemetrySource.ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Last trace received at: %s\n", statsResult.TelemetrySourceStats.LastTraceAt)
}
