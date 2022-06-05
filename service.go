package ilert

import (
	"encoding/json"
	"errors"
)

type Service struct {
	ID                  int64         `json:"id"`
	Name                string        `json:"name"`
	Status              string        `json:"status"`
	Description         string        `json:"description"`
	OneOpenIncidentOnly bool          `json:"oneOpenIncidentOnly"`
	ShowUptimeHistory   bool          `json:"showUptimeHistory"`
	Teams               []TeamShort   `json:"teams"`
	Subscribed          bool          `json:"subscribed,omitempty"`
	Uptime              ServiceUptime `json:"uptime,omitempty"`
	Incidents           Incident      `json:"incidents,omitempty"`
}

//ServiceUptime defines services uptime
type ServiceUptime struct {
	UptimePercentage ServiceUptimePercentage `json:"uptimePercentage"`
}

//ServiceUptimePercentage defines service uptime percentage
type ServiceUptimePercentage struct {
	RangeStart string        `json:"rangeStart"` // Date time string in ISO format
	RangeEnd   string        `json:"rangeEnd"`   // Date time string in ISO format
	Outages    ServiceOutage `json:"outages"`
}

type ServiceOutage struct {
	Status string `json:"status"`
	From   string `json:"from"`  // Date time string in ISO format
	Until  string `json:"until"` // Date time string in ISO format
}

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

//ServiceUptimeOnly defines services
type ServiceUptimeOnly struct {
	Id                  int64         `json:"id"`
	Name                string        `json:"name"`
	Status              string        `json:"status"`
	Description         string        `json:"description"`
	OneOpenIncidentOnly bool          `json:"oneOpenIncidentOnly"`
	ShowUptimeHistory   bool          `json:"showUptimeHistory"`
	Teams               []TeamShort   `json:"teams"`
	Uptime              ServiceUptime `json:"uptime"`
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
