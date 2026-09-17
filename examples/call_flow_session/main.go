package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetCallFlowSessions(&ilert.GetCallFlowSessionsInput{
		// the API accepts a single state only
		State:      ilert.String(ilert.CallFlowSessionStatus.Ended),
		From:       ilert.String("2026-09-01T00:00:00Z"),
		MaxResults: ilert.Int(25),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d call flow sessions\n\n", len(result.CallFlowSessions))

	for _, session := range result.CallFlowSessions {
		log.Printf("session %d: %s -> %s (%s)\n", session.ID, session.FromNumber, session.ToNumber, session.Status)
	}

	if len(result.CallFlowSessions) == 0 {
		return
	}

	session := result.CallFlowSessions[0]
	detail, err := client.GetCallFlowSession(&ilert.GetCallFlowSessionInput{CallFlowSessionID: &session.ID})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Session:\n\n %+v\n", detail.CallFlowSession)

	// the keys of the collected data come from the nodes of the flow that ran
	data, err := client.GetCallFlowSessionData(&ilert.GetCallFlowSessionDataInput{CallFlowSessionID: &session.ID})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for key, value := range data.Data {
		log.Printf("%s = %v\n", key, value)
	}
}
