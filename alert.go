package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Alert definition
type Alert struct {
	ID                 int64                  `json:"id"`
	Summary            string                 `json:"summary"`
	Details            string                 `json:"details"`
	ReportTime         string                 `json:"reportTime"` // Date time string in ISO format
	ResolvedOn         string                 `json:"resolvedOn"` // Date time string in ISO format
	Status             string                 `json:"status"`
	AlertSource        *AlertSource           `json:"alertSource,omitempty"`
	EscalationPolicy   *EscalationPolicy      `json:"scalationPolicy,omitempty"`
	Priority           string                 `json:"priority"`
	AlertKey           string                 `json:"alertKey"`
	AssignedTo         *User                  `json:"assignedTo,omitempty"`
	NextEscalation     string                 `json:"nextEscalation"` // Date time string in ISO format
	CallRoutingNumber  *CallRoutingNumber     `json:"callRoutingNumber,omitempty"`
	AcknowledgedBy     *User                  `json:"acknowledgedBy,omitempty"`
	AcknowledgedByType string                 `json:"acknowledgedByType,omitempty"`
	ResolvedBy         *User                  `json:"resolvedBy,omitempty"`
	ResolvedByType     string                 `json:"resolvedByType,omitempty"`
	Images             []AlertImage           `json:"images,omitempty"`
	Links              []AlertLink            `json:"links,omitempty"`
	CustomDetails      map[string]interface{} `json:"customDetails,omitempty"`
}

// AlertImage represents event image
type AlertImage struct {
	Src  string `json:"src"`
	Href string `json:"href"`
	Alt  string `json:"alt"`
}

// AlertLink represents event link
type AlertLink struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

// AlertComment definition
type AlertComment struct {
	ID             string `json:"id"`
	Content        string `json:"content"`
	Creator        *User  `json:"creator"`
	TriggerType    string `json:"triggerType"`
	ResolveComment bool   `json:"resolveComment"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
}

// CallRoutingNumber definition
type CallRoutingNumber struct {
	ID                  int          `json:"id"`
	Number              Phone        `json:"number"`
	VoiceLanguageLocale string       `json:"voiceLanguageLocale"`
	AlertSource         *AlertSource `json:"alertSource"`
}

// AlertStatuses defines alert statuses
var AlertStatuses = struct {
	New      string
	Pending  string
	Accepted string
	Resolved string
}{
	New:      "NEW",
	Pending:  "PENDING",
	Accepted: "ACCEPTED",
	Resolved: "RESOLVED",
}

// AlertPriorities defines alert priorities
var AlertPriorities = struct {
	High string
	Low  string
}{
	High: "HIGH",
	Low:  "LOW",
}

// AlertResponderTypes defines alert responder types
var AlertResponderTypes = struct {
	User        string
	AlertSource string
}{
	User:        "USER",
	AlertSource: "SOURCE",
}

// AlertResponder definition
type AlertResponder struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Disabled bool   `json:"disabled"`
}

// AlertResponderGroups defines alert responder groups
var AlertResponderGroups = struct {
	Suggested        string
	User             string
	EscalationPolicy string
	OnCallSchedule   string
}{
	Suggested:        "SUGGESTED",
	User:             "USER",
	EscalationPolicy: "ESCALATION_POLICY",
	OnCallSchedule:   "ON_CALL_SCHEDULE",
}

// AlertLogEntry definition
type AlertLogEntry struct {
	ID           int64  `json:"id"`
	Timestamp    string `json:"timestamp"` // Date time string in ISO format
	LogEntryType string `json:"logEntryType"`
	Text         string `json:"text"`
	AlertID      int64  `json:"alertId"`
}

// AlertLogEntryTypes defines alert log entry types
var AlertLogEntryTypes = struct {
	AlertReceivedLogEntry         string
	AlertSourceResponseLogEntry   string
	EmailReceivedLogEntry         string
	AlertAssignedBySystemLogEntry string
	AlertAssignedByUserLogEntry   string
	AlertCreatedByUserLogEntry    string
	NotificationLogEntry          string
	UserResponseLogEntry          string
}{
	AlertReceivedLogEntry:         "AlertReceivedLogEntry",
	AlertSourceResponseLogEntry:   "AlertSourceResponseLogEntry",
	EmailReceivedLogEntry:         "EmailReceivedLogEntry",
	AlertAssignedBySystemLogEntry: "AlertAssignedBySystemLogEntry",
	AlertAssignedByUserLogEntry:   "AlertAssignedByUserLogEntry",
	AlertCreatedByUserLogEntry:    "AlertCreatedByUserLogEntry",
	NotificationLogEntry:          "NotificationLogEntry",
	UserResponseLogEntry:          "UserResponseLogEntry",
}

var AlertInclude = struct {
	EscalationRules    string
	NextEscalationUser string
}{
	EscalationRules:    "escalationRules",
	NextEscalationUser: "nextEscalationUser",
}

// // AlertAction definition
// type AlertAction struct {
// 	Name        string              `json:"name"`
// 	WebhookID   string              `json:"webhookId"`
// 	ExtensionID string              `json:"extensionId"`
// 	IconURL     string              `json:"iconUrl"`
// 	History     []AlertActionResult `json:"history"`
// }

// // AlertActionResult definition
// type AlertActionResult struct {
// 	ID          string `json:"id"`
// 	AlertID     int64  `json:"alertId"`
// 	WebhookID   string `json:"webhookId"`
// 	ExtensionID string `json:"extensionId"`
// 	Actor       User   `json:"actor"`
// 	Success     bool   `json:"success"`
// }

// GetAlertInput represents the input of a GetAlert operation.
type GetAlertInput struct {
	_       struct{}
	AlertID *int64

	// describes optional properties that should be included in the response
	Include []*string
}

// GetAlertOutput represents the output of a GetAlert operation.
type GetAlertOutput struct {
	_     struct{}
	Alert *Alert
}

// GetAlert gets the alert with specified id. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}/get
func (c *Client) GetAlert(input *GetAlertInput) (*GetAlertOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertID == nil {
		return nil, errors.New("Alert id is required")
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.alerts, *input.AlertID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alert := &Alert{}
	err = json.Unmarshal(resp.Body(), alert)
	if err != nil {
		return nil, err
	}

	return &GetAlertOutput{Alert: alert}, nil
}

// GetAlertsInput represents the input of a GetAlerts operation.
type GetAlertsInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50
	MaxResults *int

	// state of the alert
	States []*string

	// alert source IDs of the alert's alert source
	AlertSources []*int64

	// user IDs of the user that the alert is assigned to
	AssignedToUserIDs []*int64

	// usernames of the user that the alert is assigned to
	AssignedToUserNames []*string

	// Date time string in ISO format
	From *string

	// Date time string in ISO format
	Until *string
}

// GetAlertsOutput represents the output of a GetAlerts operation.
type GetAlertsOutput struct {
	_      struct{}
	Alerts []*Alert
}

// GetAlerts lists alert sources. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts/get
func (c *Client) GetAlerts(input *GetAlertsInput) (*GetAlertsOutput, error) {
	if input == nil {
		input = &GetAlertsInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}
	if input.From != nil {
		q.Add("from", *input.From)
	}
	if input.Until != nil {
		q.Add("until", *input.Until)
	}

	for _, state := range input.States {
		q.Add("state", *state)
	}

	for _, alertSourceID := range input.AlertSources {
		q.Add("alert-source", strconv.FormatInt(*alertSourceID, 10))
	}

	for _, userID := range input.AssignedToUserIDs {
		q.Add("assigned-to", strconv.FormatInt(*userID, 10))
	}

	for _, username := range input.AssignedToUserNames {
		q.Add("assigned-to", *username)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.alerts, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alerts := make([]*Alert, 0)
	err = json.Unmarshal(resp.Body(), &alerts)
	if err != nil {
		return nil, err
	}

	return &GetAlertsOutput{Alerts: alerts}, nil
}

// GetAlertsCountInput represents the input of a GetAlertsCount operation.
type GetAlertsCountInput struct {
	_ struct{}

	// state of the alert
	States []*string

	// alert source IDs of the alert's alert source
	AlertSources []*int64

	// user IDs of the user that the alert is assigned to
	AssignedToUserIDs []*int64

	// usernames of the user that the alert is assigned to
	AssignedToUserNames []*string

	// Date time string in ISO format
	From *string

	// Date time string in ISO format
	Until *string
}

// GetAlertsCountOutput represents the output of a GetAlertsCount operation.
type GetAlertsCountOutput struct {
	_     struct{}
	Count int
}

// GetAlertsCount gets list uptime monitors. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1count/get
func (c *Client) GetAlertsCount(input *GetAlertsCountInput) (*GetAlertsCountOutput, error) {
	if input == nil {
		input = &GetAlertsCountInput{}
	}

	q := url.Values{}
	if input.From != nil {
		q.Add("from", *input.From)
	}
	if input.Until != nil {
		q.Add("until", *input.From)
	}

	for _, state := range input.States {
		q.Add("state", *state)
	}

	for _, alertSourceID := range input.AlertSources {
		q.Add("alert-source", strconv.FormatInt(*alertSourceID, 10))
	}

	for _, userID := range input.AssignedToUserIDs {
		q.Add("assigned-to", strconv.FormatInt(*userID, 10))
	}

	for _, username := range input.AssignedToUserNames {
		q.Add("assigned-to", *username)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/count?%s", apiRoutes.alerts, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	body := &GenericCountResponse{}
	err = json.Unmarshal(resp.Body(), body)
	if err != nil {
		return nil, err
	}

	return &GetAlertsCountOutput{Count: body.Count}, nil
}

// GetAlertResponderInput represents the input of a GetAlertResponder operation.
type GetAlertResponderInput struct {
	_        struct{}
	AlertID  *int64
	Language *string
}

// GetAlertResponderOutput represents the output of a GetAlertResponder operation.
type GetAlertResponderOutput struct {
	_          struct{}
	Responders []*AlertResponder
}

// GetAlertResponder gets the alert source with specified id. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}~1responder/get
func (c *Client) GetAlertResponder(input *GetAlertResponderInput) (*GetAlertResponderOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertID == nil {
		return nil, errors.New("Alert id is required")
	}

	q := url.Values{}
	if input.Language != nil {
		if *input.Language == "en" {
			q.Add("lng", "en")
		} else if *input.Language == "de" {
			q.Add("lng", "de")
		}
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/responder", apiRoutes.alerts, *input.AlertID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertResponders := make([]*AlertResponder, 0)
	err = json.Unmarshal(resp.Body(), &alertResponders)
	if err != nil {
		return nil, err
	}

	return &GetAlertResponderOutput{Responders: alertResponders}, nil
}

// AssignAlertInput represents the input of a AssignAlert operation.
type AssignAlertInput struct {
	_                  struct{}
	AlertID            *int64
	UserID             *int64
	Username           *string
	EscalationPolicyID *int64
	ScheduleID         *int64
}

// AssignAlertOutput represents the output of a AssignAlert operation.
type AssignAlertOutput struct {
	_     struct{}
	Alert *Alert
}

// AssignAlert gets the alert source with specified id. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}~1assign/put
func (c *Client) AssignAlert(input *AssignAlertInput) (*AssignAlertOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertID == nil {
		return nil, errors.New("Alert id is required")
	}

	if input.UserID == nil && input.Username == nil && input.EscalationPolicyID == nil && input.ScheduleID == nil {
		return nil, errors.New("one of assignments is required")
	}

	q := url.Values{}
	if input.UserID != nil {
		q.Add("user-id", strconv.FormatInt(*input.UserID, 10))
	}
	if input.Username != nil {
		q.Add("user-id", *input.Username)
	}
	if input.EscalationPolicyID != nil {
		q.Add("policy-id", strconv.FormatInt(*input.EscalationPolicyID, 10))
	}
	if input.ScheduleID != nil {
		q.Add("schedule-id", strconv.FormatInt(*input.ScheduleID, 10))
	}

	resp, err := c.httpClient.R().Put(fmt.Sprintf("%s/%d/assign?%s", apiRoutes.alerts, *input.AlertID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alert := &Alert{}
	err = json.Unmarshal(resp.Body(), alert)
	if err != nil {
		return nil, err
	}

	return &AssignAlertOutput{Alert: alert}, nil
}

// AcceptAlertInput represents the input of a AcceptAlert operation.
type AcceptAlertInput struct {
	_       struct{}
	AlertID *int64
}

// AcceptAlertOutput represents the output of a AcceptAlert operation.
type AcceptAlertOutput struct {
	_     struct{}
	Alert *Alert
}

// AcceptAlert gets the alert source with specified id. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}~1accept/put
func (c *Client) AcceptAlert(input *AcceptAlertInput) (*AcceptAlertOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertID == nil {
		return nil, errors.New("Alert id is required")
	}

	resp, err := c.httpClient.R().Put(fmt.Sprintf("%s/%d/accept", apiRoutes.alerts, *input.AlertID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alert := &Alert{}
	err = json.Unmarshal(resp.Body(), alert)
	if err != nil {
		return nil, err
	}

	return &AcceptAlertOutput{Alert: alert}, nil
}

// ResolveAlertInput represents the input of a ResolveAlert operation.
type ResolveAlertInput struct {
	_       struct{}
	AlertID *int64
}

// ResolveAlertOutput represents the output of a ResolveAlert operation.
type ResolveAlertOutput struct {
	_     struct{}
	Alert *Alert
}

// ResolveAlert gets the alert source with specified id. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}~1resolve/put
func (c *Client) ResolveAlert(input *ResolveAlertInput) (*ResolveAlertOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertID == nil {
		return nil, errors.New("Alert id is required")
	}

	resp, err := c.httpClient.R().Put(fmt.Sprintf("%s/%d/resolve", apiRoutes.alerts, *input.AlertID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alert := &Alert{}
	err = json.Unmarshal(resp.Body(), alert)
	if err != nil {
		return nil, err
	}

	return &ResolveAlertOutput{Alert: alert}, nil
}

// GetAlertLogEntriesInput represents the input of a GetAlertLogEntries operation.
type GetAlertLogEntriesInput struct {
	_        struct{}
	AlertID  *int64
	Language *string
}

// GetAlertLogEntriesOutput represents the output of a GetAlertLogEntries operation.
type GetAlertLogEntriesOutput struct {
	_          struct{}
	LogEntries []*AlertLogEntry
}

// GetAlertLogEntries gets log entries for the specified alert. https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}~1log-entries/get
func (c *Client) GetAlertLogEntries(input *GetAlertLogEntriesInput) (*GetAlertLogEntriesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertID == nil {
		return nil, errors.New("Alert id is required")
	}

	q := url.Values{}
	if input.Language != nil {
		if *input.Language == "en" {
			q.Add("lng", "en")
		} else if *input.Language == "de" {
			q.Add("lng", "de")
		}
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/log-entries", apiRoutes.alerts, *input.AlertID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertLogEntries := make([]*AlertLogEntry, 0)
	err = json.Unmarshal(resp.Body(), &alertLogEntries)
	if err != nil {
		return nil, err
	}

	return &GetAlertLogEntriesOutput{LogEntries: alertLogEntries}, nil
}

// TODO https://api.ilert.com/api-docs/#tag/Alerts/paths/~1alerts~1{id}~1notifications/get

// TODO https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alerts~1{id}~1actions/get
// // GetAlertActionsInput represents the input of a GetAlertsAction operation.
// type GetAlertActionsInput struct {
// 	_       struct{}
// 	AlertID *int64
// }

// // GetAlertActionsOutput represents the output of a GetAlertsAction operation.
// type GetAlertActionsOutput struct {
// 	_       struct{}
// 	Actions []*AlertAction
// }

// // GetAlertActions gets available actions for specified alert. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alerts~1{id}~1actions/get
// func (c *Client) GetAlertActions(input *GetAlertActionsInput) (*GetAlertActionsOutput, error) {
// 	if input == nil {
// 		return nil, errors.New("input is required")
// 	}
// 	if input.AlertID == nil {
// 		return nil, errors.New("Alert id is required")
// 	}

// 	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/actions", apiRoutes.alerts, *input.AlertID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
// 		return nil, apiErr
// 	}

// 	alertActions := make([]*AlertAction, 0)
// 	err = json.Unmarshal(resp.Body(), &alertActions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GetAlertActionsOutput{Actions: alertActions}, nil
// }

// TODO https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alerts~1{id}~1actions/post
// // InvokeAlertActionInput represents the input of a InvokeAlertAction operation.
// type InvokeAlertActionInput struct {
// 	_       struct{}
// 	AlertID *int64
// 	Action  *AlertAction
// }

// // InvokeAlertActionOutput represents the output of a InvokeAlertAction operation.
// type InvokeAlertActionOutput struct {
// 	_      struct{}
// 	Action *AlertAction
// }

// // InvokeAlertAction invokes a specific alert action. https://api.ilert.com/api-docs/#tag/Alert-Actions/paths/~1alerts~1{id}~1actions/post
// func (c *Client) InvokeAlertAction(input *InvokeAlertActionInput) (*InvokeAlertActionOutput, error) {
// 	if input == nil {
// 		return nil, errors.New("input is required")
// 	}
// 	if input.AlertID == nil {
// 		return nil, errors.New("Alert id is required")
// 	}
// 	if input.Action == nil {
// 		return nil, errors.New("action input is required")
// 	}

// 	resp, err := c.httpClient.R().SetBody(input.Action).Post(fmt.Sprintf("%s/%d/actions", apiRoutes.alerts, *input.AlertID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
// 		return nil, apiErr
// 	}

// 	alertAction := &AlertAction{}
// 	err = json.Unmarshal(resp.Body(), alertAction)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &InvokeAlertActionOutput{Action: alertAction}, nil
// }
