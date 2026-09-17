package main

import (
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
	log.Printf("Service:\n\n %+v\n", result.Service)

	// a status change that was kept internal reads back with a public status that differs
	// from the status, and can be published onto the status pages afterwards
	statusResult, err := client.GetService(&ilert.GetServiceInput{
		ServiceID: &serviceId,
		Include:   []*string{&ilert.ServiceInclude.PublicStatus},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	if statusResult.Service.PublicStatus != statusResult.Service.Status {
		// the status guards against publishing a change that happened in the meantime,
		// the API answers 409 when it no longer matches
		published, err := client.PublishServiceStatus(&ilert.PublishServiceStatusInput{
			ServiceID: &serviceId,
			Status:    &statusResult.Service.Status,
		})
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		log.Printf("Published status of service %d: %s\n", published.Service.ID, published.Service.Status)
	}
}
