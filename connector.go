package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Connector definition
type Connector struct {
	ID        string      `json:"id,omitempty"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	CreatedAt string      `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt string      `json:"updatedAt,omitempty"` // date time string in ISO 8601
	Params    interface{} `json:"params"`
}

// ConnectorOutput definition
type ConnectorOutput struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Type      string                `json:"type"`
	CreatedAt string                `json:"createdAt"` // date time string in ISO 8601
	UpdatedAt string                `json:"updatedAt"` // date time string in ISO 8601
	Params    ConnectorOutputParams `json:"params"`
}

// ConnectorOutputParams definition
type ConnectorOutputParams struct {
	APIKey        string `json:"apiKey,omitempty"`        // Datadog or Zendesk or Github or Serverless or Autotask api key
	Authorization string `json:"authorization,omitempty"` // Serverless
	URL           string `json:"url,omitempty"`           // Jira or Microsoft Teams or Zendesk or Discord or Autotask server url
	Email         string `json:"email,omitempty"`         // Jira or ServiceNow or Zendesk username or email
	Username      string `json:"username,omitempty"`      // TOPdesk or ServiceNow or Autotask username
	Password      string `json:"password,omitempty"`      // Jira or ServiceNow or Autotask user password or api token
	Secret        string `json:"secret,omitempty"`        // DingTalk
}

// ConnectorParamsDatadog definition
type ConnectorParamsDatadog struct {
	APIKey string `json:"apiKey"`
}

// ConnectorParamsJira definition
type ConnectorParamsJira struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ConnectorParamsMicrosoftTeams definition
type ConnectorParamsMicrosoftTeams struct {
	URL string `json:"url"`
}

// ConnectorParamsServiceNow definition
type ConnectorParamsServiceNow struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ConnectorParamsSlack definition
type ConnectorParamsSlack struct{}

// ConnectorParamsZendesk definition
type ConnectorParamsZendesk struct {
	URL    string `json:"url"`
	Email  string `json:"email"`
	APIKey string `json:"apiKey"`
}

// ConnectorParamsDiscord definition
type ConnectorParamsDiscord struct {
	URL string `json:"url"`
}

// ConnectorParamsGithub definition
type ConnectorParamsGithub struct {
	APIKey string `json:"apiKey"`
}

// ConnectorParamsTopdesk definition
type ConnectorParamsTopdesk struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ConnectorParamsAWSLambda definition
type ConnectorParamsAWSLambda struct {
	Authorization string `json:"authorization,omitempty"`
}

// ConnectorParamsAzureFunction definition
type ConnectorParamsAzureFunction struct {
	Authorization string `json:"authorization,omitempty"`
}

// ConnectorParamsGoogleFunction definition
type ConnectorParamsGoogleFunction struct {
	Authorization string `json:"authorization,omitempty"`
}

// ConnectorParamsSysdig definition
type ConnectorParamsSysdig struct {
	APIKey string `json:"apiKey"`
}

// ConnectorParamsAutotask definition
type ConnectorParamsAutotask struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ConnectorParamsMattermost definition
type ConnectorParamsMattermost struct {
	URL string `json:"url"`
}

// ConnectorParamsZammad definition
type ConnectorParamsZammad struct {
	URL    string `json:"url"`
	APIKey string `json:"apiKey"`
}

// ConnectorParamsStatusPageIO definition
type ConnectorParamsStatusPageIO struct {
	APIKey string `json:"apiKey"`
}

// ConnectorParamsDingTalk definition
type ConnectorParamsDingTalk struct {
	URL    string `json:"url,omitempty"`
	Secret string `json:"secret,omitempty"`
}

// ConnectorTypes defines connector types
var ConnectorTypes = struct {
	AutomationRule        string
	Autotask              string
	DingTalk              string
	DingTalkAction        string
	Discord               string
	Email                 string
	Github                string
	Jira                  string
	Mattermost            string
	MicrosoftTeams        string
	MicrosoftTeamsBot     string
	MicrosoftTeamsWebhook string
	ServiceNow            string
	Slack                 string
	SlackWebhook          string
	Telegram              string
	Topdesk               string
	Webhook               string
	Zammad                string
	Zendesk               string
}{
	AutomationRule:        "automation_rule",
	Autotask:              "autotask",
	DingTalk:              "dingtalk",
	DingTalkAction:        "dingtalk_action",
	Discord:               "discord",
	Email:                 "email",
	Github:                "github",
	Jira:                  "jira",
	Mattermost:            "mattermost",
	MicrosoftTeams:        "microsoft_teams",
	MicrosoftTeamsBot:     "microsoft_teams_bot",
	MicrosoftTeamsWebhook: "microsoft_teams_webhook",
	ServiceNow:            "servicenow",
	Slack:                 "slack",
	SlackWebhook:          "slack_webhook",
	Telegram:              "telegram",
	Topdesk:               "topdesk",
	Webhook:               "webhook",
	Zammad:                "zammad",
	Zendesk:               "zendesk",
}

// ConnectorTypesAll defines connector all types list
var ConnectorTypesAll = []string{
	ConnectorTypes.AutomationRule,
	ConnectorTypes.Autotask,
	ConnectorTypes.DingTalk,
	ConnectorTypes.DingTalkAction,
	ConnectorTypes.Discord,
	ConnectorTypes.Email,
	ConnectorTypes.Github,
	ConnectorTypes.Jira,
	ConnectorTypes.Mattermost,
	ConnectorTypes.MicrosoftTeams,
	ConnectorTypes.MicrosoftTeamsBot,
	ConnectorTypes.MicrosoftTeamsWebhook,
	ConnectorTypes.ServiceNow,
	ConnectorTypes.Slack,
	ConnectorTypes.SlackWebhook,
	ConnectorTypes.Telegram,
	ConnectorTypes.Topdesk,
	ConnectorTypes.Webhook,
	ConnectorTypes.Zammad,
	ConnectorTypes.Zendesk,
}

// CreateConnectorInput represents the input of a CreateConnector operation.
type CreateConnectorInput struct {
	_         struct{}
	Connector *Connector
}

// CreateConnectorOutput represents the output of a CreateConnector operation.
type CreateConnectorOutput struct {
	_         struct{}
	Connector *ConnectorOutput
}

// CreateConnector creates a new connector. https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
func (c *Client) CreateConnector(input *CreateConnectorInput) (*CreateConnectorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Connector == nil {
		return nil, errors.New("connector input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Connector).Post(apiRoutes.connectors)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	connector := &ConnectorOutput{}
	err = json.Unmarshal(resp.Body(), connector)
	if err != nil {
		return nil, err
	}

	return &CreateConnectorOutput{Connector: connector}, nil
}

// GetConnectorInput represents the input of a GetConnector operation.
type GetConnectorInput struct {
	_           struct{}
	ConnectorID *string
}

// GetConnectorOutput represents the output of a GetConnector operation.
type GetConnectorOutput struct {
	_         struct{}
	Connector *ConnectorOutput
}

// GetConnector gets the connector with specified id. https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors~1{id}/get
func (c *Client) GetConnector(input *GetConnectorInput) (*GetConnectorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ConnectorID == nil {
		return nil, errors.New("connector id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%s", apiRoutes.connectors, *input.ConnectorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	connector := &ConnectorOutput{}
	err = json.Unmarshal(resp.Body(), connector)
	if err != nil {
		return nil, err
	}

	return &GetConnectorOutput{Connector: connector}, nil
}

// GetConnectorsInput represents the input of a GetConnectors operation.
type GetConnectorsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	// Default: 0
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults *int
}

// GetConnectorsOutput represents the output of a GetConnectors operation.
type GetConnectorsOutput struct {
	_          struct{}
	Connectors []*ConnectorOutput
}

// GetConnectors lists existing connectors. https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/get
func (c *Client) GetConnectors(input *GetConnectorsInput) (*GetConnectorsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	} else {
		q.Add("start-index", "0")
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	} else {
		q.Add("max-results", "50")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.connectors, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	connectors := make([]*ConnectorOutput, 0)
	err = json.Unmarshal(resp.Body(), &connectors)
	if err != nil {
		return nil, err
	}

	return &GetConnectorsOutput{Connectors: connectors}, nil
}

// SearchConnectorInput represents the input of a SearchConnector operation.
type SearchConnectorInput struct {
	_             struct{}
	ConnectorName *string
}

// SearchConnectorOutput represents the output of a SearchConnector operation.
type SearchConnectorOutput struct {
	_         struct{}
	Connector *ConnectorOutput
}

// SearchConnector gets the connector with specified name.
func (c *Client) SearchConnector(input *SearchConnectorInput) (*SearchConnectorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ConnectorName == nil {
		return nil, errors.New("connector name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.connectors, *input.ConnectorName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	connector := &ConnectorOutput{}
	err = json.Unmarshal(resp.Body(), connector)
	if err != nil {
		return nil, err
	}

	return &SearchConnectorOutput{Connector: connector}, nil
}

// UpdateConnectorInput represents the input of a UpdateConnector operation.
type UpdateConnectorInput struct {
	_           struct{}
	ConnectorID *string
	Connector   *Connector
}

// UpdateConnectorOutput represents the output of a UpdateConnector operation.
type UpdateConnectorOutput struct {
	_         struct{}
	Connector *ConnectorOutput
}

// UpdateConnector updates an existing connector. https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors~1{id}/put
func (c *Client) UpdateConnector(input *UpdateConnectorInput) (*UpdateConnectorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Connector == nil {
		return nil, errors.New("connector input is required")
	}
	if input.ConnectorID == nil {
		return nil, errors.New("connector id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Connector).Put(fmt.Sprintf("%s/%s", apiRoutes.connectors, *input.ConnectorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	connector := &ConnectorOutput{}
	err = json.Unmarshal(resp.Body(), connector)
	if err != nil {
		return nil, err
	}

	return &UpdateConnectorOutput{Connector: connector}, nil
}

// DeleteConnectorInput represents the input of a DeleteConnector operation.
type DeleteConnectorInput struct {
	_           struct{}
	ConnectorID *string
}

// DeleteConnectorOutput represents the output of a DeleteConnector operation.
type DeleteConnectorOutput struct {
	_ struct{}
}

// DeleteConnector deletes the specified connector. https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors~1{id}/delete
func (c *Client) DeleteConnector(input *DeleteConnectorInput) (*DeleteConnectorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ConnectorID == nil {
		return nil, errors.New("connector id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%s", apiRoutes.connectors, *input.ConnectorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteConnectorOutput{}, nil
}
