package main

import (
	"log"

	"github.com/iLert/ilert-go/v2"
)

func main() {
	var apiKey = "heartbeat API Key"
	client := ilert.NewClient()
	result, err := client.PingHeartbeat(&ilert.PingHeartbeatInput{
		APIKey: ilert.String(apiKey),
		Method: ilert.String(ilert.HeartbeatMethods.HEAD),
	})
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println("Heartbeat is ok!")
}
