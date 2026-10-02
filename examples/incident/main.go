package main

import (
	"log"

	"github.com/iLert/ilert-go/v4"
)

func main() {
	// the incident operations need an API token on API contract version 2 or later
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	serviceID := int64(0) // your specific service id
	incident := &ilert.Incident{
		Title:    "checkout unavailable",
		Summary:  "Payments fail for all customers.",
		Severity: 1,
		AffectedServices: []ilert.AffectedServices{
			{Impact: ilert.ServiceStatus.MajorOutage, Service: ilert.Service{ID: serviceID}},
		},
	}

	// forecasts who would be notified, only the affected services are taken into account
	affected, err := client.GetIncidentAffected(&ilert.GetIncidentAffectedInput{Incident: incident})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Would notify %d public and %d private subscribers\n", affected.Affected.PublicSubscribers, affected.Affected.PrivateSubscribers)

	created, err := client.CreateIncident(&ilert.CreateIncidentInput{Incident: incident})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	incidentID := created.Incident.ID
	log.Printf("Declared incident #%d\n", created.Incident.Number)

	open, err := client.GetIncidents(&ilert.GetIncidentsInput{
		States: []*string{
			ilert.String(ilert.IncidentStatus.Declared),
			ilert.String(ilert.IncidentStatus.Investigating),
			ilert.String(ilert.IncidentStatus.Identified),
			ilert.String(ilert.IncidentStatus.Monitoring),
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d open incidents\n", len(open.Incidents))

	result, err := client.GetIncident(&ilert.GetIncidentInput{
		IncidentID: ilert.Int64(incidentID),
		Include:    []*string{ilert.String(ilert.IncidentInclude.ResponderCountJoined)},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Incident #%d %s (%s), %d responders joined\n", result.Incident.Number, result.Incident.Title, result.Incident.Status, result.Incident.ResponderCountJoined)

	// fields left empty keep their value, the ETag makes the API reject the update with a 412
	// when the incident changed since it was read
	updated, err := client.UpdateIncident(&ilert.UpdateIncidentInput{
		IncidentID: ilert.Int64(incidentID),
		Incident:   &ilert.Incident{Status: ilert.IncidentStatus.Identified},
		ETag:       result.ETag,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Incident is now %s\n", updated.Incident.Status)

	userID := int64(0) // your specific user id
	_, err = client.AddIncidentSubscribers(&ilert.AddIncidentSubscribersInput{
		IncidentID:  ilert.Int64(incidentID),
		Subscribers: &[]ilert.Subscriber{{ID: userID, Type: ilert.SubscriberType.User}},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	subscribers, err := client.GetIncidentSubscribers(&ilert.GetIncidentSubscribersInput{IncidentID: ilert.Int64(incidentID)})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for _, subscriber := range subscribers.Subscribers {
		log.Printf("Subscriber: %+v\n", *subscriber)
	}

	// a status update has to affect at least one service, the API rejects it otherwise
	statusUpdate, err := client.CreateIncidentStatusUpdate(&ilert.CreateIncidentStatusUpdateInput{
		IncidentID: ilert.Int64(incidentID),
		StatusUpdate: &ilert.StatusUpdate{
			Summary:          result.Incident.Title,
			Status:           ilert.StatusUpdateStatus.Identified,
			Message:          "We identified the cause and are working on a fix.",
			SendNotification: true,
			AffectedServices: incident.AffectedServices,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Published status update %d\n", statusUpdate.StatusUpdate.ID)

	logEntries, err := client.GetIncidentLogEntries(&ilert.GetIncidentLogEntriesInput{IncidentID: ilert.Int64(incidentID)})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for _, entry := range logEntries.LogEntries {
		log.Printf("%s %s/%s\n", entry.Timestamp, entry.Type, entry.SubType)
	}
}
