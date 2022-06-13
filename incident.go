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
		return nil, errors.New("Incident input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Incident).Post(apiRoutes.incidents)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
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

var Include = struct {
	Subscribed    string
	AffectedTeams string
	History       string
}{
	Subscribed:    "subscribed",
	AffectedTeams: "affectedTeams",
	History:       "history",
}

// GetIncident gets an incident by ID. https://api.ilert.com/api-docs/#tag/Incidents/paths/~1incidents~1{id}/get
func (c *Client) GetIncident(input *GetIncidentInput) (*GetIncidentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentID == nil {
		return nil, errors.New("Incident id is required")
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
