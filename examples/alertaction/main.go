package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	client := ilert.NewClient()

	// result, err := client.GetAlertActions(&ilert.GetAlertActionsInput{})
	// if err != nil {
	// 	log.Println(result)
	// 	log.Fatalln("ERROR:", err)
	// }
	// log.Println(fmt.Sprintf("Found %d connections\n\n ", len(result.AlertActions)))
	// for _, connection := range result.AlertActions {
	// 	s, _ := json.Marshal(connection)
	// 	log.Println(fmt.Sprintf("%+v\n", string(s)))
	// }

	// log.Fatalln("STOP")

	rep, err := client.GetEscalationPolicies(&ilert.GetEscalationPoliciesInput{})
	if err != nil {
		log.Println(rep)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d escalation policies\n\n ", len(rep.EscalationPolicies)))

	if len(rep.EscalationPolicies) == 0 {
		log.Fatalln("Escalation policy is required for this test")
	}

	ras, err := client.CreateAlertSource(&ilert.CreateAlertSourceInput{
		AlertSource: &ilert.AlertSource{
			Name:            "Test API Alert Source",
			IntegrationType: ilert.AlertSourceIntegrationTypes.API,
			EscalationPolicy: &ilert.EscalationPolicy{
				ID: rep.EscalationPolicies[0].ID,
			},
		},
	})
	if err != nil {
		log.Println(ras)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("New alert source is created:\n%+v\n", *ras.AlertSource))

	rcr, err := client.CreateConnector(&ilert.CreateConnectorInput{
		Connector: &ilert.Connector{
			Name: "Test GitHub Connector",
			Type: ilert.ConnectorTypes.Github,
			Params: &ilert.ConnectorParamsGithub{
				APIKey: "my api key",
			},
		},
	})
	if err != nil {
		log.Println(rcr)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("New connector is created:\n%+v\n", *rcr.Connector))

	rcn, err := client.CreateAlertAction(&ilert.CreateAlertActionInput{
		AlertAction: &ilert.AlertAction{
			Name:           "Test GitHub AlertAction",
			ConnectorType:  ilert.ConnectorTypes.Github,
			ConnectorID:    rcr.Connector.ID,
			TriggerMode:    ilert.AlertActionTriggerModes.Automatic,
			TriggerTypes:   ilert.AlertActionTriggerTypesAll,
			AlertSourceIDs: []int64{ras.AlertSource.ID},
			Params: ilert.AlertActionParamsGithub{
				Owner:      "my-org",
				Repository: "my-repo",
			},
		},
	})
	if err != nil {
		log.Println(rcn)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("New alert action is created:\n%+v\n", *rcn.AlertAction))
}
