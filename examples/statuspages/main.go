package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	newStatuspage := &ilert.StatusPage{
		Name:       "your statuspage name",
		Subdomain:  "your subdomain name .ilert.io",
		Visibility: "PRIVATE",
		Services:   []ilert.Service{{ID: 0}}, // your service id
		Timezone:   "Europe/Berlin",
	}
	input := &ilert.CreateStatusPageInput{
		StatusPage: newStatuspage,
	}
	result, err := client.CreateStatusPage(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Statuspage:\n\n %+v\n", result.StatusPage))
}
