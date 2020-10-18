package main

import (
	"log"

	"github.com/iLert/ilert-go"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	result, err := client.GetCurrentUser()
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println("User id:", result.User.ID)
}
