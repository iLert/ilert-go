package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Incident struct {
	ID               int64               `json:"id"`
	Summary          string              `json:"summary"`
	Status           string              `json:"status"`
	Message          string              `json:"message"`
	SendNotification bool                `json:"sendNotification"`
	CreatedAt        string              `json:"createdAt"` // Date time string in ISO format
	UpdatedAt        string              `json:"updatedAt"` // Date time string in ISO format
	AffectedServices *[]AffectedServices `json:"affectedServices"`
	ResolvedOn       string              `json:"resolvedOn,omitempty"` // Date time string in ISO format
	Subscribed       bool                `json:"subscribed,omitempty"`
	AffectedTeams    *[]TeamShort        `json:"affectedTeams,omitempty"`
}

type AffectedServices struct {
	Impact  string  `json:"impact"`
	Service Service `json:"service"`
}

type Subscriber struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UIMenuItem struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

type Affected struct {
	StatusPagesInfo    UIMenuItem `json:"statusPagesInfo"`
	PrivateStatusPages int64      `json:"privateStatusPages"`
	PublicStatusPages  int64      `json:"publicStatusPages"`
	PrivateSubscribers int64      `json:"privateSubscribers"`
	PublicSubscribers  int64      `json:"publicSubscribers"`
}

var IncidentInclude = struct {
	Subscribed    string
	AffectedTeams string
	History       string
}{
	Subscribed:    "subscribed",
	AffectedTeams: "affectedTeams",
	History:       "history",
}

var IncidentType = struct {
	User string
	Team string
}{
	User: "USER",
	Team: "TEAM",
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

// CreateIncident creates a new incident. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents/post
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
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
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
	// Default: 50, Maximum: 100
	MaxResults *int

	// describes optional properties that should be included in the response
	Include []*string

	// state of the incident
	States []*string

	// service IDs of the incident's affected services
	Services []*int64

	// Date time string in ISO format
	From *string

	// Date time string in ISO format
	Until *string
}

// GetIncidentsOutput represents the output of a GetIncidents operation.
type GetIncidentsOutput struct {
	_         struct{}
	Incidents []*Incident
}

// GetIncidents lists incident sources. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents/get
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
		q.Add("state", *state)
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

	// describes optional properties that should be included in the response
	Include []*string
}

// GetIncidentOutput represents the output of a GetIncident operation.
type GetIncidentOutput struct {
	_        struct{}
	Incident *Incident
}

// GetIncident gets an incident by ID. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents~1{id}/get
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

	var url = fmt.Sprintf("%s/%d?%s", apiRoutes.incidents, *input.IncidentID, q.Encode())

	resp, err := c.httpClient.R().Get(url)
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

	return &GetIncidentOutput{Incident: incident}, nil
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

// GetIncidentSubscribers gets subscribers of an incident by ID. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents~1{id}~1private-subscribers/get
func (c *Client) GetIncidentSubscribers(input *GetIncidentSubscribersInput) (*GetIncidentSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("incident id is required")
	}

	var url = fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.incidents, *input.IncidentID)

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

	return &GetIncidentSubscribersOutput{Subscribers: subscribers}, nil
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

// GetIncidentAffected forecasts the affected subscribers and statuspages. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents~1publish-info/post
func (c *Client) GetIncidentAffected(input *GetIncidentAffectedInput) (*GetIncidentAffectedOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Incident == nil {
		return nil, errors.New("incident is required")
	}

	url := fmt.Sprintf("%s/publish-info", apiRoutes.incidents)

	resp, err := c.httpClient.R().SetBody(input.Incident).Post(url)
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

// AddIncidentSubscribers adds a new subscriber to an incident. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents~1{id}~1private-subscribers/post
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

	url := fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.incidents, *input.IncidentID)

	resp, err := c.httpClient.R().SetBody(input.Subscribers).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return nil, apiErr
	}

	subscribers := make([]*Subscriber, 0)
	err = json.Unmarshal(resp.Body(), &subscribers)
	if err != nil {
		return nil, err
	}

	return &AddIncidentSubscribersOutput{}, nil
}

// UpdateIncidentInput represents the input of a UpdateIncident operation.
type UpdateIncidentInput struct {
	_          struct{}
	IncidentID *int64
	Incident   *Incident
}

// UpdateIncidentOutput represents the output of a UpdateIncident operation.
type UpdateIncidentOutput struct {
	_        struct{}
	Incident *Incident
}

// UpdateIncident updates the specific incident. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents~1{id}/put
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

	url := fmt.Sprintf("%s/%d", apiRoutes.incidents, *input.IncidentID)

	resp, err := c.httpClient.R().SetBody(input.Incident).Put(url)
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
