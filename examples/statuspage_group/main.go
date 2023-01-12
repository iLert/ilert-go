package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	newStatusPageGroup := &ilert.StatusPageGroup{
		Name: "your status page group name",
	}
	statusPage := &ilert.StatusPage{ID: 0} // your status page id
	input := &ilert.CreateStatusPageGroupInput{
		StatusPageGroup: newStatusPageGroup,
		StatusPageID:    &statusPage.ID,
	}
	result, err := client.CreateStatusPageGroup(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Statuspage Group:\n\n %+v\n", result.StatusPageGroup))
}
