package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetAlertSources(&ilert.GetAlertSourcesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d alert sources\n\n ", len(result.AlertSources))
	for _, alertSource := range result.AlertSources {
		log.Printf("%+v\n", *alertSource)
	}

	rep, err := client.GetEscalationPolicies(&ilert.GetEscalationPoliciesInput{})
	if err != nil {
		log.Println(rep)
		log.Fatalln("ERROR:", err)
	}
	if len(rep.EscalationPolicies) == 0 {
		log.Fatalln("Escalation policy is required for this test")
	}

	ras, err := client.CreateAlertSource(&ilert.CreateAlertSourceInput{
		AlertSource: &ilert.AlertSource{
			Name:            "Test API Alert Source",
			IntegrationType: ilert.AlertSourceIntegrationTypes.API,
			SetupStatus:     ilert.AlertSourceSetupStatuses.Finished,
			EscalationPolicy: &ilert.EscalationPolicy{
				ID: rep.EscalationPolicies[0].ID,
			},
		},
	})
	if err != nil {
		log.Println(ras)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("New alert source is created with setup status %q:\n%+v\n", ras.AlertSource.SetupStatus, *ras.AlertSource)
}
