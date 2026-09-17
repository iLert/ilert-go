package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetCallFlowNumbers(&ilert.GetCallFlowNumbersInput{
		// only the numbers that are not assigned to a call flow yet
		State: ilert.String(ilert.CallFlowNumberStateFilter.Available),
		// the call flow a number is assigned to is only part of the response when requested
		Include: []*string{&ilert.CallFlowNumberInclude.AssignedTo},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d call flow numbers\n\n", len(result.CallFlowNumbers))
	for _, number := range result.CallFlowNumbers {
		log.Printf("%+v\n", *number)
	}

	if len(result.CallFlowNumbers) == 0 {
		return
	}

	getResult, err := client.GetCallFlowNumber(&ilert.GetCallFlowNumberInput{
		CallFlowNumberID: &result.CallFlowNumbers[0].ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Call flow number by id:\n\n %+v\n", getResult.CallFlowNumber)

	// the name lookup matches the phone number first and falls back to the name, and a
	// number matched by its phone number comes back without a state
	searchResult, err := client.SearchCallFlowNumber(&ilert.SearchCallFlowNumberInput{
		CallFlowNumberName: &result.CallFlowNumbers[0].Name,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Call flow number by name:\n\n %+v\n", searchResult.CallFlowNumber)
}
