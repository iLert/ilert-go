package main

import (
	"log"

	"github.com/iLert/ilert-go/v4"
)

func main() {
	var apiToken = "your API token"
	// agents and credentials are early access APIs, so the client has to opt in to
	// both, and then sends the ilert-beta header the API requires for them
	client := ilert.NewClient(
		ilert.WithAPIToken(apiToken),
		ilert.WithEarlyAccess(ilert.EarlyAccess.Agents, ilert.EarlyAccess.Credentials),
	)

	// a credential of a built-in MCP server is in the MCP category, with the short
	// name of the server as sub category
	credentialResult, err := client.CreateCredential(&ilert.CreateCredentialInput{
		Credential: &ilert.Credential{
			Name:        "github mcp",
			TypeHint:    ilert.CredentialType.Bearer,
			SecretValue: "your GitHub token",
			Category:    "MCP",
			SubCategory: "github",
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Credential:\n\n %+v\n", credentialResult.Credential)

	createResult, err := client.CreateAgent(&ilert.CreateAgentInput{
		Agent: &ilert.Agent{
			Name:   "triage",
			Status: ilert.AgentStatus.Enabled,
			// servers are referenced by their name, not their short name
			McpServers: []ilert.AgentMcpServer{{
				Name:        "ilert-workforce-github-mcp",
				Active:      true,
				Credentials: &ilert.AgentMcpServerCredentials{ID: credentialResult.Credential.ID},
			}},
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Agent:\n\n %+v\n", createResult.Agent)

	// an update replaces the agent, so change it as it was read
	agent := createResult.Agent
	agent.Status = ilert.AgentStatus.Disabled
	updateResult, err := client.UpdateAgent(&ilert.UpdateAgentInput{AgentID: &agent.ID, Agent: agent})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Updated agent:\n\n %+v\n", updateResult.Agent)
}
