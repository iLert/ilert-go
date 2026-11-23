package main

import (
	"log"

	"github.com/iLert/ilert-go/v4"
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

	rsh, err := client.CreateSupportHour(&ilert.CreateSupportHourInput{
		SupportHour: &ilert.SupportHour{
			Name:     "Business hours",
			Timezone: "Europe/Berlin",
			SupportDays: &ilert.SupportDays{
				MONDAY: &ilert.SupportDay{Start: "08:00", End: "17:00"},
				FRIDAY: &ilert.SupportDay{Start: "08:00", End: "17:00"},
			},
		},
	})
	if err != nil {
		log.Println(rsh)
		log.Fatalln("ERROR:", err)
	}

	// with support hours referenced by their id, every write sets AutoRaiseAlerts:
	// left out, it turns raising the pending alerts off, so send the value read to keep it
	rws, err := client.CreateAlertSource(&ilert.CreateAlertSourceInput{
		AlertSource: &ilert.AlertSource{
			Name:              "Test API Alert Source with support hours",
			IntegrationType:   ilert.AlertSourceIntegrationTypes.API,
			EscalationPolicy:  &ilert.EscalationPolicy{ID: rep.EscalationPolicies[0].ID},
			AlertPriorityRule: "HIGH_DURING_SUPPORT_HOURS",
			SupportHours:      &ilert.SupportHour{ID: rsh.SupportHour.ID},
			AutoRaiseAlerts:   ilert.Bool(true),
		},
	})
	if err != nil {
		log.Println(rws)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("New alert source is created with auto raise alerts %t:\n%+v\n", *rws.AlertSource.AutoRaiseAlerts, *rws.AlertSource)
}
