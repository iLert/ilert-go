package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	var org = "your organization"
	var username = "your username"
	var password = "your password"
	client := ilert.NewClient(ilert.WithBasicAuth(org, username, password))
	result, err := client.GetAlertSources(&ilert.GetAlertSourcesInput{})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Found %d alert sources\n\n ", len(result.AlertSources)))
	for _, alertSource := range result.AlertSources {
		log.Println(fmt.Sprintf("%+v\n", *alertSource))
	}

}
