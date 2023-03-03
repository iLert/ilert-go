package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	result, err := client.GetCurrentUser()
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err.Error())
	}
	log.Printf("User:\n\n %+v\n", result.User)
}
