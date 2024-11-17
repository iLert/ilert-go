package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createDeploymentPipelineInput := &ilert.CreateDeploymentPipelineInput{
		DeploymentPipeline: &ilert.DeploymentPipeline{
			Name:            "example",
			IntegrationType: ilert.IntegrationType.GitHub,
			Params: &ilert.DeploymentPipelineGitHubParams{
				BranchFilters: []string{"main", "master"},
				EventFilters:  []string{ilert.GitHubEventFilterType.Release},
			},
		},
	}

	result, err := client.CreateDeploymentPipeline(createDeploymentPipelineInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Deployment pipeline:\n\n %+v\n", result.DeploymentPipeline)
}
