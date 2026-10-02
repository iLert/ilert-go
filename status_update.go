package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// StatusUpdate definition https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
//
// A status update communicates an outage of one or more services to status pages and subscribers.
// Up to version 3 of this SDK it was called Incident and read from /api/incidents, which the API
// now reserves for the operational Incident.
type StatusUpdate struct {
	ID               int64              `json:"id"`
	Summary          string             `json:"summary"`
	Status           string             `json:"status"`
	Message          string             `json:"message"`
	SendNotification bool               `json:"sendNotification"`
	CreatedAt        string             `json:"createdAt"` // Date time string in ISO format
	UpdatedAt        string             `json:"updatedAt"` // Date time string in ISO format
	AffectedServices []AffectedServices `json:"affectedServices"`
	ResolvedOn       string             `json:"resolvedOn,omitempty"` // Date time string in ISO format
	Subscribed       bool               `json:"subscribed,omitempty"`
	AffectedTeams    []TeamShort        `json:"affectedTeams,omitempty"`

	// the messages published on the status update, newest first. Only part of a response when
	// StatusUpdateInclude.History is requested from GetStatusUpdate.
	History []StatusUpdateHistoryEntry `json:"history,omitempty"`
}

// StatusUpdateHistoryEntry is one message published on a status update
type StatusUpdateHistoryEntry struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
	Creator *User  `json:"creator,omitempty"`

	// the status the status update had when the message was published
	// possible values: "INVESTIGATING", "IDENTIFIED", "MONITORING", "RESOLVED"
	Status string `json:"incidentStatus,omitempty"`

	SendNotification bool   `json:"sendNotification,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"` // Date time string in ISO format
}

// AffectedServices defines affected services
type AffectedServices struct {
	Impact  string  `json:"impact"`
	Service Service `json:"service"`
}

// Subscriber defines a subscriber
type Subscriber struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// UIMenuItem defines the ui menu item
type UIMenuItem struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Affected defines the status pages and subscribers a status update or an incident would notify
type Affected struct {
	// one entry per status page that would show it
	StatusPagesInfo    []UIMenuItem `json:"statusPagesInfo"`
	PrivateStatusPages int64        `json:"privateStatusPages"`
	PublicStatusPages  int64        `json:"publicStatusPages"`
	PrivateSubscribers int64        `json:"privateSubscribers"`
	PublicSubscribers  int64        `json:"publicSubscribers"`
}

// SubscriberType defines subscriber type
var SubscriberType = struct {
	User string
	Team string
}{
	User: "USER",
	Team: "TEAM",
}

// SubscriberTypeAll defines subscriber type list
var SubscriberTypeAll = []string{
	SubscriberType.User,
	SubscriberType.Team,
}

// StatusUpdateInclude defines status update includes
var StatusUpdateInclude = struct {
	Subscribed    string
	AffectedTeams string
	History       string
}{
	Subscribed:    "subscribed",
	AffectedTeams: "affectedTeams",
	History:       "history",
}

// StatusUpdateIncludeAll defines status update includes list
var StatusUpdateIncludeAll = []string{
	StatusUpdateInclude.Subscribed,
	StatusUpdateInclude.AffectedTeams,
	StatusUpdateInclude.History,
}

// StatusUpdateStatus defines status update status
var StatusUpdateStatus = struct {
	Investigating string
	Identified    string
	Monitoring    string
	Resolved      string
}{
	Investigating: "INVESTIGATING",
	Identified:    "IDENTIFIED",
	Monitoring:    "MONITORING",
	Resolved:      "RESOLVED",
}

// StatusUpdateStatusAll defines status update status list
var StatusUpdateStatusAll = []string{
	StatusUpdateStatus.Investigating,
	StatusUpdateStatus.Identified,
	StatusUpdateStatus.Monitoring,
	StatusUpdateStatus.Resolved,
}

// StatusUpdateLogEntry is one entry on the timeline of a status update
type StatusUpdateLogEntry struct {
	ID        string `json:"id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"` // Date time string in ISO format
	Type      string `json:"type,omitempty"`
	SubType   string `json:"subType,omitempty"`

	// what happened, its keys depend on the sub type
	EventContext map[string]any `json:"eventContext,omitempty"`

	// what the entry changed, when it changed anything
	Diff map[string]any `json:"diff,omitempty"`

	ActorID   string `json:"actorId,omitempty"`
	ActorType string `json:"actorType,omitempty"`
}

// CreateStatusUpdateInput represents the input of a CreateStatusUpdate operation.
type CreateStatusUpdateInput struct {
	_            struct{}
	StatusUpdate *StatusUpdate
}

// CreateStatusUpdateOutput represents the output of a CreateStatusUpdate operation.
type CreateStatusUpdateOutput struct {
	_            struct{}
	StatusUpdate *StatusUpdate
}

// CreateStatusUpdate creates a new status update. It has to affect at least one service, the API
// rejects it with a 400 otherwise. Depending on the affected services this publishes notifications
// to subscribers, GetStatusUpdateAffected forecasts them. To publish a status update for an
// operational incident use CreateIncidentStatusUpdate instead. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) CreateStatusUpdate(input *CreateStatusUpdateInput) (*CreateStatusUpdateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdate == nil {
		return nil, errors.New("status update input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.StatusUpdate).Post(apiRoutes.statusUpdates)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200, 201); apiErr != nil {
		return nil, apiErr
	}

	statusUpdate := &StatusUpdate{}
	err = json.Unmarshal(resp.Body(), statusUpdate)
	if err != nil {
		return nil, err
	}

	return &CreateStatusUpdateOutput{StatusUpdate: statusUpdate}, nil
}

// GetStatusUpdatesInput represents the input of a GetStatusUpdates operation.
type GetStatusUpdatesInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	// Default: 0
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 10, Maximum: 25 or 100 without include
	MaxResults *int

	// describes optional properties that should be included in the response
	// possible values: "subscribed"
	Include []*string

	// status of the status update
	// possible values: "INVESTIGATING", "IDENTIFIED", "MONITORING", "RESOLVED"
	States []*string

	// service IDs of the status update's affected services
	Services []*int64

	// Date time string in ISO format, based on the creation of the status update
	From *string

	// Date time string in ISO format, based on the creation of the status update
	Until *string
}

// GetStatusUpdatesOutput represents the output of a GetStatusUpdates operation.
type GetStatusUpdatesOutput struct {
	_             struct{}
	StatusUpdates []*StatusUpdate
}

// GetStatusUpdates lists existing status updates. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) GetStatusUpdates(input *GetStatusUpdatesInput) (*GetStatusUpdatesOutput, error) {
	if input == nil {
		input = &GetStatusUpdatesInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	} else {
		q.Add("start-index", "0")
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	} else {
		q.Add("max-results", "10")
	}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	for _, state := range input.States {
		q.Add("states", *state)
	}

	for _, services := range input.Services {
		q.Add("services", strconv.FormatInt(*services, 10))
	}

	if input.From != nil {
		q.Add("from", *input.From)
	}
	if input.Until != nil {
		q.Add("until", *input.Until)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.statusUpdates, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusUpdates := make([]*StatusUpdate, 0)
	err = json.Unmarshal(resp.Body(), &statusUpdates)
	if err != nil {
		return nil, err
	}

	return &GetStatusUpdatesOutput{StatusUpdates: statusUpdates}, nil
}

// GetStatusUpdateInput represents the input of a GetStatusUpdate operation.
type GetStatusUpdateInput struct {
	_              struct{}
	StatusUpdateID *int64

	// describes optional properties that should be included in the response
	// possible values: "subscribed", "affectedTeams", "history"
	Include []*string
}

// GetStatusUpdateOutput represents the output of a GetStatusUpdate operation.
type GetStatusUpdateOutput struct {
	_            struct{}
	StatusUpdate *StatusUpdate

	// send it back as UpdateStatusUpdateInput.ETag to have the update rejected with a 412 when the
	// status update or its affected services changed in the meantime
	ETag *string
}

// GetStatusUpdate gets a status update by id. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) GetStatusUpdate(input *GetStatusUpdateInput) (*GetStatusUpdateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdateID == nil {
		return nil, errors.New("status update id is required")
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	var url = fmt.Sprintf("%s/%d?%s", apiRoutes.statusUpdates, *input.StatusUpdateID, q.Encode())

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusUpdate := &StatusUpdate{}
	err = json.Unmarshal(resp.Body(), statusUpdate)
	if err != nil {
		return nil, err
	}

	output := &GetStatusUpdateOutput{StatusUpdate: statusUpdate}
	etag := resp.Header().Get("ETag")
	if etag != "" {
		output.ETag = String(etag)
	}

	return output, nil
}

// GetStatusUpdateSubscribersInput represents the input of a GetStatusUpdateSubscribers operation.
type GetStatusUpdateSubscribersInput struct {
	_              struct{}
	StatusUpdateID *int64
}

// GetStatusUpdateSubscribersOutput represents the output of a GetStatusUpdateSubscribers operation.
type GetStatusUpdateSubscribersOutput struct {
	_           struct{}
	Subscribers []*Subscriber
}

// GetStatusUpdateSubscribers gets subscribers of a status update by id. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) GetStatusUpdateSubscribers(input *GetStatusUpdateSubscribersInput) (*GetStatusUpdateSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdateID == nil {
		return nil, errors.New("status update id is required")
	}

	var url = fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.statusUpdates, *input.StatusUpdateID)

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	subscribers := make([]*Subscriber, 0)
	err = json.Unmarshal(resp.Body(), &subscribers)
	if err != nil {
		return nil, err
	}

	return &GetStatusUpdateSubscribersOutput{Subscribers: subscribers}, nil
}

// GetStatusUpdateAffectedInput represents the input of a GetStatusUpdateAffected operation.
type GetStatusUpdateAffectedInput struct {
	_            struct{}
	StatusUpdate *StatusUpdate
}

// GetStatusUpdateAffectedOutput represents the output of a GetStatusUpdateAffected operation.
type GetStatusUpdateAffectedOutput struct {
	_        struct{}
	Affected *Affected
}

// GetStatusUpdateAffected forecasts the subscribers and status pages that creating or updating the
// status update would notify. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) GetStatusUpdateAffected(input *GetStatusUpdateAffectedInput) (*GetStatusUpdateAffectedOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdate == nil {
		return nil, errors.New("status update is required")
	}

	url := fmt.Sprintf("%s/publish-info", apiRoutes.statusUpdates)

	resp, err := c.httpClient.R().SetBody(input.StatusUpdate).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	affected := &Affected{}
	err = json.Unmarshal(resp.Body(), affected)
	if err != nil {
		return nil, err
	}

	return &GetStatusUpdateAffectedOutput{Affected: affected}, nil
}

// AddStatusUpdateSubscribersInput represents the input of a AddStatusUpdateSubscribers operation.
type AddStatusUpdateSubscribersInput struct {
	_              struct{}
	StatusUpdateID *int64
	Subscribers    *[]Subscriber
}

// AddStatusUpdateSubscribersOutput represents the output of a AddStatusUpdateSubscribers operation.
type AddStatusUpdateSubscribersOutput struct {
	_ struct{}
}

// AddStatusUpdateSubscribers adds subscribers to a status update. Adding one that is already
// subscribed is rejected with a 400. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) AddStatusUpdateSubscribers(input *AddStatusUpdateSubscribersInput) (*AddStatusUpdateSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdateID == nil {
		return nil, errors.New("status update id is required")
	}
	if input.Subscribers == nil {
		return nil, errors.New("subscriber input is required")
	}

	url := fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.statusUpdates, *input.StatusUpdateID)

	resp, err := c.httpClient.R().SetBody(input.Subscribers).Post(url)
	if err != nil {
		return nil, err
	}
	// the API answers 202 without a body, the spec documents 204
	if apiErr := getGenericAPIError(resp, 202, 204); apiErr != nil {
		return nil, apiErr
	}

	return &AddStatusUpdateSubscribersOutput{}, nil
}

// UpdateStatusUpdateInput represents the input of a UpdateStatusUpdate operation.
type UpdateStatusUpdateInput struct {
	_              struct{}
	StatusUpdateID *int64
	StatusUpdate   *StatusUpdate

	// the ETag returned by GetStatusUpdate, sent as If-Match
	ETag *string
}

// UpdateStatusUpdateOutput represents the output of a UpdateStatusUpdate operation.
type UpdateStatusUpdateOutput struct {
	_            struct{}
	StatusUpdate *StatusUpdate
}

// UpdateStatusUpdate updates the specific status update. The update is appended to the history of
// the status update and publishes notifications to subscribers. https://docs.ilert.com/developer-docs/rest-api/api-reference/status-updates
func (c *Client) UpdateStatusUpdate(input *UpdateStatusUpdateInput) (*UpdateStatusUpdateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdateID == nil {
		return nil, errors.New("status update id is required")
	}
	if input.StatusUpdate == nil {
		return nil, errors.New("status update input is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.statusUpdates, *input.StatusUpdateID)

	req := c.httpClient.R()
	if input.ETag != nil && *input.ETag != "" {
		req.SetHeader("If-Match", *input.ETag)
	}
	resp, err := req.SetBody(input.StatusUpdate).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusUpdate := &StatusUpdate{}
	err = json.Unmarshal(resp.Body(), statusUpdate)
	if err != nil {
		return nil, err
	}

	return &UpdateStatusUpdateOutput{StatusUpdate: statusUpdate}, nil
}

// GetStatusUpdateLogEntriesInput represents the input of a GetStatusUpdateLogEntries operation.
type GetStatusUpdateLogEntriesInput struct {
	_              struct{}
	StatusUpdateID *int64

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities.
	// The API currently skips entries on every page after the first, so prefer a MaxResults that covers
	// the whole timeline, or narrow it down with From and Until.
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 100
	MaxResults *int

	// start of the time window, has to be set together with Until. Date time string in ISO format.
	From *string

	// end of the time window, has to be set together with From. Date time string in ISO format.
	Until *string

	// restricts the result to these sub types
	SubTypes []*string
}

// GetStatusUpdateLogEntriesOutput represents the output of a GetStatusUpdateLogEntries operation.
type GetStatusUpdateLogEntriesOutput struct {
	_          struct{}
	LogEntries []*StatusUpdateLogEntry
}

// GetStatusUpdateLogEntries lists the timeline of a status update, oldest first. The endpoint throttles bursts of
// requests and answers them with a 403 "authentication failed" for about a minute, so a 403 here does
// not necessarily mean the token is wrong. https://api.ilert.com/api/beta/openapi.json
func (c *Client) GetStatusUpdateLogEntries(input *GetStatusUpdateLogEntriesInput) (*GetStatusUpdateLogEntriesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusUpdateID == nil {
		return nil, errors.New("status update id is required")
	}
	q, err := logEntriesQuery(input.StartIndex, input.MaxResults, input.From, input.Until, input.SubTypes)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/log-entries?%s", apiRoutes.statusUpdates, *input.StatusUpdateID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	logEntries := make([]*StatusUpdateLogEntry, 0)
	err = json.Unmarshal(resp.Body(), &logEntries)
	if err != nil {
		return nil, err
	}

	return &GetStatusUpdateLogEntriesOutput{LogEntries: logEntries}, nil
}

// logEntriesQuery builds the paging and filter parameters the incident and status update timelines share
func logEntriesQuery(startIndex *int, maxResults *int, from *string, until *string, subTypes []*string) (url.Values, error) {
	if (from == nil) != (until == nil) {
		return nil, errors.New("from and until have to be set together")
	}

	q := url.Values{}
	if startIndex != nil {
		q.Add("start-index", strconv.Itoa(*startIndex))
	}
	if maxResults != nil {
		q.Add("max-results", strconv.Itoa(*maxResults))
	}
	if from != nil {
		q.Add("from", *from)
		q.Add("until", *until)
	}
	for _, subType := range subTypes {
		q.Add("sub-types", *subType)
	}
	return q, nil
}
