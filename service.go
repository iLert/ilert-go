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
	Alias               string         `json:"alias,omitempty"`
	Status              string         `json:"status"`
	Description         string         `json:"description"`
	OneOpenIncidentOnly bool           `json:"oneOpenIncidentOnly"`
	ShowUptimeHistory   bool           `json:"showUptimeHistory"`
	Teams               []TeamShort    `json:"teams"`
	Subscribed          bool           `json:"subscribed,omitempty"`
	Uptime              *ServiceUptime `json:"uptime,omitempty"`
	Incidents           []Incident     `json:"incidents,omitempty"`

	// free-form key-value labels assigned to this service. A nil map disappears from the
	// payload and leaves the labels untouched, a non-nil empty map marshals to
	// "labels":{} and clears them.
	Labels *map[string]string `json:"labels,omitempty"`

	IconUrl string `json:"iconUrl,omitempty"`

	// only returned when "links" is requested through Include. A nil slice disappears from
	// the payload and leaves the links untouched, a non-nil empty slice marshals to
	// "links":[] and clears them.
	Links *[]ServiceLink `json:"links,omitempty"`

	// the status shown on status pages, only returned when "publicStatus" is requested
	// through Include. Read-only, it is derived by the API from the service status.
	PublicStatus string `json:"publicStatus,omitempty"`

	// only returned when "dependencies" is requested through Include. Read-only on this
	// endpoint, the service create and update payloads do not accept it: dependencies are
	// written through the /services/{id}/dependencies endpoints.
	Dependencies []ServiceDependency `json:"dependencies,omitempty"`

	ManagedBy *ManagedBy `json:"managedBy,omitempty"`
	CreatedBy *ManagedBy `json:"createdBy,omitempty"`
}

// ServiceDependencyType defines the strength of a service dependency
var ServiceDependencyType = struct {
	Hard string
	Soft string
}{
	Hard: "HARD",
	Soft: "SOFT",
}

// ServiceDependencyTypeAll defines the service dependency type list
var ServiceDependencyTypeAll = []string{
	ServiceDependencyType.Hard,
	ServiceDependencyType.Soft,
}

// ServiceLink defines a link displayed on a service
type ServiceLink struct {
	Href string `json:"href"`
	Text string `json:"text,omitempty"`
}

// ServiceDependency defines a directed edge of the service dependency graph
type ServiceDependency struct {
	ID int64 `json:"id,omitempty"`

	// the service that depends on the target service
	SourceServiceID int64 `json:"sourceServiceId,omitempty"`

	// the service that is being depended upon
	TargetServiceID int64 `json:"targetServiceId"`

	// the strength of the dependency, "HARD" when the API is left to default it.
	// possible values: "HARD", "SOFT"
	Type string `json:"type,omitempty"`

	Notes string `json:"notes,omitempty"`

	// point in time after which this dependency is no longer considered valid, as a date
	// time string in ISO format. Read-only in practice: the API documents it as writable
	// but ignores it on both POST and PUT, in every date format, and always returns null.
	// Verified against the API on 09.09.2026.
	InvalidAfter string `json:"invalidAfter,omitempty"`

	CreatedAt string     `json:"createdAt,omitempty"`
	UpdatedAt string     `json:"updatedAt,omitempty"`
	CreatedBy *ManagedBy `json:"createdBy,omitempty"`
	ManagedBy *ManagedBy `json:"managedBy,omitempty"`
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
	Subscribed   string
	Uptime       string
	Incidents    string
	Links        string
	Dependencies string
	PublicStatus string
}{
	Subscribed:   "subscribed",
	Uptime:       "uptime",
	Incidents:    "incidents",
	Links:        "links",
	Dependencies: "dependencies",
	PublicStatus: "publicStatus",
}

// ServiceIncludeAll defines included services list
var ServiceIncludeAll = []string{
	ServiceInclude.Subscribed,
	ServiceInclude.Uptime,
	ServiceInclude.Incidents,
	ServiceInclude.Links,
	ServiceInclude.Dependencies,
	ServiceInclude.PublicStatus,
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
		return nil, errors.New("service input is required")
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

// GetServices lists existing services. https://api.ilert.com/api-docs/#tag/Services/paths/~1services/get
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

// GetService gets a service by id. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}/get
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

// GetServiceSubscribers gets subscribers of a service by id. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}~1private-subscribers/get
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

	// describes optional properties that should be included in the response
	// possible values: "subscribed", "uptime", "incidents", "links", "dependencies", "publicStatus"
	Include []*string
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

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s?%s", apiRoutes.services, url.PathEscape(*input.ServiceName), q.Encode()))
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

// AddServiceSubscribers adds a new subscriber to a service. https://api.ilert.com/api-docs/#tag/Services/paths/~1services~1{id}~1private-subscribers/post
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

// GetServiceDependenciesInput represents the input of a GetServiceDependencies operation.
type GetServiceDependenciesInput struct {
	_         struct{}
	ServiceID *int64

	// whether to also return the dependencies of the dependencies
	Transitive *bool

	// how many levels to follow when Transitive is set
	Depth *int
}

// GetServiceDependenciesOutput represents the output of a GetServiceDependencies operation.
type GetServiceDependenciesOutput struct {
	_            struct{}
	Dependencies []*ServiceDependency
}

// GetServiceDependencies lists the dependencies of a service. https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) GetServiceDependencies(input *GetServiceDependenciesInput) (*GetServiceDependenciesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}

	q := url.Values{}
	if input.Transitive != nil {
		q.Add("transitive", strconv.FormatBool(*input.Transitive))
	}
	if input.Depth != nil {
		q.Add("depth", strconv.Itoa(*input.Depth))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/dependencies?%s", apiRoutes.services, *input.ServiceID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	dependencies := make([]*ServiceDependency, 0)
	err = json.Unmarshal(resp.Body(), &dependencies)
	if err != nil {
		return nil, err
	}

	return &GetServiceDependenciesOutput{Dependencies: dependencies}, nil
}

// GetServiceDependencyInput represents the input of a GetServiceDependency operation.
type GetServiceDependencyInput struct {
	_         struct{}
	ServiceID *int64

	// the id of the dependency edge, not the id of either service
	EdgeID *int64
}

// GetServiceDependencyOutput represents the output of a GetServiceDependency operation.
type GetServiceDependencyOutput struct {
	_          struct{}
	Dependency *ServiceDependency
}

// GetServiceDependency gets a single dependency of a service. https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) GetServiceDependency(input *GetServiceDependencyInput) (*GetServiceDependencyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}
	if input.EdgeID == nil {
		return nil, errors.New("service dependency edge id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/dependencies/%d", apiRoutes.services, *input.ServiceID, *input.EdgeID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	dependency := &ServiceDependency{}
	err = json.Unmarshal(resp.Body(), dependency)
	if err != nil {
		return nil, err
	}

	return &GetServiceDependencyOutput{Dependency: dependency}, nil
}

// CreateServiceDependencyInput represents the input of a CreateServiceDependency operation.
type CreateServiceDependencyInput struct {
	_          struct{}
	ServiceID  *int64
	Dependency *ServiceDependency
}

// CreateServiceDependencyOutput represents the output of a CreateServiceDependency operation.
type CreateServiceDependencyOutput struct {
	_          struct{}
	Dependency *ServiceDependency
}

// CreateServiceDependency creates a dependency on a service. https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) CreateServiceDependency(input *CreateServiceDependencyInput) (*CreateServiceDependencyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}
	if input.Dependency == nil {
		return nil, errors.New("service dependency input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Dependency).Post(fmt.Sprintf("%s/%d/dependencies", apiRoutes.services, *input.ServiceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	dependency := &ServiceDependency{}
	err = json.Unmarshal(resp.Body(), dependency)
	if err != nil {
		return nil, err
	}

	return &CreateServiceDependencyOutput{Dependency: dependency}, nil
}

// SetServiceDependenciesInput represents the input of a SetServiceDependencies operation.
type SetServiceDependenciesInput struct {
	_         struct{}
	ServiceID *int64

	// replaces every dependency of the service, an empty slice removes them all
	Dependencies []ServiceDependency
}

// SetServiceDependenciesOutput represents the output of a SetServiceDependencies operation.
type SetServiceDependenciesOutput struct {
	_            struct{}
	Dependencies []*ServiceDependency
}

// SetServiceDependencies replaces all dependencies of a service. https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) SetServiceDependencies(input *SetServiceDependenciesInput) (*SetServiceDependenciesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}
	if input.Dependencies == nil {
		input.Dependencies = []ServiceDependency{}
	}

	resp, err := c.httpClient.R().SetBody(input.Dependencies).Put(fmt.Sprintf("%s/%d/dependencies", apiRoutes.services, *input.ServiceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	dependencies := make([]*ServiceDependency, 0)
	err = json.Unmarshal(resp.Body(), &dependencies)
	if err != nil {
		return nil, err
	}

	return &SetServiceDependenciesOutput{Dependencies: dependencies}, nil
}

// DeleteServiceDependencyInput represents the input of a DeleteServiceDependency operation.
type DeleteServiceDependencyInput struct {
	_         struct{}
	ServiceID *int64

	// the id of the dependency edge, not the id of either service
	EdgeID *int64
}

// DeleteServiceDependencyOutput represents the output of a DeleteServiceDependency operation.
type DeleteServiceDependencyOutput struct {
	_ struct{}
}

// DeleteServiceDependency deletes a dependency of a service. https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) DeleteServiceDependency(input *DeleteServiceDependencyInput) (*DeleteServiceDependencyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ServiceID == nil {
		return nil, errors.New("service id is required")
	}
	if input.EdgeID == nil {
		return nil, errors.New("service dependency edge id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d/dependencies/%d", apiRoutes.services, *input.ServiceID, *input.EdgeID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteServiceDependencyOutput{}, nil
}
