package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	result, err := client.GetServiceTopology(&ilert.GetServiceTopologyInput{
		// a single label condition expression, not the repeated "key:value" list filter
		Labels: ilert.String("environment == 'production'"),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Graph of %d services and %d dependencies\n\n", len(result.ServiceGraph.Nodes), len(result.ServiceGraph.Edges))
	for _, node := range result.ServiceGraph.Nodes {
		log.Printf("service %d %q: %s\n", node.ID, node.Name, node.Status)
	}
	for _, edge := range result.ServiceGraph.Edges {
		// the graph only carries the rendering projection of an edge, read the full
		// dependency through GetServiceDependency
		log.Printf("%d depends on %d (%s)\n", edge.SourceServiceID, edge.TargetServiceID, edge.Type)
	}
}
