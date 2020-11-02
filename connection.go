package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Connection definition https://api.ilert.com/api-docs/#!/Connections
type Connection struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	AlertSourceIDs []int64     `json:"alertSourceIds"`
	ConnectorID    string      `json:"connectorId"`
	ConnectorType  string      `json:"connectorType"`
	TriggerMode    string      `json:"triggerMode"`
	TriggerTypes   []string    `json:"triggerTypes,omitempty"`
	CreatedAt      string      `json:"createdAt"` // date time string in ISO 8601
	UpdatedAt      string      `json:"updatedAt"` // date time string in ISO 8601
	Params         interface{} `json:"params"`
}

// ConnectionOutput definition https://api.ilert.com/api-docs/#!/Connections
type ConnectionOutput struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	AlertSourceIDs []int64                 `json:"alertSourceIds"`
	ConnectorID    string                  `json:"connectorId"`
	ConnectorType  string                  `json:"connectorType"`
	TriggerMode    string                  `json:"triggerMode"`
	TriggerTypes   []string                `json:"triggerTypes,omitempty"`
	Params         *ConnectionOutputParams `json:"params"`
}

// ConnectionOutputParams definition
type ConnectionOutputParams struct {
	BodyTemplate string   `json:"bodyTemplate,omitempty"` // Custom, Jira, Email. Used to post data to external server
	CallerID     string   `json:"callerId,omitempty"`     // ServiceNow: user email
	ChannelID    string   `json:"channelId,omitempty"`    // Slack
	ChannelName  string   `json:"channelName,omitempty"`  // Slack
	EventFilter  string   `json:"eventFilter,omitempty"`  // Sysdig
	Impact       string   `json:"impact,omitempty"`       // ServiceNow: 1 - High, 2 - Medium, 3 - Low (Default)
	IssueType    string   `json:"issueType,omitempty"`    // Jira: "Bug" | "Epic" | "Subtask" | "Story" | "Task"
	Labels       []string `json:"labels,omitempty"`       // Github
	Name         string   `json:"name,omitempty"`         // Jira or MicrosoftTeams or Zendesk
	Owner        string   `json:"owner,omitempty"`        // Github
	Priority     string   `json:"priority,omitempty"`     // Datadog: "normal" | "low". Zendesk: "urgent" | "high" | "normal" | "low".
	Project      string   `json:"project,omitempty"`      // Jira
	QueueID      int      `json:"queueID,omitempty"`      // Autotask: Queue ID
	Recipients   []string `json:"recipients,omitempty"`   // Email
	Repository   string   `json:"repository,omitempty"`   // Github
	Site         string   `json:"site,omitempty"`         // Datadog: default `US`. Values: `US` or `EU`
	Status       string   `json:"status,omitempty"`       // Topdesk: firstLine, secondLine, partial
	Subject      string   `json:"subject,omitempty"`      // Email
	Tags         []string `json:"tags,omitempty"`         // Datadog or Sysdig
	TeamDomain   string   `json:"teamDomain,omitempty"`   // Slack
	TeamID       string   `json:"teamId,omitempty"`       // Slack
	Urgency      string   `json:"urgency,omitempty"`      // ServiceNow: 1 - High, 2 - Medium, 3 - Low (Default)
	WebhookURL   string   `json:"webhookUrl,omitempty"`   // Custom
}

// ConnectionParamsDatadog definition
type ConnectionParamsDatadog struct {
	Tags     []string `json:"tags,omitempty"`
	Priority string   `json:"priority,omitempty"` // "normal" | "low"
	Site     string   `json:"site,omitempty"`     // `US` | `EU`
}

// ConnectionParamsJira definition
type ConnectionParamsJira struct {
	Project      string `json:"project,omitempty"`
	IssueType    string `json:"issueType,omitempty"` // "Bug" | "Epic" | "Subtask" | "Story" | "Task"
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// ConnectionParamsMicrosoftTeams definition
type ConnectionParamsMicrosoftTeams struct{}

// ConnectionParamsServiceNow definition
type ConnectionParamsServiceNow struct {
	CallerID string `json:"callerId,omitempty"` // user email
	Impact   string `json:"impact,omitempty"`   // 1 - High, 2 - Medium, 3 - Low (Default)
	Urgency  string `json:"urgency,omitempty"`  // 1 - High, 2 - Medium, 3 - Low (Default)
}

// ConnectionParamsSlack definition
type ConnectionParamsSlack struct {
	ChannelID   string `json:"channelId,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
	TeamDomain  string `json:"teamDomain,omitempty"`
	TeamID      string `json:"teamId,omitempty"`
}

// ConnectionParamsWebhook definition
type ConnectionParamsWebhook struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// ConnectionParamsZendesk definition
type ConnectionParamsZendesk struct {
	Priority string `json:"priority,omitempty"` // "urgent" | "high" | "normal" | "low"
}

// ConnectionParamsDiscord definition
type ConnectionParamsDiscord struct{}

// ConnectionParamsGithub definition
type ConnectionParamsGithub struct {
	Owner      string   `json:"owner,omitempty"`
	Repository string   `json:"repository,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

// ConnectionParamsTopdesk definition
type ConnectionParamsTopdesk struct {
	Status string `json:"status,omitempty"` // `firstLine`| `secondLine` | `partial`
}

// ConnectionParamsAWSLambda definition
type ConnectionParamsAWSLambda struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// ConnectionParamsAzureFunction definition
type ConnectionParamsAzureFunction struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// ConnectionParamsGoogleFunction definition
type ConnectionParamsGoogleFunction struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// ConnectionParamsEmail definition
type ConnectionParamsEmail struct {
	Recipients   []string `json:"recipients,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	BodyTemplate string   `json:"bodyTemplate,omitempty"`
}

// ConnectionParamsSysdig definition
type ConnectionParamsSysdig struct {
	Tags        []string `json:"tags,omitempty"`
	EventFilter string   `json:"eventFilter,omitempty"`
}

// ConnectionParamsZapier definition
type ConnectionParamsZapier struct {
	WebhookURL string `json:"webhookUrl,omitempty"`
}

// ConnectionTriggerModes defines connection trigger modes
var ConnectionTriggerModes = struct {
	Automatic string
	Manual    string
}{
	Automatic: "AUTOMATIC",
	Manual:    "MANUAL",
}

// ConnectionTriggerTypes defines connection trigger types
var ConnectionTriggerTypes = struct {
	IncidentCreated       string
	IncidentAssigned      string
	IncidentAutoEscalated string
	IncidentAcknowledged  string
	IncidentRaised        string
	IncidentCommentAdded  string
	IncidentResolved      string
}{
	IncidentCreated:       "incident-created",
	IncidentAssigned:      "incident-assigned",
	IncidentAutoEscalated: "incident-auto-escalated",
	IncidentAcknowledged:  "incident-acknowledged",
	IncidentRaised:        "incident-raised",
	IncidentCommentAdded:  "incident-comment-added",
	IncidentResolved:      "incident-resolved",
}

// ConnectionTriggerTypesAll defines all connection trigger types
var ConnectionTriggerTypesAll = []string{
	ConnectionTriggerTypes.IncidentCreated,
	ConnectionTriggerTypes.IncidentAssigned,
	ConnectionTriggerTypes.IncidentAutoEscalated,
	ConnectionTriggerTypes.IncidentAcknowledged,
	ConnectionTriggerTypes.IncidentRaised,
	ConnectionTriggerTypes.IncidentCommentAdded,
	ConnectionTriggerTypes.IncidentResolved,
}

// CreateConnectionInput represents the input of a CreateConnection operation.
type CreateConnectionInput struct {
	_          struct{}
	Connection *Connection
}

// CreateConnectionOutput represents the output of a CreateConnection operation.
type CreateConnectionOutput struct {
	_          struct{}
	Connection *ConnectionOutput
}

// CreateConnection creates a new connection. https://api.ilert.com/api-docs/#tag/Connections/paths/~1connections/post
func (c *Client) CreateConnection(input *CreateConnectionInput) (*CreateConnectionOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.Connection == nil {
		return nil, errors.New("Connection input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Connection).Post(fmt.Sprintf("%s", apiRoutes.connections))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 201); err != nil {
		return nil, err
	}

	connection := &ConnectionOutput{}
	err = json.Unmarshal(resp.Body(), connection)
	if err != nil {
		return nil, err
	}

	output := &CreateConnectionOutput{Connection: connection}

	return output, nil
}

// GetConnectionInput represents the input of a GetConnection operation.
type GetConnectionInput struct {
	_            struct{}
	ConnectionID *string
}

// GetConnectionOutput represents the output of a GetConnection operation.
type GetConnectionOutput struct {
	_          struct{}
	Connection *ConnectionOutput
}

// GetConnection gets the connection with specified id. https://api.ilert.com/api-docs/#tag/Connections/paths/~1connections~1{id}/get
func (c *Client) GetConnection(input *GetConnectionInput) (*GetConnectionOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.ConnectionID == nil {
		return nil, errors.New("Connection id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%s", apiRoutes.connections, *input.ConnectionID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	connection := &ConnectionOutput{}
	err = json.Unmarshal(resp.Body(), connection)
	if err != nil {
		return nil, err
	}

	output := &GetConnectionOutput{
		Connection: connection,
	}

	return output, nil
}

// GetConnectionsInput represents the input of a GetConnections operation.
type GetConnectionsInput struct {
	_ struct{}
}

// GetConnectionsOutput represents the output of a GetConnections operation.
type GetConnectionsOutput struct {
	_           struct{}
	Connections []*ConnectionOutput
}

// GetConnections lists connections. https://api.ilert.com/api-docs/#tag/Connections/paths/~1connections/get
func (c *Client) GetConnections(input *GetConnectionsInput) (*GetConnectionsOutput, error) {
	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s", apiRoutes.connections))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	connections := make([]*ConnectionOutput, 0)
	err = json.Unmarshal(resp.Body(), &connections)
	if err != nil {
		return nil, err
	}

	output := &GetConnectionsOutput{Connections: connections}

	return output, nil
}

// UpdateConnectionInput represents the input of a UpdateConnection operation.
type UpdateConnectionInput struct {
	_            struct{}
	ConnectionID *string
	Connection   *Connection
}

// UpdateConnectionOutput represents the output of a UpdateConnection operation.
type UpdateConnectionOutput struct {
	_          struct{}
	Connection *ConnectionOutput
}

// UpdateConnection updates an existing connection. https://api.ilert.com/api-docs/#tag/Connections/paths/~1connections~1{id}/put
func (c *Client) UpdateConnection(input *UpdateConnectionInput) (*UpdateConnectionOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.Connection == nil {
		return nil, errors.New("Connection input is required")
	}
	if input.ConnectionID == nil {
		return nil, errors.New("Connection id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Connection).Put(fmt.Sprintf("%s/%s", apiRoutes.connections, *input.ConnectionID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	connection := &ConnectionOutput{}
	err = json.Unmarshal(resp.Body(), connection)
	if err != nil {
		return nil, err
	}

	output := &UpdateConnectionOutput{Connection: connection}

	return output, nil
}

// DeleteConnectionInput represents the input of a DeleteConnection operation.
type DeleteConnectionInput struct {
	_            struct{}
	ConnectionID *string
}

// DeleteConnectionOutput represents the output of a DeleteConnection operation.
type DeleteConnectionOutput struct {
	_ struct{}
}

// DeleteConnection deletes the specified alert source. https://api.ilert.com/api-docs/#tag/Connections/paths/~1connections~1{id}/delete
func (c *Client) DeleteConnection(input *DeleteConnectionInput) (*DeleteConnectionOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.ConnectionID == nil {
		return nil, errors.New("Connection id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%s", apiRoutes.connections, *input.ConnectionID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 204); err != nil {
		return nil, err
	}

	output := &DeleteConnectionOutput{}
	return output, nil
}
