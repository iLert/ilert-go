package main

import (
	"fmt"
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))
	serviceId := int64(0) //your specific service id
	input := &ilert.GetServiceInput{
		ServiceID: &serviceId,
		Include:   []*string{&ilert.ServiceInclude.Uptime},
	}
	result, err := client.GetService(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Println(fmt.Sprintf("Service:\n\n %+v\n", result.Service))
}
