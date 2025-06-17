package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createCallFlowInput := &ilert.CreateCallFlowInput{
		CallFlow: &ilert.CallFlow{
			Name:     "call flow example",
			Language: "en",
			RootNode: &ilert.CallFlowNode{
				NodeType: "ROOT",
				Branches: []ilert.CallFlowBranch{
					{
						BranchType: "ANSWERED",
						Target: &ilert.CallFlowNode{
							NodeType: "CREATE_ALERT",
							Name:     "Create alert",
							Metadata: &ilert.CallFlowNodeMetadata{
								AlertSourceId: -1, // your call flow alert source id
							},
						},
					},
				},
			},
		},
	}

	result, err := client.CreateCallFlow(createCallFlowInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Call flow:\n\n %+v\n", result.CallFlow)
}
