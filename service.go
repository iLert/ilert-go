package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Service definition https://api.ilert.com/api-docs/#tag/Services
type Service struct {
	ID                  int64          `json:"id"`
	Name                string         `json:"name"`
	Status              string         `json:"status"`
	Description         string         `json:"description"`
	OneOpenIncidentOnly bool           `json:"oneOpenIncidentOnly"`
	ShowUptimeHistory   bool           `json:"showUptimeHistory"`
	Teams               []TeamShort    `json:"teams"`
	Subscribed          bool           `json:"subscribed,omitempty"`
	Uptime              *ServiceUptime `json:"uptime,omitempty"`
	Incidents           []Incident     `json:"incidents,omitempty"`
}

// ServiceUptime defines services uptime
type ServiceUptime struct {
	RangeStart       string                   `json:"rangeStart"` // Date time string in ISO format
	RangeEnd         string                   `json:"rangeEnd"`   // Date time string in ISO format
	Outages          []ServiceOutage          `json:"outages"`
	UptimePercentage *ServiceUptimePercentage `json:"uptimePercentage"`
}

// ServiceUptimePercentage defines service uptime percentage
type ServiceUptimePercentage struct {
	P90 float64 `json:"p90"`
	P60 float64 `json:"p60"`
	P30 float64 `json:"p30"`
}

// ServiceOutage defines services outage
type ServiceOutage struct {
	Status string `json:"status"`
	From   string `json:"from"`  // Date time string in ISO format
	Until  string `json:"until"` // Date time string in ISO format
}

// ServiceStatus defines services status
var ServiceStatus = struct {
	Operational      string
	UnderMaintenance string
	Degraded         string
	PartialOutage    string
	MajorOutage      string
}{
	Operational:      "OPERATIONAL",
	UnderMaintenance: "UNDER_MAINTENANCE",
	Degraded:         "DEGRADED",
	PartialOutage:    "PARTIAL_OUTAGE",
	MajorOutage:      "MAJOR_OUTAGE",
}

// ServiceStatusAll defines services status list
var ServiceStatusAll = []string{
	ServiceStatus.Operational,
	ServiceStatus.UnderMaintenance,
	ServiceStatus.Degraded,
	ServiceStatus.PartialOutage,
	ServiceStatus.MajorOutage,
}

// ServiceInclude defines included services
var ServiceInclude = struct {
	Subscribed string
	Uptime     string
	Incidents  string
}{
	Subscribed: "subscribed",
	Uptime:     "uptime",
	Incidents:  "incidents",
}

// ServiceIncludeAll defines included services list
var ServiceIncludeAll = []string{
	ServiceInclude.Subscribed,
	ServiceInclude.Uptime,
	ServiceInclude.Incidents,
}

// CreateServiceInput represents the input of a CreateService operation.
type CreateServiceInput struct {
	_       struct{}
	Service *Service
}

// CreateServiceOutput represents the output of a CreateService operation.
type CreateServiceOutput struct {
	_       struct{}
	Service *Service
}

// CreateService creates a new service. https://api.ilert.com/api-docs/#tag/Services/paths/~1services/post
func (c *Client) CreateService(input *CreateServiceInput) (*CreateServiceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Service == nil {
		return nil, errors.New("Service input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Service).Post(apiRoutes.services)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	service := &Service{}
	err = json.Unmarshal(resp.Body(), service)
	if err != nil {
		return nil, err
	}

	return &CreateServiceOutput{Service: service}, nil
}

// GetServicesInput represents the input of a GetServices operation.
type GetServicesInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	// Default: 0
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 10, Maximum: 25 or 100 without include
	MaxResults *int

	// describes optional properties that should be included in the response
	Include []*string
}

// GetServicesOutput represents the output of a GetServices operation.
type GetServicesOutput struct {
	_        struct{}
	Services []*Service
}

// GetServices lists service sources. https://api.ilert.com/api-docs/#tag/Services/paths/~1services/get
func (c *Client) GetServices(input *GetServicesInput) (*GetServicesOutput, error) {
	if input == nil {
		input = &GetServicesInput{}
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

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.services, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	services := make([]*Service, 0)
	err = json.Unmarshal(resp.Body(), &services)
	if err != nil {
		return nil, err
	}

	return &GetServicesOutput{Services: services}, nil
}

// GetServiceInput represents the input of a GetService operation.
type GetServiceInput struct {
	_         struct{}
	ServiceID *int64

	// describes optional properties that should be included in the response
	Include []*string
}

// GetServiceOutput represents the output of a GetService operation.
type GetServiceOutput struct {
	_       struct{}
	Service *Service
}

// GetService gets a service by ID. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}/get
func (c *Client) GetService(input *GetServiceInput) (*GetServiceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	var url = fmt.Sprintf("%s/%d?%s", apiRoutes.services, *input.ServiceID, q.Encode())

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	service := &Service{}
	err = json.Unmarshal(resp.Body(), service)
	if err != nil {
		return nil, err
	}

	return &GetServiceOutput{Service: service}, nil
}

// GetServiceSubscribersInput represents the input of a GetServiceSubscribers operation.
type GetServiceSubscribersInput struct {
	_         struct{}
	ServiceID *int64
}

// GetServiceSubscribersOutput represents the output of a GetServiceSubscribers operation.
type GetServiceSubscribersOutput struct {
	_           struct{}
	Subscribers []*Subscriber
}

// GetServiceSubscribers gets subscribers of a service by ID. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}~1private-subscribers/get
func (c *Client) GetServiceSubscribers(input *GetServiceSubscribersInput) (*GetServiceSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}

	var url = fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.services, *input.ServiceID)

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

	return &GetServiceSubscribersOutput{Subscribers: subscribers}, nil
}

// SearchServiceInput represents the input of a SearchService operation.
type SearchServiceInput struct {
	_           struct{}
	ServiceName *string
}

// SearchServiceOutput represents the output of a SearchService operation.
type SearchServiceOutput struct {
	_       struct{}
	Service *Service
}

// SearchService gets the service with specified name.
func (c *Client) SearchService(input *SearchServiceInput) (*SearchServiceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceName == nil {
		return nil, errors.New("service name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.services, *input.ServiceName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	service := &Service{}
	err = json.Unmarshal(resp.Body(), service)
	if err != nil {
		return nil, err
	}

	return &SearchServiceOutput{Service: service}, nil
}

// UpdateServiceInput represents the input of a UpdateService operation.
type UpdateServiceInput struct {
	_         struct{}
	ServiceID *int64
	Service   *Service
}

// UpdateServiceOutput represents the output of a UpdateService operation.
type UpdateServiceOutput struct {
	_       struct{}
	Service *Service
}

// UpdateService updates the specific service. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}/put
func (c *Client) UpdateService(input *UpdateServiceInput) (*UpdateServiceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}
	if input.Service == nil {
		return nil, errors.New("service input is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.services, *input.ServiceID)

	resp, err := c.httpClient.R().SetBody(input.Service).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	service := &Service{}
	err = json.Unmarshal(resp.Body(), service)
	if err != nil {
		return nil, err
	}

	return &UpdateServiceOutput{Service: service}, nil
}

// AddServiceSubscribersInput represents the input of a AddServiceSubscribers operation.
type AddServiceSubscribersInput struct {
	_           struct{}
	ServiceID   *int64
	Subscribers *[]Subscriber
}

// AddServiceSubscribersOutput represents the output of a AddServiceSubscribers operation.
type AddServiceSubscribersOutput struct {
	_ struct{}
}

// AddServiceSubscribers adds a new subscriber to an service. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}~1private-subscribers/post
func (c *Client) AddServiceSubscribers(input *AddServiceSubscribersInput) (*AddServiceSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}
	if input.Subscribers == nil {
		return nil, errors.New("subscriber input is required")
	}

	url := fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.services, *input.ServiceID)

	resp, err := c.httpClient.R().SetBody(input.Subscribers).Put(url)
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

	return &AddServiceSubscribersOutput{}, nil
}

// DeleteServiceInput represents the input of a DeleteService operation.
type DeleteServiceInput struct {
	_         struct{}
	ServiceID *int64
}

// DeleteServiceOutput represents the output of a DeleteService operation.
type DeleteServiceOutput struct {
	_ struct{}
}

// DeleteService deletes the specified service. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}/delete
func (c *Client) DeleteService(input *DeleteServiceInput) (*DeleteServiceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.services, *input.ServiceID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteServiceOutput{}, nil
}
