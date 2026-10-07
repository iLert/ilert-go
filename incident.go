package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Incident definition https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
//
// An incident is the operational incident a team declares and works on, with its own status,
// severity, responders and chat channel. What it communicates to status pages and subscribers is
// published separately as a StatusUpdate, see CreateIncidentStatusUpdate.
//
// The API decides by the API contract version of the token which entity /api/incidents serves:
// version 2 and later get incidents, a token still on version 1 gets status updates on the same
// paths. Those do not decode into this type, so every Incident operation needs a token on version 2
// or later. The StatusUpdate operations use /api/status-updates and work with either.
//
// Incidents cannot be deleted through the API, it answers with a 400 and asks to resolve them instead.
type Incident struct {
	ID int64 `json:"id,omitempty"`

	// sequential number of the incident within the account
	Number int64 `json:"number,omitempty"`

	// required when creating an incident
	Title string `json:"title,omitempty"`

	Summary string `json:"summary,omitempty"`

	// possible values: "DECLARED", "INVESTIGATING", "IDENTIFIED", "MONITORING", "RESOLVED"
	Status string `json:"status,omitempty"`

	// required when creating an incident, in range 1..5. The published spec documents it as
	// optional, but the API rejects a create without it.
	Severity int `json:"severity,omitempty"`

	DeclaredAt string `json:"declaredAt,omitempty"` // Date time string in ISO format
	ResolvedAt string `json:"resolvedAt,omitempty"` // Date time string in ISO format
	DeclaredBy *User  `json:"declaredBy,omitempty"`
	ResolvedBy *User  `json:"resolvedBy,omitempty"`

	// the chat channel the incident is worked on in, read-only through these operations
	Channel *IncidentChannel `json:"channel,omitempty"`

	// the conference bridge of the incident, read-only through these operations
	ConferenceBridge *IncidentConferenceBridge `json:"conferenceBridge,omitempty"`

	// the services the incident affects. Deliberately sent as null when nil: the API leaves the
	// affected services untouched on null and replaces them with whatever list it is given, so a
	// non-nil empty slice removes every one of them.
	AffectedServices []AffectedServices `json:"affectedServices"`

	// read-only through these operations
	Responders []IncidentResponder `json:"responders,omitempty"`

	// only part of a response when IncidentInclude.AlertCount is requested
	AlertCount int `json:"alertCount,omitempty"`

	// only part of a response when IncidentInclude.ResponderCountJoined is requested
	ResponderCountJoined int `json:"responderCountJoined,omitempty"`

	// only part of a response when IncidentInclude.ResponderCountTotal is requested
	ResponderCountTotal int `json:"responderCountTotal,omitempty"`
}

// IncidentChannel defines the chat channel of an incident
type IncidentChannel struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	// possible values: "SLACK", "MS_TEAMS", "GOOGLE_CHAT"
	Type string `json:"type,omitempty"`

	Href      string `json:"href,omitempty"`
	GroupID   string `json:"groupId,omitempty"`
	IsPrivate bool   `json:"isPrivate,omitempty"`
}

// IncidentConferenceBridge defines the conference bridge of an incident
type IncidentConferenceBridge struct {
	// possible values: "ZOOM", "GOOGLE_MEET", "MS_TEAMS", "CUSTOM"
	Type string `json:"type,omitempty"`

	Href   string `json:"href,omitempty"`
	ID     string `json:"id,omitempty"`
	DialIn string `json:"dialIn,omitempty"`
}

// IncidentResponder defines a responder of an incident
type IncidentResponder struct {
	User *User `json:"user,omitempty"`

	// possible values: "ADDED", "NOTIFIED", "JOINED", "DECLINED"
	Status string `json:"status,omitempty"`

	// possible values: "INCIDENT_COMMANDER", "COMMUNICATIONS_LEAD", "ACCOUNTABLE_EXECUTIVE", "LEGAL_LEAD",
	// "CUSTOMER_SUCCESS_LEAD", "SUPPORT_LEAD", "FINANCE_LEAD", "RESPONDER"
	Role string `json:"role,omitempty"`

	NotifiedAt    string `json:"notifiedAt,omitempty"` // Date time string in ISO format
	AcceptedAt    string `json:"acceptedAt,omitempty"` // Date time string in ISO format
	DeclinedAt    string `json:"declinedAt,omitempty"` // Date time string in ISO format
	PagedByUserID int64  `json:"pagedByUserId,omitempty"`
}

// IncidentLogEntry is one entry on the timeline of an incident
type IncidentLogEntry struct {
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

// IncidentInclude defines incident includes
var IncidentInclude = struct {
	AffectedServices     string
	Responders           string
	AlertCount           string
	ResponderCountJoined string
	ResponderCountTotal  string
}{
	AffectedServices:     "affectedServices",
	Responders:           "responders",
	AlertCount:           "alertCount",
	ResponderCountJoined: "responderCountJoined",
	ResponderCountTotal:  "responderCountTotal",
}

// IncidentIncludeAll defines incident includes list
var IncidentIncludeAll = []string{
	IncidentInclude.AffectedServices,
	IncidentInclude.Responders,
	IncidentInclude.AlertCount,
	IncidentInclude.ResponderCountJoined,
	IncidentInclude.ResponderCountTotal,
}

// IncidentStatus defines incident status
var IncidentStatus = struct {
	Declared      string
	Investigating string
	Identified    string
	Monitoring    string
	Resolved      string
}{
	Declared:      "DECLARED",
	Investigating: "INVESTIGATING",
	Identified:    "IDENTIFIED",
	Monitoring:    "MONITORING",
	Resolved:      "RESOLVED",
}

// IncidentStatusAll defines incident status list
var IncidentStatusAll = []string{
	IncidentStatus.Declared,
	IncidentStatus.Investigating,
	IncidentStatus.Identified,
	IncidentStatus.Monitoring,
	IncidentStatus.Resolved,
}

// IncidentResponderStatus defines incident responder status
var IncidentResponderStatus = struct {
	Added    string
	Notified string
	Joined   string
	Declined string
}{
	Added:    "ADDED",
	Notified: "NOTIFIED",
	Joined:   "JOINED",
	Declined: "DECLINED",
}

// IncidentResponderStatusAll defines incident responder status list
var IncidentResponderStatusAll = []string{
	IncidentResponderStatus.Added,
	IncidentResponderStatus.Notified,
	IncidentResponderStatus.Joined,
	IncidentResponderStatus.Declined,
}

// IncidentResponderRole defines incident responder role
var IncidentResponderRole = struct {
	IncidentCommander    string
	CommunicationsLead   string
	AccountableExecutive string
	LegalLead            string
	CustomerSuccessLead  string
	SupportLead          string
	FinanceLead          string
	Responder            string
}{
	IncidentCommander:    "INCIDENT_COMMANDER",
	CommunicationsLead:   "COMMUNICATIONS_LEAD",
	AccountableExecutive: "ACCOUNTABLE_EXECUTIVE",
	LegalLead:            "LEGAL_LEAD",
	CustomerSuccessLead:  "CUSTOMER_SUCCESS_LEAD",
	SupportLead:          "SUPPORT_LEAD",
	FinanceLead:          "FINANCE_LEAD",
	Responder:            "RESPONDER",
}

// IncidentResponderRoleAll defines incident responder role list
var IncidentResponderRoleAll = []string{
	IncidentResponderRole.IncidentCommander,
	IncidentResponderRole.CommunicationsLead,
	IncidentResponderRole.AccountableExecutive,
	IncidentResponderRole.LegalLead,
	IncidentResponderRole.CustomerSuccessLead,
	IncidentResponderRole.SupportLead,
	IncidentResponderRole.FinanceLead,
	IncidentResponderRole.Responder,
}

// IncidentChannelType defines incident channel type
var IncidentChannelType = struct {
	Slack      string
	MsTeams    string
	GoogleChat string
}{
	Slack:      "SLACK",
	MsTeams:    "MS_TEAMS",
	GoogleChat: "GOOGLE_CHAT",
}

// IncidentChannelTypeAll defines incident channel type list
var IncidentChannelTypeAll = []string{
	IncidentChannelType.Slack,
	IncidentChannelType.MsTeams,
	IncidentChannelType.GoogleChat,
}

// IncidentConferenceBridgeType defines incident conference bridge type
var IncidentConferenceBridgeType = struct {
	Zoom       string
	GoogleMeet string
	MsTeams    string
	Custom     string
}{
	Zoom:       "ZOOM",
	GoogleMeet: "GOOGLE_MEET",
	MsTeams:    "MS_TEAMS",
	Custom:     "CUSTOM",
}

// IncidentConferenceBridgeTypeAll defines incident conference bridge type list
var IncidentConferenceBridgeTypeAll = []string{
	IncidentConferenceBridgeType.Zoom,
	IncidentConferenceBridgeType.GoogleMeet,
	IncidentConferenceBridgeType.MsTeams,
	IncidentConferenceBridgeType.Custom,
}

// CreateIncidentInput represents the input of a CreateIncident operation.
type CreateIncidentInput struct {
	_        struct{}
	Incident *Incident
}

// CreateIncidentOutput represents the output of a CreateIncident operation.
type CreateIncidentOutput struct {
	_        struct{}
	Incident *Incident
}

// CreateIncident creates a new incident, a title and a severity are required. Depending on the
// affected services this publishes notifications to subscribers, GetIncidentAffected forecasts
// them. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) CreateIncident(input *CreateIncidentInput) (*CreateIncidentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Incident == nil {
		return nil, errors.New("incident input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Incident).Post(apiRoutes.incidents)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200, 201); apiErr != nil {
		return nil, apiErr
	}

	incident := &Incident{}
	err = json.Unmarshal(resp.Body(), incident)
	if err != nil {
		return nil, err
	}

	return &CreateIncidentOutput{Incident: incident}, nil
}

// GetIncidentsInput represents the input of a GetIncidents operation.
type GetIncidentsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 10, Maximum: 100, or 25 with include and 11 when including alertCount
	MaxResults *int

	// describes optional properties that should be included in the response
	// possible values: "affectedServices", "responders", "alertCount", "responderCountJoined", "responderCountTotal"
	Include []*string

	// status of the incident
	// possible values: "DECLARED", "INVESTIGATING", "IDENTIFIED", "MONITORING", "RESOLVED"
	States []*string

	// service IDs of the incident's affected services
	Services []*int64

	// severity of the incident, in range 1..5
	Severity *int

	// Date time string in ISO format, based on the declaration of the incident
	From *string

	// Date time string in ISO format, based on the declaration of the incident
	Until *string
}

// GetIncidentsOutput represents the output of a GetIncidents operation.
type GetIncidentsOutput struct {
	_         struct{}
	Incidents []*Incident
}

// GetIncidents lists existing incidents. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) GetIncidents(input *GetIncidentsInput) (*GetIncidentsOutput, error) {
	if input == nil {
		input = &GetIncidentsInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}
	for _, include := range input.Include {
		q.Add("include", *include)
	}
	for _, state := range input.States {
		q.Add("states", *state)
	}
	for _, service := range input.Services {
		q.Add("services", strconv.FormatInt(*service, 10))
	}
	if input.Severity != nil {
		q.Add("severity", strconv.Itoa(*input.Severity))
	}
	if input.From != nil {
		q.Add("from", *input.From)
	}
	if input.Until != nil {
		q.Add("until", *input.Until)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.incidents, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incidents := make([]*Incident, 0)
	err = json.Unmarshal(resp.Body(), &incidents)
	if err != nil {
		return nil, err
	}

	return &GetIncidentsOutput{Incidents: incidents}, nil
}

// GetIncidentInput represents the input of a GetIncident operation.
type GetIncidentInput struct {
	_          struct{}
	IncidentID *int64

	// describes optional properties that should be included in the response, the affected services
	// and the responders are always included
	// possible values: "alertCount", "responderCountJoined", "responderCountTotal"
	Include []*string
}

// GetIncidentOutput represents the output of a GetIncident operation.
type GetIncidentOutput struct {
	_        struct{}
	Incident *Incident

	// send it back as UpdateIncidentInput.ETag to have the update rejected with a 412 when the
	// incident or its affected services changed in the meantime
	ETag *string
}

// GetIncident gets an incident by id. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) GetIncident(input *GetIncidentInput) (*GetIncidentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.incidents, *input.IncidentID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incident := &Incident{}
	err = json.Unmarshal(resp.Body(), incident)
	if err != nil {
		return nil, err
	}

	output := &GetIncidentOutput{Incident: incident}
	etag := resp.Header().Get("ETag")
	if etag != "" {
		output.ETag = String(etag)
	}

	return output, nil
}

// UpdateIncidentInput represents the input of a UpdateIncident operation.
type UpdateIncidentInput struct {
	_          struct{}
	IncidentID *int64
	Incident   *Incident

	// the ETag returned by GetIncident, sent as If-Match
	ETag *string
}

// UpdateIncidentOutput represents the output of a UpdateIncident operation.
type UpdateIncidentOutput struct {
	_        struct{}
	Incident *Incident
}

// UpdateIncident updates the specific incident. Only the title, summary, status, severity and
// affected services are written. The SDK leaves out an empty title, summary, status or severity,
// which the API keeps as they are: an update does not have to carry the whole incident, but it
// cannot clear the summary either. Affected services left nil are kept too, see
// Incident.AffectedServices. Status updates are published separately through
// CreateIncidentStatusUpdate. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) UpdateIncident(input *UpdateIncidentInput) (*UpdateIncidentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}
	if input.Incident == nil {
		return nil, errors.New("incident input is required")
	}

	req := c.httpClient.R()
	if input.ETag != nil && *input.ETag != "" {
		req.SetHeader("If-Match", *input.ETag)
	}
	resp, err := req.SetBody(input.Incident).Put(fmt.Sprintf("%s/%d", apiRoutes.incidents, *input.IncidentID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incident := &Incident{}
	err = json.Unmarshal(resp.Body(), incident)
	if err != nil {
		return nil, err
	}

	return &UpdateIncidentOutput{Incident: incident}, nil
}

// GetIncidentAffectedInput represents the input of a GetIncidentAffected operation.
type GetIncidentAffectedInput struct {
	_        struct{}
	Incident *Incident
}

// GetIncidentAffectedOutput represents the output of a GetIncidentAffected operation.
type GetIncidentAffectedOutput struct {
	_        struct{}
	Affected *Affected
}

// GetIncidentAffected forecasts the subscribers and status pages that creating or updating the
// incident would notify. Only the affected services of the incident are taken into account. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) GetIncidentAffected(input *GetIncidentAffectedInput) (*GetIncidentAffectedOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Incident == nil {
		return nil, errors.New("incident is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Incident).Post(fmt.Sprintf("%s/publish-info", apiRoutes.incidents))
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

	return &GetIncidentAffectedOutput{Affected: affected}, nil
}

// GetIncidentSubscribersInput represents the input of a GetIncidentSubscribers operation.
type GetIncidentSubscribersInput struct {
	_          struct{}
	IncidentID *int64
}

// GetIncidentSubscribersOutput represents the output of a GetIncidentSubscribers operation.
type GetIncidentSubscribersOutput struct {
	_           struct{}
	Subscribers []*Subscriber
}

// GetIncidentSubscribers gets subscribers of an incident by id. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) GetIncidentSubscribers(input *GetIncidentSubscribersInput) (*GetIncidentSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.incidents, *input.IncidentID))
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

	return &GetIncidentSubscribersOutput{Subscribers: subscribers}, nil
}

// AddIncidentSubscribersInput represents the input of a AddIncidentSubscribers operation.
type AddIncidentSubscribersInput struct {
	_           struct{}
	IncidentID  *int64
	Subscribers *[]Subscriber
}

// AddIncidentSubscribersOutput represents the output of a AddIncidentSubscribers operation.
type AddIncidentSubscribersOutput struct {
	_ struct{}
}

// AddIncidentSubscribers adds subscribers to an incident. Adding one that is already subscribed is
// rejected with a 400. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) AddIncidentSubscribers(input *AddIncidentSubscribersInput) (*AddIncidentSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}
	if input.Subscribers == nil {
		return nil, errors.New("subscriber input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Subscribers).Post(fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.incidents, *input.IncidentID))
	if err != nil {
		return nil, err
	}
	// the API answers 202 without a body, the spec documents 204
	if apiErr := getGenericAPIError(resp, 202, 204); apiErr != nil {
		return nil, apiErr
	}

	return &AddIncidentSubscribersOutput{}, nil
}

// CreateIncidentStatusUpdateInput represents the input of a CreateIncidentStatusUpdate operation.
type CreateIncidentStatusUpdateInput struct {
	_            struct{}
	IncidentID   *int64
	StatusUpdate *StatusUpdate
}

// CreateIncidentStatusUpdateOutput represents the output of a CreateIncidentStatusUpdate operation.
type CreateIncidentStatusUpdateOutput struct {
	_            struct{}
	StatusUpdate *StatusUpdate
}

// CreateIncidentStatusUpdate publishes a status update for an incident. It is the same entity
// CreateStatusUpdate creates, linked to the incident, and like it has to affect at least one service. https://docs.ilert.com/developer-docs/rest-api/api-reference/incidents
func (c *Client) CreateIncidentStatusUpdate(input *CreateIncidentStatusUpdateInput) (*CreateIncidentStatusUpdateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}
	if input.StatusUpdate == nil {
		return nil, errors.New("status update input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.StatusUpdate).Post(fmt.Sprintf("%s/%d/status-updates", apiRoutes.incidents, *input.IncidentID))
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

	return &CreateIncidentStatusUpdateOutput{StatusUpdate: statusUpdate}, nil
}

// GetIncidentLogEntriesInput represents the input of a GetIncidentLogEntries operation.
type GetIncidentLogEntriesInput struct {
	_          struct{}
	IncidentID *int64

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

// GetIncidentLogEntriesOutput represents the output of a GetIncidentLogEntries operation.
type GetIncidentLogEntriesOutput struct {
	_          struct{}
	LogEntries []*IncidentLogEntry
}

// GetIncidentLogEntries lists the timeline of an incident, newest first. Every call of an API token counts
// towards its limit of 120 calls per minute. Once the token exceeds it, the timeline endpoints
// answer with a 429 and then with a 403 "authentication failed" for about a minute, even after the
// rest of the API accepts the token again, so a 403 here does not necessarily mean the token is
// wrong. https://api.ilert.com/api/beta/openapi.json
func (c *Client) GetIncidentLogEntries(input *GetIncidentLogEntriesInput) (*GetIncidentLogEntriesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}
	q, err := logEntriesQuery(input.StartIndex, input.MaxResults, input.From, input.Until, input.SubTypes)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/log-entries?%s", apiRoutes.incidents, *input.IncidentID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	logEntries := make([]*IncidentLogEntry, 0)
	err = json.Unmarshal(resp.Body(), &logEntries)
	if err != nil {
		return nil, err
	}

	return &GetIncidentLogEntriesOutput{LogEntries: logEntries}, nil
}
