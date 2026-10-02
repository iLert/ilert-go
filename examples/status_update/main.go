package main

import (
	"log"

	"github.com/iLert/ilert-go/v4"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	serviceID := int64(0) // your specific service id
	statusUpdate := &ilert.StatusUpdate{
		Summary:          "checkout unavailable",
		Status:           ilert.StatusUpdateStatus.Investigating,
		Message:          "We are investigating failing payments.",
		SendNotification: true,
		AffectedServices: []ilert.AffectedServices{
			{Impact: ilert.ServiceStatus.MajorOutage, Service: ilert.Service{ID: serviceID}},
		},
	}

	affected, err := client.GetStatusUpdateAffected(&ilert.GetStatusUpdateAffectedInput{StatusUpdate: statusUpdate})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Would notify %d public and %d private subscribers\n", affected.Affected.PublicSubscribers, affected.Affected.PrivateSubscribers)

	created, err := client.CreateStatusUpdate(&ilert.CreateStatusUpdateInput{StatusUpdate: statusUpdate})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	statusUpdateID := created.StatusUpdate.ID
	log.Printf("Created status update %d\n", statusUpdateID)

	open, err := client.GetStatusUpdates(&ilert.GetStatusUpdatesInput{
		States: []*string{ilert.String(ilert.StatusUpdateStatus.Investigating)},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d status updates under investigation\n", len(open.StatusUpdates))

	result, err := client.GetStatusUpdate(&ilert.GetStatusUpdateInput{
		StatusUpdateID: ilert.Int64(statusUpdateID),
		Include:        []*string{ilert.String(ilert.StatusUpdateInclude.AffectedTeams), ilert.String(ilert.StatusUpdateInclude.History)},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Status update:\n\n %+v\n", result.StatusUpdate)

	userID := int64(0) // your specific user id
	_, err = client.AddStatusUpdateSubscribers(&ilert.AddStatusUpdateSubscribersInput{
		StatusUpdateID: ilert.Int64(statusUpdateID),
		Subscribers:    &[]ilert.Subscriber{{ID: userID, Type: ilert.SubscriberType.User}},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	subscribers, err := client.GetStatusUpdateSubscribers(&ilert.GetStatusUpdateSubscribersInput{StatusUpdateID: ilert.Int64(statusUpdateID)})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for _, subscriber := range subscribers.Subscribers {
		log.Printf("Subscriber: %+v\n", *subscriber)
	}

	// changing the status or the message appends to the history, and notifies the subscribers
	// because the status update was created with SendNotification
	resolved := result.StatusUpdate
	resolved.Status = ilert.StatusUpdateStatus.Resolved
	resolved.Message = "The issue has been resolved."
	updated, err := client.UpdateStatusUpdate(&ilert.UpdateStatusUpdateInput{
		StatusUpdateID: ilert.Int64(statusUpdateID),
		StatusUpdate:   resolved,
		ETag:           result.ETag,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Status update is now %s\n", updated.StatusUpdate.Status)

	logEntries, err := client.GetStatusUpdateLogEntries(&ilert.GetStatusUpdateLogEntriesInput{StatusUpdateID: ilert.Int64(statusUpdateID)})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for _, entry := range logEntries.LogEntries {
		log.Printf("%s %s/%s\n", entry.Timestamp, entry.Type, entry.SubType)
	}
}
