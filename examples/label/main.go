package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	keysResult, err := client.GetServiceLabelKeys(&ilert.GetServiceLabelKeysInput{
		StartIndex: ilert.Int(0),
		MaxResults: ilert.Int(100),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d service label keys\n\n", len(keysResult.LabelKeys))

	for _, labelKey := range keysResult.LabelKeys {
		valuesResult, err := client.GetServiceLabelValues(&ilert.GetServiceLabelValuesInput{
			LabelKey: &labelKey.Key,
		})
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		for _, labelValue := range valuesResult.LabelValues {
			log.Printf("%s = %s\n", labelKey.Key, labelValue.Value)
		}
	}

	// alerts and telemetry sources have their own label namespaces. The alert endpoints
	// page in fixed steps of ilert.AlertLabelPageSize and read through a five minute cache
	// unless a query is set
	alertKeysResult, err := client.GetAlertLabelKeys(&ilert.GetAlertLabelKeysInput{})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d alert label keys\n", len(alertKeysResult.LabelKeys))

	for _, labelKey := range alertKeysResult.LabelKeys {
		alertValuesResult, err := client.GetAlertLabelValues(&ilert.GetAlertLabelValuesInput{
			LabelKey:   &labelKey.Key,
			MaxResults: ilert.Int(ilert.AlertLabelPageSize),
		})
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		log.Printf("alert %s has %d distinct values\n", labelKey.Key, len(alertValuesResult.LabelValues))
	}

	telemetryKeysResult, err := client.GetTelemetrySourceLabelKeys(&ilert.GetTelemetrySourceLabelKeysInput{})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d telemetry source label keys\n", len(telemetryKeysResult.LabelKeys))

	for _, labelKey := range telemetryKeysResult.LabelKeys {
		telemetryValuesResult, err := client.GetTelemetrySourceLabelValues(&ilert.GetTelemetrySourceLabelValuesInput{
			LabelKey: &labelKey.Key,
		})
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		log.Printf("telemetry source %s has %d distinct values\n", labelKey.Key, len(telemetryValuesResult.LabelValues))
	}
}
