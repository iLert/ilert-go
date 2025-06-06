package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// HeartbeatMonitor definition https://api.ilert.com/api-docs/#tag/heartbeat-monitors
type HeartbeatMonitor struct {
	ID             int64        `json:"id,omitempty"`
	Name           string       `json:"name"`
	State          string       `json:"state,omitempty"`
	IntervalSec    int64        `json:"intervalSec"`
	AlertSummary   string       `json:"alertSummary,omitempty"`
	CreatedAt      string       `json:"createdAt,omitempty"`
	UpdatedAt      string       `json:"updatedAt,omitempty"`
	AlertSource    *AlertSource `json:"alertSource,omitempty"`
	Teams          []TeamShort  `json:"teams,omitempty"`
	IntegrationKey string       `json:"integrationKey,omitempty"`
	IntegrationUrl string       `json:"integrationUrl,omitempty"`
}

// CreateHeartbeatMonitorInput represents the input of a CreateHeartbeatMonitor operation.
type CreateHeartbeatMonitorInput struct {
	_                struct{}
	HeartbeatMonitor *HeartbeatMonitor

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// CreateHeartbeatMonitorOutput represents the output of a CreateHeartbeatMonitor operation.
type CreateHeartbeatMonitorOutput struct {
	_                struct{}
	HeartbeatMonitor *HeartbeatMonitor
}

// CreateHeartbeatMonitor creates a new heartbeat monitor resource. https://api.ilert.com/api-docs/#tag/heartbeat-monitors/post/heartbeat-monitors
func (c *Client) CreateHeartbeatMonitor(input *CreateHeartbeatMonitorInput) (*CreateHeartbeatMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.HeartbeatMonitor == nil {
		return nil, errors.New("heartbeat monitor input is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.HeartbeatMonitor).Post(fmt.Sprintf("%s?%s", apiRoutes.heartbeatMonitors, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	HeartbeatMonitor := &HeartbeatMonitor{}
	err = json.Unmarshal(resp.Body(), HeartbeatMonitor)
	if err != nil {
		return nil, err
	}

	return &CreateHeartbeatMonitorOutput{HeartbeatMonitor: HeartbeatMonitor}, nil
}

// GetHeartbeatMonitorInput represents the input of a GetHeartbeatMonitor operation.
type GetHeartbeatMonitorInput struct {
	_                  struct{}
	HeartbeatMonitorID *int64

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// GetHeartbeatMonitorOutput represents the output of a GetHeartbeatMonitor operation.
type GetHeartbeatMonitorOutput struct {
	_                struct{}
	HeartbeatMonitor *HeartbeatMonitor
}

// GetHeartbeatMonitor gets the heartbeat monitors resource with specified id. https://api.ilert.com/api-docs/#tag/heartbeat-monitors/get/heartbeat-monitors/{id}
func (c *Client) GetHeartbeatMonitor(input *GetHeartbeatMonitorInput) (*GetHeartbeatMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.HeartbeatMonitorID == nil {
		return nil, errors.New("heartbeat monitor id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.heartbeatMonitors, *input.HeartbeatMonitorID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	HeartbeatMonitor := &HeartbeatMonitor{}
	err = json.Unmarshal(resp.Body(), HeartbeatMonitor)
	if err != nil {
		return nil, err
	}

	return &GetHeartbeatMonitorOutput{HeartbeatMonitor: HeartbeatMonitor}, nil
}

// GetHeartbeatMonitorsInput represents the input of a GetHeartbeatMonitors operation.
type GetHeartbeatMonitorsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetHeartbeatMonitorsOutput represents the output of a GetHeartbeatMonitors operation.
type GetHeartbeatMonitorsOutput struct {
	_                 struct{}
	HeartbeatMonitors []*HeartbeatMonitor
}

// GetHeartbeatMonitors lists existing heartbeat monitor resources. https://api.ilert.com/api-docs/#tag/heartbeat-monitors/get/heartbeat-monitors
func (c *Client) GetHeartbeatMonitors(input *GetHeartbeatMonitorsInput) (*GetHeartbeatMonitorsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.heartbeatMonitors, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	HeartbeatMonitors := make([]*HeartbeatMonitor, 0)
	err = json.Unmarshal(resp.Body(), &HeartbeatMonitors)
	if err != nil {
		return nil, err
	}

	return &GetHeartbeatMonitorsOutput{HeartbeatMonitors: HeartbeatMonitors}, nil
}

// SearchHeartbeatMonitorInput represents the input of a SearchHeartbeatMonitor operation.
type SearchHeartbeatMonitorInput struct {
	_                    struct{}
	HeartbeatMonitorName *string
}

// SearchHeartbeatMonitorOutput represents the output of a SearchHeartbeatMonitor operation.
type SearchHeartbeatMonitorOutput struct {
	_                struct{}
	HeartbeatMonitor *HeartbeatMonitor
}

// SearchHeartbeatMonitor gets the heartbeat monitor resource with specified name.
func (c *Client) SearchHeartbeatMonitor(input *SearchHeartbeatMonitorInput) (*SearchHeartbeatMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.HeartbeatMonitorName == nil {
		return nil, errors.New("heartbeat monitor name is required")
	}

	q := url.Values{}
	q.Add("include", "integrationUrl")

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s?%s", apiRoutes.heartbeatMonitors, *input.HeartbeatMonitorName, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	HeartbeatMonitor := &HeartbeatMonitor{}
	err = json.Unmarshal(resp.Body(), HeartbeatMonitor)
	if err != nil {
		return nil, err
	}

	return &SearchHeartbeatMonitorOutput{HeartbeatMonitor: HeartbeatMonitor}, nil
}

// UpdateHeartbeatMonitorInput represents the input of a UpdateHeartbeatMonitor operation.
type UpdateHeartbeatMonitorInput struct {
	_                  struct{}
	HeartbeatMonitorID *int64
	HeartbeatMonitor   *HeartbeatMonitor

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// UpdateHeartbeatMonitorOutput represents the output of a UpdateHeartbeatMonitor operation.
type UpdateHeartbeatMonitorOutput struct {
	_                struct{}
	HeartbeatMonitor *HeartbeatMonitor
}

// UpdateHeartbeatMonitor updates an existing heartbeat monitor resource. https://api.ilert.com/api-docs/#tag/heartbeat-monitors/put/heartbeat-monitors/{id}
func (c *Client) UpdateHeartbeatMonitor(input *UpdateHeartbeatMonitorInput) (*UpdateHeartbeatMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.HeartbeatMonitor == nil {
		return nil, errors.New("heartbeat monitor input is required")
	}
	if input.HeartbeatMonitorID == nil {
		return nil, errors.New("heartbeat monitor id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.HeartbeatMonitor).Put(fmt.Sprintf("%s/%d?%s", apiRoutes.heartbeatMonitors, *input.HeartbeatMonitorID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	HeartbeatMonitor := &HeartbeatMonitor{}
	err = json.Unmarshal(resp.Body(), HeartbeatMonitor)
	if err != nil {
		return nil, err
	}

	return &UpdateHeartbeatMonitorOutput{HeartbeatMonitor: HeartbeatMonitor}, nil
}

// DeleteHeartbeatMonitorInput represents the input of a DeleteHeartbeatMonitor operation.
type DeleteHeartbeatMonitorInput struct {
	_                  struct{}
	HeartbeatMonitorID *int64
}

// DeleteHeartbeatMonitorOutput represents the output of a DeleteHeartbeatMonitor operation.
type DeleteHeartbeatMonitorOutput struct {
	_ struct{}
}

// DeleteHeartbeatMonitor deletes the specified heartbeat monitor resource. https://api.ilert.com/api-docs/#tag/heartbeat-monitors/delete/heartbeat-monitors/{id}
func (c *Client) DeleteHeartbeatMonitor(input *DeleteHeartbeatMonitorInput) (*DeleteHeartbeatMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.HeartbeatMonitorID == nil {
		return nil, errors.New("heartbeat monitor id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.heartbeatMonitors, *input.HeartbeatMonitorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteHeartbeatMonitorOutput{}, nil
}
