package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// AlertAction definition https://api.ilert.com/api-docs/#tag/Alert-Actions
type AlertAction struct {
	ID                      string         `json:"id,omitempty"`
	Name                    string         `json:"name"`
	AlertSourceIDs          []int64        `json:"alertSourceIds,omitempty"` // @deprecated
	AlertSources            *[]AlertSource `json:"alertSources,omitempty"`
	ConnectorID             string         `json:"connectorId,omitempty"`
	ConnectorType           string         `json:"connectorType"`
	TriggerMode             string         `json:"triggerMode"`
	DelaySec                int            `json:"delaySec,omitempty"`                // @deprecated
	EscalationEndedDelaySec int            `json:"escalationEndedDelaySec,omitempty"` // between 0 and 7200, used with triggerType 'AlertEscalationEnded'
	NotResolvedDelaySec     int            `json:"notResolvedDelaySec,omitempty"`     // between 0 and 7200, used with triggerType 'AlertNotResolved'
	TriggerTypes            []string       `json:"triggerTypes,omitempty"`
	CreatedAt               string         `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt               string         `json:"updatedAt,omitempty"` // date time string in ISO 8601
	Params                  interface{}    `json:"params"`
	AlertFilter             *AlertFilter   `json:"alertFilter,omitempty"` // @deprecated
	Teams                   *[]TeamShort   `json:"teams,omitempty"`
	Conditions              string         `json:"conditions,omitempty"`
}

// AlertActionOutput definition https://api.ilert.com/api-docs/#tag/Alert-Actions
type AlertActionOutput struct {
	ID                      string                   `json:"id"`
	Name                    string                   `json:"name"`
	AlertSourceIDs          []int64                  `json:"alertSourceIds,omitempty"` // @deprecated
	AlertSources            *[]AlertSource           `json:"alertSources,omitempty"`
	ConnectorID             string                   `json:"connectorId"`
	ConnectorType           string                   `json:"connectorType"`
	TriggerMode             string                   `json:"triggerMode"`
	DelaySec                int                      `json:"delaySec,omitempty"`                // @deprecated
	EscalationEndedDelaySec int                      `json:"escalationEndedDelaySec,omitempty"` // between 0 and 7200, used with triggerType 'AlertEscalationEnded'
	NotResolvedDelaySec     int                      `json:"notResolvedDelaySec,omitempty"`     // between 0 and 7200, used with triggerType 'AlertNotResolved'
	TriggerTypes            []string                 `json:"triggerTypes,omitempty"`
	CreatedAt               string                   `json:"createdAt"` // date time string in ISO 8601
	UpdatedAt               string                   `json:"updatedAt"` // date time string in ISO 8601
	Params                  *AlertActionOutputParams `json:"params"`
	AlertFilter             *AlertFilter             `json:"alertFilter,omitempty"` // @deprecated
	Teams                   *[]TeamShort             `json:"teams,omitempty"`
	Conditions              string                   `json:"conditions,omitempty"`
}

// AlertActionOutputParams definition
type AlertActionOutputParams struct {
	AlertType        string   `json:"alertType"`                  // Automation rule
	AtMobiles        []string `json:"atMobiles,omitempty"`        // DingTalk
	BodyTemplate     string   `json:"bodyTemplate,omitempty"`     // Custom, Jira, Email. Used to post data to external server
	CallerID         string   `json:"callerId,omitempty"`         // ServiceNow: user email
	ChannelID        string   `json:"channelId,omitempty"`        // Slack, Telegram
	ChannelName      string   `json:"channelName,omitempty"`      // Slack
	CompanyID        int64    `json:"companyId,omitempty"`        // Autotask: Company ID
	Email            string   `json:"email,omitempty"`            // Zammad
	EventFilter      string   `json:"eventFilter,omitempty"`      // Sysdig
	Impact           string   `json:"impact,omitempty"`           // ServiceNow: 1 - High, 2 - Medium, 3 - Low (Default)
	IsAtAll          bool     `json:"isAtAll,omitempty"`          // DingTalk
	IssueType        string   `json:"issueType,omitempty"`        // Jira: "Bug" | "Epic" | "Subtask" | "Story" | "Task"
	IssueTypeNumber  int64    `json:"issueTypeNumber,omitempty"`  // Autotask: Issue type
	Labels           []string `json:"labels,omitempty"`           // Github
	Name             string   `json:"name,omitempty"`             // Jira or MicrosoftTeams or Zendesk
	Owner            string   `json:"owner,omitempty"`            // Github
	PageID           string   `json:"pageId,omitempty"`           // StatusPage.io
	Priority         string   `json:"priority,omitempty"`         // Datadog: "normal" | "low". Zendesk: "urgent" | "high" | "normal" | "low".
	Project          string   `json:"project,omitempty"`          // Jira
	QueueID          int64    `json:"queueId,omitempty"`          // Autotask: Queue ID
	Recipients       []string `json:"recipients,omitempty"`       // Email
	Repository       string   `json:"repository,omitempty"`       // Github
	ResolveIncident  bool     `json:"resolveIncident,omitempty"`  // Automation rule
	Secret           string   `json:"secret,omitempty"`           // DingTalk
	SendNotification bool     `json:"sendNotification,omitempty"` // Automation rule
	ServiceIds       []int64  `json:"serviceIds"`                 // Automation rule
	ServiceStatus    string   `json:"serviceStatus"`              // Automation rule
	Site             string   `json:"site,omitempty"`             // Datadog: default `US`. Values: `US` or `EU`
	Status           string   `json:"status,omitempty"`           // Topdesk: firstLine, secondLine, partial
	Subject          string   `json:"subject,omitempty"`          // Email
	Tags             []string `json:"tags,omitempty"`             // Datadog or Sysdig
	TeamDomain       string   `json:"teamDomain,omitempty"`       // Slack
	TeamID           string   `json:"teamId,omitempty"`           // Slack, Microsoft Teams Bot
	TeamName         string   `json:"teamName,omitempty"`         // Slack, Microsoft Teams Bot
	TemplateId       int64    `json:"templateId,omitempty"`       // Automation rule
	TicketCategory   string   `json:"ticketCategory,omitempty"`   // Autotask
	TicketType       string   `json:"ticketType,omitempty"`       // Autotask
	Type             string   `json:"type,omitempty"`             // Microsoft Teams Bot
	Urgency          string   `json:"urgency,omitempty"`          // ServiceNow: 1 - High, 2 - Medium, 3 - Low (Default)
	WebhookURL       string   `json:"webhookUrl,omitempty"`       // Custom
	URL              string   `json:"url,omitempty"`              // DingTalk
}

// AlertActionParamsAutotask definition
type AlertActionParamsAutotask struct {
	CompanyID      string `json:"companyId,omitempty"`      // Autotask: Company ID
	IssueType      string `json:"issueType,omitempty"`      // Autotask: Issue type
	QueueID        int64  `json:"queueId,omitempty"`        // Autotask: Queue ID
	TicketCategory string `json:"ticketCategory,omitempty"` // Autotask ticket category
	TicketType     string `json:"ticketType,omitempty"`     // Autotask ticket type
}

// AlertActionParamsJira definition
type AlertActionParamsJira struct {
	Project      string `json:"project,omitempty"`
	IssueType    string `json:"issueType,omitempty"` // "Bug" | "Epic" | "Subtask" | "Story" | "Task"
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsMicrosoftTeams definition
type AlertActionParamsMicrosoftTeams struct{}

// AlertActionParamsMicrosoftTeamsBot definition
type AlertActionParamsMicrosoftTeamsBot struct {
	ChannelID   string `json:"channelId,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
	TeamID      string `json:"teamId,omitempty"`
	TeamName    string `json:"teamName,omitempty"`
	Type        string `json:"type,omitempty"` // "chat" | "meeting"
}

// AlertActionParamsMicrosoftTeamsWebhook definition
type AlertActionParamsMicrosoftTeamsWebhook struct {
	URL          string `json:"url,omitempty"`
	BodyTemplate string `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsSlackWebhook definition
type AlertActionParamsSlackWebhook struct {
	URL string `json:"url,omitempty"`
}

// AlertActionParamsServiceNow definition
type AlertActionParamsServiceNow struct {
	CallerID     string `json:"callerId,omitempty"` // user email
	Impact       string `json:"impact,omitempty"`   // 1 - High, 2 - Medium, 3 - Low (Default)
	Urgency      string `json:"urgency,omitempty"`  // 1 - High, 2 - Medium, 3 - Low (Default)
	BodyTemplate string `json:"bodyTemplate,omitempty"`
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

// AlertActionParamsEmail definition
type AlertActionParamsEmail struct {
	Recipients   []string `json:"recipients,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	BodyTemplate string   `json:"bodyTemplate,omitempty"`
}

// AlertActionParamsZammad definition
type AlertActionParamsZammad struct {
	Email string `json:"email,omitempty"`
}

// AlertActionParamsDingTalk definition
type AlertActionParamsDingTalk struct {
	IsAtAll   bool     `json:"isAtAll,omitempty"`
	AtMobiles []string `json:"atMobiles,omitempty"`
}

// AlertActionParamsDingTalkAction definition
type AlertActionParamsDingTalkAction struct {
	URL       string   `json:"url,omitempty"`
	Secret    string   `json:"secret,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
	AtMobiles []string `json:"atMobiles,omitempty"`
}

// AlertActionParamsAutomationRule definition
type AlertActionParamsAutomationRule struct {
	AlertType        string  `json:"alertType"`
	ResolveIncident  bool    `json:"resolveIncident,omitempty"`
	ServiceStatus    string  `json:"serviceStatus"`
	TemplateId       int64   `json:"templateId,omitempty"`
	SendNotification bool    `json:"sendNotification,omitempty"`
	ServiceIds       []int64 `json:"serviceIds"`
}

// AlertActionParamsTelegram definition
type AlertActionParamsTelegram struct {
	ChannelID string `json:"channelId,omitempty"`
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

// AlertFilter definition
type AlertFilter struct {
	Operator   string                 `json:"operator"`
	Predicates []AlertFilterPredicate `json:"predicates"`
}

// AlertFilterPredicate definition
type AlertFilterPredicate struct {
	Field    string `json:"field"`
	Criteria string `json:"criteria"`
	Value    string `json:"value"`
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
	AlertCreated          string
	AlertAssigned         string
	AlertAutoEscalated    string
	AlertAcknowledged     string
	AlertRaised           string
	AlertCommentAdded     string
	AlertEscalationEnded  string
	AlertResolved         string
	AlertAutoResolved     string
	AlertResponderAdded   string
	AlertResponderRemoved string
	AlertChannelAttached  string
	AlertChannelDetached  string
	AlertNotResolved      string
}{
	AlertCreated:          "alert-created",
	AlertAssigned:         "alert-assigned",
	AlertAutoEscalated:    "alert-auto-escalated",
	AlertAcknowledged:     "alert-acknowledged",
	AlertRaised:           "alert-raised",
	AlertCommentAdded:     "alert-comment-added",
	AlertEscalationEnded:  "alert-escalation-ended",
	AlertResolved:         "alert-resolved",
	AlertAutoResolved:     "alert-auto-resolved",
	AlertResponderAdded:   "alert-responder-added",
	AlertResponderRemoved: "alert-responder-removed",
	AlertChannelAttached:  "alert-channel-attached",
	AlertChannelDetached:  "alert-channel-detached",
	AlertNotResolved:      "v-alert-not-resolved",
}

// AlertActionTriggerTypesAll defines all alertAction trigger types
var AlertActionTriggerTypesAll = []string{
	AlertActionTriggerTypes.AlertCreated,
	AlertActionTriggerTypes.AlertAssigned,
	AlertActionTriggerTypes.AlertAutoEscalated,
	AlertActionTriggerTypes.AlertAcknowledged,
	AlertActionTriggerTypes.AlertRaised,
	AlertActionTriggerTypes.AlertCommentAdded,
	AlertActionTriggerTypes.AlertEscalationEnded,
	AlertActionTriggerTypes.AlertResolved,
	AlertActionTriggerTypes.AlertAutoResolved,
	AlertActionTriggerTypes.AlertResponderAdded,
	AlertActionTriggerTypes.AlertResponderRemoved,
	AlertActionTriggerTypes.AlertChannelAttached,
	AlertActionTriggerTypes.AlertChannelDetached,
	AlertActionTriggerTypes.AlertNotResolved,
}

// AlertFilterOperator defines alertFilter operator
var AlertFilterOperator = struct {
	And string
	Or  string
}{
	And: "AND",
	Or:  "OR",
}

// AlertFilterOperatorAll defines all alertFilter operator
var AlertFilterOperatorAll = []string{
	AlertFilterOperator.And,
	AlertFilterOperator.Or,
}

// AlertFilterPredicateFields defines alertFilter predicate fields
var AlertFilterPredicateFields = struct {
	AlertSummary     string
	AlertDetails     string
	EscalationPolicy string
	AlertPriority    string
}{
	AlertSummary:     "ALERT_SUMMARY",
	AlertDetails:     "ALERT_DETAILS",
	EscalationPolicy: "ESCALATION_POLICY",
	AlertPriority:    "ALERT_PRIORITY",
}

// AlertFilterPredicateFieldsAll defines all alertFilter predicate fields
var AlertFilterPredicateFieldsAll = []string{
	AlertFilterPredicateFields.AlertSummary,
	AlertFilterPredicateFields.AlertDetails,
	AlertFilterPredicateFields.EscalationPolicy,
	AlertFilterPredicateFields.AlertPriority,
}

// AlertFilterPredicateCriteria defines alertFilter predicate criteria
var AlertFilterPredicateCriteria = struct {
	ContainsAnyWords  string
	ContainsNotWords  string
	ContainsString    string
	ContainsNotString string
	IsString          string
	IsNotString       string
	MatchesRegex      string
	MatchesNotRegex   string
}{
	ContainsAnyWords:  "CONTAINS_ANY_WORDS",
	ContainsNotWords:  "CONTAINS_NOT_WORDS",
	ContainsString:    "CONTAINS_STRING",
	ContainsNotString: "CONTAINS_NOT_STRING",
	IsString:          "IS_STRING",
	IsNotString:       "IS_NOT_STRING",
	MatchesRegex:      "MATCHES_REGEX",
	MatchesNotRegex:   "MATCHES_NOT_REGEX",
}

// AlertFilterPredicateCriteriaAll defines all alertFilter predicate criteria
var AlertFilterPredicateCriteriaAll = []string{
	AlertFilterPredicateCriteria.ContainsAnyWords,
	AlertFilterPredicateCriteria.ContainsNotWords,
	AlertFilterPredicateCriteria.ContainsString,
	AlertFilterPredicateCriteria.ContainsNotString,
	AlertFilterPredicateCriteria.IsString,
	AlertFilterPredicateCriteria.IsNotString,
	AlertFilterPredicateCriteria.MatchesRegex,
	AlertFilterPredicateCriteria.MatchesNotRegex,
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

// CreateAlertAction creates a new alert action. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions/post
func (c *Client) CreateAlertAction(input *CreateAlertActionInput) (*CreateAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertAction == nil {
		return nil, errors.New("alert action input is required")
	}
	if input.AlertAction.AlertSources != nil && len(*input.AlertAction.AlertSources) == 1 && (input.AlertAction.Teams == nil || len(*input.AlertAction.Teams) == 0) && input.AlertAction.Conditions == "" {
		sourceId := (*input.AlertAction.AlertSources)[0].ID

		// manually set fields to ensure backwards compatibility with api v1
		input.AlertAction.AlertSourceIDs = []int64{sourceId}
		input.AlertAction.AlertSources = nil
		input.AlertAction.Teams = nil
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
	Version       *int
}

// GetAlertActionOutput represents the output of a GetAlertAction operation.
type GetAlertActionOutput struct {
	_           struct{}
	AlertAction *AlertActionOutput
}

// GetAlertAction gets the alert action with specified id. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions~1{id}/get
func (c *Client) GetAlertAction(input *GetAlertActionInput) (*GetAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertActionID == nil {
		return nil, errors.New("alert action id is required")
	}
	q := url.Values{}
	if input.Version != nil {
		q.Add("version", strconv.Itoa(*input.Version))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%s?%s", apiRoutes.alertActions, *input.AlertActionID, q.Encode()))
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

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	// Default: 0
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults *int
}

// GetAlertActionsOutput represents the output of a GetAlertActions operation.
type GetAlertActionsOutput struct {
	_            struct{}
	AlertActions []*AlertActionOutput
}

// GetAlertActions lists existing alert actions. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions/get
func (c *Client) GetAlertActions(input *GetAlertActionsInput) (*GetAlertActionsOutput, error) {
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

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.alertActions, q.Encode()))
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
	AlertAction *AlertActionOutput
}

// SearchAlertAction gets the alert action with specified name.
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

	alertAction := &AlertActionOutput{}
	err = json.Unmarshal(resp.Body(), alertAction)
	if err != nil {
		return nil, err
	}

	return &SearchAlertActionOutput{AlertAction: alertAction}, nil
}

// UpdateAlertActionInput represents the input of an UpdateAlertAction operation.
type UpdateAlertActionInput struct {
	_             struct{}
	AlertActionID *string
	AlertAction   *AlertAction
}

// UpdateAlertActionOutput represents the output of an UpdateAlertAction operation.
type UpdateAlertActionOutput struct {
	_           struct{}
	AlertAction *AlertActionOutput
}

// UpdateAlertAction updates an existing alert action. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions~1{id}/put
func (c *Client) UpdateAlertAction(input *UpdateAlertActionInput) (*UpdateAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertAction == nil {
		return nil, errors.New("alert action input is required")
	}
	if input.AlertActionID == nil {
		return nil, errors.New("alert action id is required")
	}
	if input.AlertAction.AlertSources != nil && len(*input.AlertAction.AlertSources) == 1 && (input.AlertAction.Teams == nil || len(*input.AlertAction.Teams) == 0) && input.AlertAction.Conditions == "" {
		sourceId := (*input.AlertAction.AlertSources)[0].ID

		// manually set fields to ensure backwards compatibility with api v1
		input.AlertAction.AlertSourceIDs = []int64{sourceId}
		input.AlertAction.AlertSources = nil
		input.AlertAction.Teams = nil
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

// DeleteAlertAction deletes the specified alert action. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alert-actions~1{id}/delete
func (c *Client) DeleteAlertAction(input *DeleteAlertActionInput) (*DeleteAlertActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertActionID == nil {
		return nil, errors.New("alert action id is required")
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
