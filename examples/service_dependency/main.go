package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	sourceServiceId := int64(0) // the service that depends on the target service
	targetServiceId := int64(0) // the service that is being depended upon

	createResult, err := client.CreateServiceDependency(&ilert.CreateServiceDependencyInput{
		ServiceID: &sourceServiceId,
		Dependency: &ilert.ServiceDependency{
			TargetServiceID: targetServiceId,
			Type:            ilert.ServiceDependencyType.Hard,
			Notes:           "the source service reads from the target service",
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Service dependency:\n\n %+v\n", createResult.Dependency)

	listResult, err := client.GetServiceDependencies(&ilert.GetServiceDependenciesInput{
		ServiceID: &sourceServiceId,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Service has %d dependencies\n", len(listResult.Dependencies))
}
