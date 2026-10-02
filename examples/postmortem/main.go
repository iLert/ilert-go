package main

import (
	"log"

	"github.com/iLert/ilert-go/v4"
)

func main() {
	// the postmortem operations need an API token on API contract version 2 or later
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	incidentID := int64(0) // your specific incident id

	// an incident has at most one postmortem, Overwrite replaces the one it already has
	linked, err := client.LinkIncidentPostmortem(&ilert.LinkIncidentPostmortemInput{
		IncidentID: ilert.Int64(incidentID),
		Postmortem: &ilert.Postmortem{
			LinkURL:    "https://wiki.example.com/postmortems/checkout-outage",
			Visibility: ilert.PostmortemVisibility.Private,
		},
		Overwrite: ilert.Bool(true),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	postmortemID := linked.Postmortem.ID
	log.Printf("Linked postmortem %d (%s)\n", postmortemID, linked.Postmortem.Status)

	// the only update the API accepts is relinking, send the visibility to keep it
	updated, err := client.UpdateIncidentPostmortem(&ilert.UpdateIncidentPostmortemInput{
		IncidentID:   ilert.Int64(incidentID),
		PostmortemID: ilert.Int64(postmortemID),
		Postmortem: &ilert.Postmortem{
			Status:     ilert.PostmortemStatus.Linked,
			LinkURL:    "https://wiki.example.com/postmortems/checkout-outage-v2",
			Visibility: ilert.PostmortemVisibility.Public,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Postmortem now links to %s\n", updated.Postmortem.LinkURL)

	postmortems, err := client.GetIncidentPostmortems(&ilert.GetIncidentPostmortemsInput{IncidentID: ilert.Int64(incidentID)})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for _, postmortem := range postmortems.Postmortems {
		log.Printf("Postmortem %d: %s %s\n", postmortem.ID, postmortem.Status, postmortem.LinkURL)
	}

	single, err := client.GetIncidentPostmortem(&ilert.GetIncidentPostmortemInput{
		IncidentID:   ilert.Int64(incidentID),
		PostmortemID: ilert.Int64(postmortemID),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Postmortem by id:\n\n %+v\n", single.Postmortem)

	_, err = client.DeleteIncidentPostmortem(&ilert.DeleteIncidentPostmortemInput{
		IncidentID:   ilert.Int64(incidentID),
		PostmortemID: ilert.Int64(postmortemID),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	// generation runs asynchronously and needs the postmortem feature on the account,
	// only Slack channels are taken into account
	requested, err := client.RequestIncidentPostmortem(&ilert.RequestIncidentPostmortemInput{
		IncidentID: ilert.Int64(incidentID),
		Request: &ilert.PostmortemRequest{
			RootCause:  "A configuration change disabled the payment provider.",
			Visibility: ilert.PostmortemVisibility.Private,
			Language:   "en",
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Requested postmortem %d (%s)\n", requested.Postmortem.ID, requested.Postmortem.Status)
}
