package main

import (
	"encoding/json"
	"log"

	"github.com/iLert/ilert-go/v4"
)

func main() {
	var apiToken = "your API token"
	var alertSourceID = "your alert source id"
	// actions and credentials are early access APIs, so the client has to opt in to
	// both by their paths under /api/beta, and then sends the ilert-beta header the
	// API requires for them
	client := ilert.NewClient(
		ilert.WithAPIToken(apiToken),
		ilert.WithEarlyAccess("actions", "credentials"),
	)

	// a webhook URL is a credential of type NONE: the API keeps its path and query
	// as the secret
	credentialResult, err := client.CreateCredential(&ilert.CreateCredentialInput{
		Credential: &ilert.Credential{
			Name:        "discord ops channel",
			TypeHint:    ilert.CredentialType.None,
			HostURL:     "https://discord.com/api/webhooks/your-webhook-id/your-webhook-token",
			Category:    "global-action",
			SubCategory: "discord",
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Credential:\n\n %+v\n", credentialResult.Credential)

	createResult, err := client.CreateAction(&ilert.CreateActionInput{
		Action: &ilert.Action{
			Name:          "post high priority alerts to discord",
			TemplateKey:   "discord_alert",
			EntityKind:    ilert.ActionEntityKind.Alert,
			TriggerEvents: []string{ilert.ActionTriggerEvents.AlertCreated, ilert.ActionTriggerEvents.AlertNotResolved},
			// alert.not_resolved is only raised after a delay
			TriggerEventDelays: []ilert.ActionTriggerEventDelay{{Event: ilert.ActionTriggerEvents.AlertNotResolved, DelaySec: 600}},
			// the condition limits the action to one alert source, with the id as a string
			ConditionICL: `entity.source.id in ["` + alertSourceID + `"] and entity.priority == "HIGH"`,
			Config:       json.RawMessage(`{"pushComments":true}`),
			CredentialID: &credentialResult.Credential.ID,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Action:\n\n %+v\n", createResult.Action)

	// an update sends the whole action and clears what it leaves empty, so change
	// it as it was read
	action := createResult.Action
	action.Enabled = ilert.Bool(false)
	updateResult, err := client.UpdateAction(&ilert.UpdateActionInput{ActionID: &action.ID, Action: action})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Updated action:\n\n %+v\n", updateResult.Action)
}
