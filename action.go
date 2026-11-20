package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Action definition, an early access resource behind /api/beta/actions. It may
// change or be removed while it is in early access, and its operations need a
// client opted in with WithEarlyAccess("actions").
//
// An action calls another system when an alert or incident changes, for example to
// open a Jira issue or post to a webhook. It is built from a template that ilert
// maintains, referenced by TemplateKey, or defines its own Steps. Actions are the
// successor to alert actions, which the API lists as read-only mirrors with the
// origin LEGACY_REF. Deleting a mirror deletes the alert action it mirrors.
type Action struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`

	// TemplateKey references the template the action is built from, e.g.
	// discord_alert. It cannot be changed or removed once set.
	TemplateKey string `json:"templateKey,omitempty"`
	// TemplateVersion pins the version of the template, nil follows the latest one
	TemplateVersion *int `json:"templateVersion,omitempty"`

	// EntityKind is ALERT or INCIDENT, see ActionEntityKind. It cannot be changed
	// once the action is created.
	EntityKind string `json:"entityKind,omitempty"`
	// DispatchMode is AUTOMATIC or MANUAL, see ActionDispatchMode. It cannot be
	// changed once the action is created.
	DispatchMode string `json:"dispatchMode,omitempty"`
	// TriggerEvents are the events the action reacts to, see ActionTriggerEvents.
	// None reacts to every event of its entity kind.
	TriggerEvents      []string                  `json:"triggerEvents,omitempty"`
	TriggerEventDelays []ActionTriggerEventDelay `json:"triggerEventDelays,omitempty"`
	// ConditionICL is the condition an alert or incident has to meet, in the ilert
	// condition language. Blank uses the condition of the template.
	ConditionICL string `json:"conditionICL,omitempty"`

	// Config holds the values for the fields the template declares, as a JSON
	// object. Some templates write values into it when the action is created or
	// updated, so send it back with them.
	Config json.RawMessage `json:"config,omitempty"`
	// Steps, FieldMapping and Inbound are the action's own definition, as JSON. An
	// action built from a template runs the template's while they are not set.
	Steps        json.RawMessage `json:"steps,omitempty"`
	FieldMapping json.RawMessage `json:"fieldMapping,omitempty"`
	Inbound      json.RawMessage `json:"inbound,omitempty"`

	// CredentialID references the credential the action authenticates with
	CredentialID *int64 `json:"credentialId,omitempty"`
	// Transport can only be set to HTTP, on an action without a template
	Transport     string `json:"transport,omitempty"`
	Bidirectional bool   `json:"bidirectional,omitempty"`

	Teams  []TeamShort       `json:"teams,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`

	// IntegrationKey and IntegrationURL authenticate and receive the callbacks of a
	// bidirectional action, and are only returned to users who may update it
	IntegrationKey string `json:"integrationKey,omitempty"`
	IntegrationURL string `json:"integrationUrl,omitempty"`
	// Brand is the vendor the action talks to, e.g. jira
	Brand string `json:"brand,omitempty"`
	// Status is the outcome of the most recent execution, see ActionStatus
	Status            string `json:"status,omitempty"`
	LastResultAt      int64  `json:"lastResultAt,omitempty"` // milliseconds since the epoch
	LastCorrelationID string `json:"lastCorrelationId,omitempty"`
	LastError         string `json:"lastError,omitempty"`
	LastErrorStepID   string `json:"lastErrorStepId,omitempty"`
	// Origin is NATIVE for an action and LEGACY_REF for the mirror of an alert
	// action, see ActionOrigin
	Origin              string `json:"origin,omitempty"`
	LegacyAlertActionID int64  `json:"legacyAlertActionId,omitempty"`
	CreatedAt           string `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt           string `json:"updatedAt,omitempty"` // date time string in ISO 8601
}

// ActionTriggerEventDelay delays the delivery of one trigger event, at most one per
// event. alert.not_resolved and incident.not_resolved take 60 to 172800 seconds, a
// shorter delay is raised to 60, and need one to be raised at all.
// alert.escalation_ended takes up to 7200 seconds, and a delay below 30 is dropped.
// Other events cannot be delayed.
type ActionTriggerEventDelay struct {
	Event    string `json:"event"`
	DelaySec int    `json:"delaySec"`
}

// ActionEntityKind defines the kinds of entity an action is bound to
var ActionEntityKind = struct {
	Alert    string
	Incident string
}{
	Alert:    "ALERT",
	Incident: "INCIDENT",
}

// ActionEntityKindAll defines all action entity kinds
var ActionEntityKindAll = []string{
	ActionEntityKind.Alert,
	ActionEntityKind.Incident,
}

// ActionDispatchMode defines when an action runs: AUTOMATIC on every matching
// event, MANUAL only when it is run on an alert or incident
var ActionDispatchMode = struct {
	Automatic string
	Manual    string
}{
	Automatic: "AUTOMATIC",
	Manual:    "MANUAL",
}

// ActionDispatchModeAll defines all action dispatch modes
var ActionDispatchModeAll = []string{
	ActionDispatchMode.Automatic,
	ActionDispatchMode.Manual,
}

// ActionTriggerEventsAll defines the events an action can react to. The alert.*
// events apply to alert actions, the incident.* events to incident actions.
var ActionTriggerEventsAll = []string{
	"alert.created",
	"alert.accepted",
	"alert.resolved",
	"alert.rerouted",
	"alert.escalated",
	"alert.escalation_ended",
	"alert.priority_raised",
	"alert.responder_added",
	"alert.responder_removed",
	"alert.comment_added",
	"alert.channel_attached",
	"alert.channel_detached",
	"alert.not_resolved",
	"incident.opened",
	"incident.updated",
	"incident.status_changed",
	"incident.resolved",
	"incident.deleted",
	"incident.responder_added",
	"incident.responder_removed",
	"incident.responder_updated",
	"incident.comment_added",
	"incident.channel_attached",
	"incident.channel_detached",
	"incident.postmortem_linked",
	"incident.not_resolved",
}

// ActionOrigin defines where an action comes from
var ActionOrigin = struct {
	Native    string
	LegacyRef string
}{
	Native:    "NATIVE",
	LegacyRef: "LEGACY_REF",
}

// ActionStatus defines the outcomes of the most recent execution of an action
var ActionStatus = struct {
	NeverTriggered string
	OK             string
	Error          string
}{
	NeverTriggered: "NEVER_TRIGGERED",
	OK:             "OK",
	Error:          "ERROR",
}

// CreateActionInput represents the input of a CreateAction operation.
type CreateActionInput struct {
	_      struct{}
	Action *Action
}

// CreateActionOutput represents the output of a CreateAction operation.
type CreateActionOutput struct {
	_      struct{}
	Action *Action
}

// CreateAction creates a new action resource. Name and EntityKind are required,
// together with either TemplateKey or Steps.
func (c *Client) CreateAction(input *CreateActionInput) (*CreateActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Action == nil {
		return nil, errors.New("action input is required")
	}

	req, err := c.betaRequest(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	resp, err := req.SetBody(input.Action).Post(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	action := &Action{}
	err = json.Unmarshal(resp.Body(), action)
	if err != nil {
		return nil, err
	}

	return &CreateActionOutput{Action: action}, nil
}

// GetActionInput represents the input of a GetAction operation.
type GetActionInput struct {
	_        struct{}
	ActionID *int64
}

// GetActionOutput represents the output of a GetAction operation.
type GetActionOutput struct {
	_      struct{}
	Action *Action
}

// GetAction gets the action resource with specified id.
func (c *Client) GetAction(input *GetActionInput) (*GetActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ActionID == nil {
		return nil, errors.New("action id is required")
	}

	req, err := c.betaRequest(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(fmt.Sprintf("%s/%d", apiRoutes.actions, *input.ActionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	action := &Action{}
	err = json.Unmarshal(resp.Body(), action)
	if err != nil {
		return nil, err
	}

	return &GetActionOutput{Action: action}, nil
}

// GetActionsInput represents the input of a GetActions operation.
type GetActionsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults *int
}

// GetActionsOutput represents the output of a GetActions operation.
type GetActionsOutput struct {
	_       struct{}
	Actions []*Action
}

// GetActions lists existing action resources, the mirrors of alert actions
// included, most recently updated first.
func (c *Client) GetActions(input *GetActionsInput) (*GetActionsOutput, error) {
	if input == nil {
		input = &GetActionsInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	req, err := c.betaRequest(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(fmt.Sprintf("%s?%s", apiRoutes.actions, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	actions := make([]*Action, 0)
	err = json.Unmarshal(resp.Body(), &actions)
	if err != nil {
		return nil, err
	}

	return &GetActionsOutput{Actions: actions}, nil
}

// SearchActionInput represents the input of a SearchAction operation.
type SearchActionInput struct {
	_          struct{}
	ActionName *string
}

// SearchActionOutput represents the output of a SearchAction operation.
type SearchActionOutput struct {
	_      struct{}
	Action *Action
}

// SearchAction gets the action resource with specified name. The mirrors of alert
// actions are not found by name.
func (c *Client) SearchAction(input *SearchActionInput) (*SearchActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ActionName == nil {
		return nil, errors.New("action name is required")
	}

	req, err := c.betaRequest(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(fmt.Sprintf("%s/name/%s", apiRoutes.actions, url.PathEscape(*input.ActionName)))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	action := &Action{}
	err = json.Unmarshal(resp.Body(), action)
	if err != nil {
		return nil, err
	}

	return &SearchActionOutput{Action: action}, nil
}

// UpdateActionInput represents the input of a UpdateAction operation.
type UpdateActionInput struct {
	_        struct{}
	ActionID *int64
	Action   *Action
}

// UpdateActionOutput represents the output of a UpdateAction operation.
type UpdateActionOutput struct {
	_      struct{}
	Action *Action
}

// actionClearedValues are what UpdateAction sends for the fields an action leaves
// empty, so the API clears them: it merges an update onto the action, and keeps
// the value of a field that is left out
var actionClearedValues = map[string]json.RawMessage{
	"description":        json.RawMessage(`""`),
	"templateVersion":    json.RawMessage(`null`),
	"triggerEvents":      json.RawMessage(`[]`),
	"triggerEventDelays": json.RawMessage(`[]`),
	"conditionICL":       json.RawMessage(`""`),
	"config":             json.RawMessage(`{}`),
	"steps":              json.RawMessage(`[]`),
	"fieldMapping":       json.RawMessage(`[]`),
	"inbound":            json.RawMessage(`{}`),
	"credentialId":       json.RawMessage(`null`),
	"bidirectional":      json.RawMessage(`false`),
	"teams":              json.RawMessage(`[]`),
	"labels":             json.RawMessage(`{}`),
}

// actionReadOnlyFields are left out of an update. The transport is too: a read
// often reports HTTP for an action whose template uses another one, and the API
// refuses a transport that differs from its template's.
var actionReadOnlyFields = []string{"id", "transport", "integrationKey", "integrationUrl", "brand", "status",
	"lastResultAt", "lastCorrelationId", "lastError", "lastErrorStepId", "origin", "legacyAlertActionId",
	"createdAt", "updatedAt"}

// UpdateAction updates an existing action resource. Like the other updates of this
// package it sends the whole action, so a field left empty is cleared: send it as
// GetAction returned it with the changes applied. The config included, since some
// templates write values into it. Enabled is kept while it is nil, and the template
// cannot be removed.
func (c *Client) UpdateAction(input *UpdateActionInput) (*UpdateActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Action == nil {
		return nil, errors.New("action input is required")
	}
	if input.ActionID == nil {
		return nil, errors.New("action id is required")
	}

	body, err := actionUpdateBody(input.Action)
	if err != nil {
		return nil, err
	}

	req, err := c.betaRequest(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	resp, err := req.SetBody(body).Put(fmt.Sprintf("%s/%d", apiRoutes.actions, *input.ActionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	action := &Action{}
	err = json.Unmarshal(resp.Body(), action)
	if err != nil {
		return nil, err
	}

	return &UpdateActionOutput{Action: action}, nil
}

// actionUpdateBody is the action as UpdateAction sends it: every field it leaves
// empty sent empty, and the read-only ones left out
func actionUpdateBody(action *Action) (map[string]json.RawMessage, error) {
	encoded, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	body := map[string]json.RawMessage{}
	if err := json.Unmarshal(encoded, &body); err != nil {
		return nil, err
	}
	for field, cleared := range actionClearedValues {
		if _, ok := body[field]; !ok {
			body[field] = cleared
		}
	}
	for _, field := range actionReadOnlyFields {
		delete(body, field)
	}
	return body, nil
}

// DeleteActionInput represents the input of a DeleteAction operation.
type DeleteActionInput struct {
	_        struct{}
	ActionID *int64
}

// DeleteActionOutput represents the output of a DeleteAction operation.
type DeleteActionOutput struct {
	_ struct{}
}

// DeleteAction deletes the specified action resource. Deleting the mirror of an
// alert action, an action with the origin LEGACY_REF, deletes the alert action.
func (c *Client) DeleteAction(input *DeleteActionInput) (*DeleteActionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ActionID == nil {
		return nil, errors.New("action id is required")
	}

	req, err := c.betaRequest(apiRoutes.actions)
	if err != nil {
		return nil, err
	}
	resp, err := req.Delete(fmt.Sprintf("%s/%d", apiRoutes.actions, *input.ActionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteActionOutput{}, nil
}
