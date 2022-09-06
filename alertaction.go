package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AlertAction definition https://api.ilert.com/api-docs/#tag/Alert-Actions
type AlertAction struct {
	ID             string      `json:"id,omitempty"`
	Name           string      `json:"name"`
	AlertSourceIDs []int64     `json:"alertSourceIds"`
	ConnectorID    string      `json:"connectorId"`
	ConnectorType  string      `json:"connectorType"`
	TriggerMode    string      `json:"triggerMode"`
	TriggerTypes   []string    `json:"triggerTypes,omitempty"`
	CreatedAt      string      `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt      string      `json:"updatedAt,omitempty"` // date time string in ISO 8601
	Params         interface{} `json:"params"`
}

// AlertActionOutput definition https://api.ilert.com/api-docs/#tag/Alert-Actions
type AlertActionOutput struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	AlertSourceIDs []int64                  `json:"alertSourceIds"`
	ConnectorID    string                   `json:"connectorId"`
	ConnectorType  string                   `json:"connectorType"`
	TriggerMode    string                   `json:"triggerMode"`
	TriggerTypes   []string                 `json:"triggerTypes,omitempty"`
	CreatedAt      string                   `json:"createdAt"` // date time string in ISO 8601
	UpdatedAt      string                   `json:"updatedAt"` // date time string in ISO 8601
	Params         *AlertActionOutputParams `json:"params"`
}

// AlertActionOutputParams definition
type AlertActionOutputParams struct {
	BodyTemplate    string   `json:"bodyTemplate,omitempty"`    // Custom, Jira, Email. Used to post data to external server
	CallerID        string   `json:"callerId,omitempty"`        // ServiceNow: user email
	ChannelID       string   `json:"channelId,omitempty"`       // Slack
	ChannelName     string   `json:"channelName,omitempty"`     // Slack
	CompanyID       int64    `json:"companyId,omitempty"`       // Autotask: Company ID
	EventFilter     string   `json:"eventFilter,omitempty"`     // Sysdig
	Impact          string   `json:"impact,omitempty"`          // ServiceNow: 1 - High, 2 - Medium, 3 - Low (Default)
	IssueType       string   `json:"issueType,omitempty"`       // Jira: "Bug" | "Epic" | "Subtask" | "Story" | "Task"
	IssueTypeNumber int64    `json:"issueTypeNumber,omitempty"` // Autotask: Issue type
	Labels          []string `json:"labels,omitempty"`          // Github
	Name            string   `json:"name,omitempty"`            // Jira or MicrosoftTeams or Zendesk
	Owner           string   `json:"owner,omitempty"`           // Github
	Priority        string   `json:"priority,omitempty"`        // Datadog: "normal" | "low". Zendesk: "urgent" | "high" | "normal" | "low".
	Project         string   `json:"project,omitempty"`         // Jira
	QueueID         int64    `json:"queueId,omitempty"`         // Autotask: Queue ID
	Recipients      []string `json:"recipients,omitempty"`      // Email
	Repository      string   `json:"repository,omitempty"`      // Github
	Site            string   `json:"site,omitempty"`            // Datadog: default `US`. Values: `US` or `EU`
	Status          string   `json:"status,omitempty"`          // Topdesk: firstLine, secondLine, partial
	Subject         string   `json:"subject,omitempty"`         // Email
	Tags            []string `json:"tags,omitempty"`            // Datadog or Sysdig
	TeamDomain      string   `json:"teamDomain,omitempty"`      // Slack
	TeamID          string   `json:"teamId,omitempty"`          // Slack
	TicketCategory  string   `json:"ticketCategory,omitempty"`  // Autotask
	TicketType      string   `json:"ticketType,omitempty"`      // Autotask
	Urgency         string   `json:"urgency,omitempty"`         // ServiceNow: 1 - High, 2 - Medium, 3 - Low (Default)
	WebhookURL      string   `json:"webhookUrl,omitempty"`      // Custom
	Email           string   `json:"email,omitempty"`           // Zammad
	PageID          string   `json:"pageId,omitempty"`          // StatusPage.io
}

// AlertActionParamsAutotask definition
type AlertActionParamsAutotask struct {
	CompanyID      string `json:"companyId,omitempty"`      // Autotask: Company ID
	IssueType      string `json:"issueType,omitempty"`      // Autotask: Issue type
	QueueID        int64  `json:"queueId,omitempty"`        // Autotask: Queue ID
	TicketCategory string `json:"ticketCategory,omitempty"` // Autotask ticket category
	TicketType     string `json:"ticketType,omitempty"`     // Autotask ticket type
}

// AlertActionParamsDatadog definition
type AlertActionParamsDatadog struct {
	Tags     []string `json:"tags,omitempty"`
	Priority string   `json:"priority,omitempty"` // "normal" | "low"
	Site     string   `json:"site,omitempty"`     // `US` | `EU`
}

// AlertActionParamsJira definition
type AlertActionParamsJira struct {
	Project      string `json:"project,omitempty"`
	IssueType    string `json:"issueType,omitempty"` // "Bug" | "Epic" | "Subtask" | "Story" | "Task"
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsMicrosoftTeams definition
type AlertActionParamsMicrosoftTeams struct{}

// AlertActionParamsServiceNow definition
type AlertActionParamsServiceNow struct {
	CallerID string `json:"callerId,omitempty"` // user email
	Impact   string `json:"impact,omitempty"`   // 1 - High, 2 - Medium, 3 - Low (Default)
	Urgency  string `json:"urgency,omitempty"`  // 1 - High, 2 - Medium, 3 - Low (Default)
}

// AlertActionParamsSlack definition
type AlertActionParamsSlack struct {
	ChannelID   string `json:"channelId,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
	TeamDomain  string `json:"teamDomain,omitempty"`
	TeamID      string `json:"teamId,omitempty"`
}

// AlertActionParamsWebhook definition
type AlertActionParamsWebhook struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsZendesk definition
type AlertActionParamsZendesk struct {
	Priority string `json:"priority,omitempty"` // "urgent" | "high" | "normal" | "low"
}

// AlertActionParamsDiscord definition
type AlertActionParamsDiscord struct{}

// AlertActionParamsGithub definition
type AlertActionParamsGithub struct {
	Owner      string   `json:"owner,omitempty"`
	Repository string   `json:"repository,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

// AlertActionParamsTopdesk definition
type AlertActionParamsTopdesk struct {
	Status string `json:"status,omitempty"` // `firstLine`| `secondLine` | `partial`
}

// AlertActionParamsAWSLambda definition
type AlertActionParamsAWSLambda struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsAzureFunction definition
type AlertActionParamsAzureFunction struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsGoogleFunction definition
type AlertActionParamsGoogleFunction struct {
	WebhookURL   string `json:"webhookUrl,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsEmail definition
type AlertActionParamsEmail struct {
	Recipients   []string `json:"recipients,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	BodyTemplate string   `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsSysdig definition
type AlertActionParamsSysdig struct {
	Tags        []string `json:"tags,omitempty"`
	EventFilter string   `json:"eventFilter,omitempty"`
}

// AlertActionParamsZapier definition
type AlertActionParamsZapier struct {
	WebhookURL string `json:"webhookUrl,omitempty"`
}

// AlertActionParamsZammad definition
type AlertActionParamsZammad struct {
	Email string `json:"email,omitempty"`
}

// AlertActionParamsStatusPageIO definition
type AlertActionParamsStatusPageIO struct {
	PageID string `json:"pageId,omitempty"`
}

// AlertActionResult definition
type AlertActionResult struct {
	ID          string `json:"id"`
	WebhookID   string `json:"webhookId"`
	ExtensionID string `json:"extensionId"`
	AlertID     int64  `json:"alertId"`
	Success     bool   `json:"success"`
	Actor       User   `json:"actor"`
}

// AlertActionTriggerModes defines alertAction trigger modes
var AlertActionTriggerModes = struct {
	Automatic string
	Manual    string
}{
	Automatic: "AUTOMATIC",
	Manual:    "MANUAL",
}

// AlertActionTriggerModesAll defines alertAction trigger modes list
var AlertActionTriggerModesAll = []string{
	AlertActionTriggerModes.Automatic,
	AlertActionTriggerModes.Manual,
}

// AlertActionTriggerTypes defines alertAction trigger types
var AlertActionTriggerTypes = struct {
	AlertCreated       string
	AlertAssigned      string
	AlertAutoEscalated string
	AlertAcknowledged  string
	AlertRaised        string
	AlertCommentAdded  string
	AlertResolved      string
}{
	AlertCreated:       "alert-created",
	AlertAssigned:      "alert-assigned",
	AlertAutoEscalated: "alert-auto-escalated",
	AlertAcknowledged:  "alert-acknowledged",
	AlertRaised:        "alert-raised",
	AlertCommentAdded:  "alert-comment-added",
	AlertResolved:      "alert-resolved",
}

// AlertActionTriggerTypesAll defines all alertAction trigger types
var AlertActionTriggerTypesAll = []string{
	AlertActionTriggerTypes.AlertCreated,
	AlertActionTriggerTypes.AlertAssigned,
	AlertActionTriggerTypes.AlertAutoEscalated,
	AlertActionTriggerTypes.AlertAcknowledged,
	AlertActionTriggerTypes.AlertRaised,
	AlertActionTriggerTypes.AlertCommentAdded,
	AlertActionTriggerTypes.AlertResolved,
}

// CreateAlertActionInput represents the input of a CreateAlertAction operation.
type CreateAlertActionInput struct {
	_           struct{}
	AlertAction *AlertAction
}

// CreateAlertActionOutput represents the output of a CreateAlertAction operation.
type CreateAlertActionOutput struct {
	_           struct{}
	AlertAction *AlertActionOutput
}

// CreateAlertAction creates a new alertAction https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions/post
func (c *Client) CreateAlertAction(input *CreateAlertActionInput) (*CreateAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertAction == nil {
		return nil, errors.New("AlertAction input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.AlertAction).Post(apiRoutes.alertActions)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	alertAction := &AlertActionOutput{}
	err = json.Unmarshal(resp.Body(), alertAction)
	if err != nil {
		return nil, err
	}

	return &CreateAlertActionOutput{AlertAction: alertAction}, nil
}

// GetAlertActionInput represents the input of a GetAlertAction operation.
type GetAlertActionInput struct {
	_             struct{}
	AlertActionID *string
}

// GetAlertActionOutput represents the output of a GetAlertAction operation.
type GetAlertActionOutput struct {
	_           struct{}
	AlertAction *AlertActionOutput
}

// GetAlertAction gets the alertAction with specified id. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions~1{id}/get
func (c *Client) GetAlertAction(input *GetAlertActionInput) (*GetAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertActionID == nil {
		return nil, errors.New("AlertAction id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%s", apiRoutes.alertActions, *input.AlertActionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertAction := &AlertActionOutput{}
	err = json.Unmarshal(resp.Body(), alertAction)
	if err != nil {
		return nil, err
	}

	return &GetAlertActionOutput{AlertAction: alertAction}, nil
}

// GetAlertActionsInput represents the input of a GetAlertActions operation.
type GetAlertActionsInput struct {
	_ struct{}
}

// GetAlertActionsOutput represents the output of a GetAlertActions operation.
type GetAlertActionsOutput struct {
	_            struct{}
	AlertActions []*AlertActionOutput
}

// GetAlertActions lists alertActions. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions/get
func (c *Client) GetAlertActions(input *GetAlertActionsInput) (*GetAlertActionsOutput, error) {
	resp, err := c.httpClient.R().Get(apiRoutes.alertActions)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertActions := make([]*AlertActionOutput, 0)
	err = json.Unmarshal(resp.Body(), &alertActions)
	if err != nil {
		return nil, err
	}

	return &GetAlertActionsOutput{AlertActions: alertActions}, nil
}

// SearchAlertActionInput represents the input of a SearchAlertAction operation.
type SearchAlertActionInput struct {
	_               struct{}
	AlertActionName *string
}

// SearchAlertActionOutput represents the output of a SearchAlertAction operation.
type SearchAlertActionOutput struct {
	_           struct{}
	AlertAction *AlertAction
}

// SearchAlertAction gets the alertAction with specified name.
func (c *Client) SearchAlertAction(input *SearchAlertActionInput) (*SearchAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertActionName == nil {
		return nil, errors.New("alert action name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.alertActions, *input.AlertActionName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertAction := &AlertAction{}
	err = json.Unmarshal(resp.Body(), alertAction)
	if err != nil {
		return nil, err
	}

	return &SearchAlertActionOutput{AlertAction: alertAction}, nil
}

// UpdateAlertActionInput represents the input of a UpdateAlertAction operation.
type UpdateAlertActionInput struct {
	_             struct{}
	AlertActionID *string
	AlertAction   *AlertAction
}

// UpdateAlertActionOutput represents the output of a UpdateAlertAction operation.
type UpdateAlertActionOutput struct {
	_           struct{}
	AlertAction *AlertActionOutput
}

// UpdateAlertAction updates an existing alertAction. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions~1{id}/put
func (c *Client) UpdateAlertAction(input *UpdateAlertActionInput) (*UpdateAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertAction == nil {
		return nil, errors.New("AlertAction input is required")
	}
	if input.AlertActionID == nil {
		return nil, errors.New("AlertAction id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.AlertAction).Put(fmt.Sprintf("%s/%s", apiRoutes.alertActions, *input.AlertActionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertAction := &AlertActionOutput{}
	err = json.Unmarshal(resp.Body(), alertAction)
	if err != nil {
		return nil, err
	}

	return &UpdateAlertActionOutput{AlertAction: alertAction}, nil
}

// DeleteAlertActionInput represents the input of a DeleteAlertAction operation.
type DeleteAlertActionInput struct {
	_             struct{}
	AlertActionID *string
}

// DeleteAlertActionOutput represents the output of a DeleteAlertAction operation.
type DeleteAlertActionOutput struct {
	_ struct{}
}

// DeleteAlertAction deletes the specified alert source. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions~1{id}/delete
func (c *Client) DeleteAlertAction(input *DeleteAlertActionInput) (*DeleteAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertActionID == nil {
		return nil, errors.New("AlertAction id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%s", apiRoutes.alertActions, *input.AlertActionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteAlertActionOutput{}, nil
}
