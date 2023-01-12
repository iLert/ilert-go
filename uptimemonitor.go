package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// UptimeMonitor definition https://api.ilert.com/api-docs/#tag/Uptime-Monitors
type UptimeMonitor struct {
	ID                              int64                    `json:"id,omitempty"`
	Name                            string                   `json:"name"`
	Region                          string                   `json:"region"`
	CheckType                       string                   `json:"checkType"`
	CheckParams                     UptimeMonitorCheckParams `json:"checkParams"`
	IntervalSec                     int                      `json:"intervalSec,omitempty"`                     // default: 300
	TimeoutMs                       int                      `json:"timeoutMs,omitempty"`                       // default: 30000
	CreateIncidentAfterFailedChecks int                      `json:"createIncidentAfterFailedChecks,omitempty"` // @deprecated
	CreateAlertAfterFailedChecks    int                      `json:"createAlertAfterFailedChecks,omitempty"`    // default: 1
	EscalationPolicy                *EscalationPolicy        `json:"escalationPolicy,omitempty"`
	Paused                          bool                     `json:"paused,omitempty"`   // default: false
	EmbedURL                        string                   `json:"embedURL,omitempty"` // read only
	ShareURL                        string                   `json:"shareURL,omitempty"` // read only
	Status                          string                   `json:"status,omitempty"`
	LastStatusChange                string                   `json:"lastStatusChange,omitempty"` // Date time string in ISO format
}

// UptimeMonitorCheckParams definition
type UptimeMonitorCheckParams struct {
	Host                     string   `json:"host,omitempty"`
	Port                     int      `json:"port,omitempty"`
	URL                      string   `json:"url,omitempty"`
	ResponseKeywords         []string `json:"responseKeywords,omitempty"` // only for `http`
	AlertBeforeSec           int      `json:"alertBeforeSec,omitempty"`   // only for `ssl`
	AlertOnFingerprintChange bool     `json:"alertOnFingerprintChange"`   // only for `ssl`
}

// UptimeMonitorStatuses defines uptime monitor statuses
var UptimeMonitorStatuses = struct {
	Up      string
	Down    string
	Warning string
	Paused  string
	Unknown string
}{
	Up:      "up",
	Down:    "down",
	Warning: "warn",
	Paused:  "paused",
	Unknown: "unknown",
}

// UptimeMonitorStatusesAll defines uptime monitor statuses list
var UptimeMonitorStatusesAll = []string{
	UptimeMonitorStatuses.Up,
	UptimeMonitorStatuses.Down,
	UptimeMonitorStatuses.Warning,
	UptimeMonitorStatuses.Paused,
	UptimeMonitorStatuses.Unknown,
}

// UptimeMonitorRegions defines uptime monitor regions
var UptimeMonitorRegions = struct {
	EU string
	US string
}{
	EU: "EU",
	US: "US",
}

// UptimeMonitorRegionsAll defines uptime monitor regions list
var UptimeMonitorRegionsAll = []string{
	UptimeMonitorRegions.EU,
	UptimeMonitorRegions.US,
}

// UptimeMonitorCheckTypes defines uptime monitor check types
var UptimeMonitorCheckTypes = struct {
	HTTP string
	Ping string
	TCP  string
	UDP  string
	SSL  string
}{
	HTTP: "http",
	Ping: "ping",
	TCP:  "tcp",
	UDP:  "udp",
	SSL:  "ssl",
}

// UptimeMonitorCheckTypesAll defines uptime monitor check types list
var UptimeMonitorCheckTypesAll = []string{
	UptimeMonitorCheckTypes.HTTP,
	UptimeMonitorCheckTypes.Ping,
	UptimeMonitorCheckTypes.TCP,
	UptimeMonitorCheckTypes.UDP,
	UptimeMonitorCheckTypes.SSL,
}

// CreateUptimeMonitorInput represents the input of a CreateUptimeMonitor operation.
type CreateUptimeMonitorInput struct {
	_             struct{}
	UptimeMonitor *UptimeMonitor
}

// CreateUptimeMonitorOutput represents the output of a CreateUptimeMonitor operation.
type CreateUptimeMonitorOutput struct {
	_             struct{}
	UptimeMonitor *UptimeMonitor
}

// CreateUptimeMonitor creates a new uptime monitor. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors/post
func (c *Client) CreateUptimeMonitor(input *CreateUptimeMonitorInput) (*CreateUptimeMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UptimeMonitor == nil {
		return nil, errors.New("uptime monitor input is required")
	}

	if input.UptimeMonitor.CreateAlertAfterFailedChecks != input.UptimeMonitor.CreateIncidentAfterFailedChecks {
		input.UptimeMonitor.CreateAlertAfterFailedChecks = input.UptimeMonitor.CreateIncidentAfterFailedChecks
	}

	resp, err := c.httpClient.R().SetBody(input.UptimeMonitor).Post(apiRoutes.uptimeMonitors)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	uptimeMonitor := &UptimeMonitor{}
	err = json.Unmarshal(resp.Body(), uptimeMonitor)
	if err != nil {
		return nil, err
	}

	return &CreateUptimeMonitorOutput{UptimeMonitor: uptimeMonitor}, nil
}

// GetUptimeMonitorInput represents the input of a GetUptimeMonitor operation.
type GetUptimeMonitorInput struct {
	_               struct{}
	UptimeMonitorID *int64
}

// GetUptimeMonitorOutput represents the output of a GetUptimeMonitor operation.
type GetUptimeMonitorOutput struct {
	_             struct{}
	UptimeMonitor *UptimeMonitor
}

// GetUptimeMonitor gets the uptime monitor with specified id. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors~1{id}/get
func (c *Client) GetUptimeMonitor(input *GetUptimeMonitorInput) (*GetUptimeMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UptimeMonitorID == nil {
		return nil, errors.New("uptime monitor id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.uptimeMonitors, *input.UptimeMonitorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	uptimeMonitor := &UptimeMonitor{}
	err = json.Unmarshal(resp.Body(), uptimeMonitor)
	if err != nil {
		return nil, err
	}

	return &GetUptimeMonitorOutput{UptimeMonitor: uptimeMonitor}, nil
}

// GetUptimeMonitorsInput represents the input of a GetUptimeMonitors operation.
type GetUptimeMonitorsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetUptimeMonitorsOutput represents the output of a GetUptimeMonitors operation.
type GetUptimeMonitorsOutput struct {
	_              struct{}
	UptimeMonitors []*UptimeMonitor
}

// GetUptimeMonitors gets list uptime monitors. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors/get
func (c *Client) GetUptimeMonitors(input *GetUptimeMonitorsInput) (*GetUptimeMonitorsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.uptimeMonitors, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	uptimeMonitors := make([]*UptimeMonitor, 0)
	err = json.Unmarshal(resp.Body(), &uptimeMonitors)
	if err != nil {
		return nil, err
	}

	return &GetUptimeMonitorsOutput{UptimeMonitors: uptimeMonitors}, nil
}

// SearchUptimeMonitorInput represents the input of a SearchUptimeMonitor operation.
type SearchUptimeMonitorInput struct {
	_                 struct{}
	UptimeMonitorName *string
}

// SearchUptimeMonitorOutput represents the output of a SearchUptimeMonitor operation.
type SearchUptimeMonitorOutput struct {
	_             struct{}
	UptimeMonitor *UptimeMonitor
}

// SearchUptimeMonitor gets the UptimeMonitor with specified name.
func (c *Client) SearchUptimeMonitor(input *SearchUptimeMonitorInput) (*SearchUptimeMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UptimeMonitorName == nil {
		return nil, errors.New("uptime monitor name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.uptimeMonitors, *input.UptimeMonitorName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	uptimeMonitor := &UptimeMonitor{}
	err = json.Unmarshal(resp.Body(), uptimeMonitor)
	if err != nil {
		return nil, err
	}

	return &SearchUptimeMonitorOutput{UptimeMonitor: uptimeMonitor}, nil
}

// UpdateUptimeMonitorInput represents the input of a UpdateUptimeMonitor operation.
type UpdateUptimeMonitorInput struct {
	_               struct{}
	UptimeMonitorID *int64
	UptimeMonitor   *UptimeMonitor
}

// UpdateUptimeMonitorOutput represents the output of a UpdateUptimeMonitor operation.
type UpdateUptimeMonitorOutput struct {
	_             struct{}
	UptimeMonitor *UptimeMonitor
}

// UpdateUptimeMonitor updates an existing uptime monitor. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors~1{id}/put
func (c *Client) UpdateUptimeMonitor(input *UpdateUptimeMonitorInput) (*UpdateUptimeMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UptimeMonitor == nil {
		return nil, errors.New("uptime monitor input is required")
	}
	if input.UptimeMonitorID == nil {
		return nil, errors.New("uptime monitor id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.UptimeMonitor).Put(fmt.Sprintf("%s/%d", apiRoutes.uptimeMonitors, *input.UptimeMonitorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	uptimeMonitor := &UptimeMonitor{}
	err = json.Unmarshal(resp.Body(), uptimeMonitor)
	if err != nil {
		return nil, err
	}

	return &UpdateUptimeMonitorOutput{UptimeMonitor: uptimeMonitor}, nil
}

// DeleteUptimeMonitorInput represents the input of a DeleteUptimeMonitor operation.
type DeleteUptimeMonitorInput struct {
	_               struct{}
	UptimeMonitorID *int64
}

// DeleteUptimeMonitorOutput represents the output of a DeleteUptimeMonitor operation.
type DeleteUptimeMonitorOutput struct {
	_ struct{}
}

// DeleteUptimeMonitor deletes the specified uptime monitor. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors~1{id}/delete
func (c *Client) DeleteUptimeMonitor(input *DeleteUptimeMonitorInput) (*DeleteUptimeMonitorOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UptimeMonitorID == nil {
		return nil, errors.New("UptimeMonitor id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.uptimeMonitors, *input.UptimeMonitorID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUptimeMonitorOutput{}, nil
}

// GetUptimeMonitorsCountInput represents the input of a GetUptimeMonitorsCount operation.
type GetUptimeMonitorsCountInput struct {
	_ struct{}
}

// GetUptimeMonitorsCountOutput represents the output of a GetUptimeMonitorsCount operation.
type GetUptimeMonitorsCountOutput struct {
	_     struct{}
	Count int
}

// GetUptimeMonitorsCount gets list uptime monitors. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors~1count/get
func (c *Client) GetUptimeMonitorsCount(input *GetUptimeMonitorsCountInput) (*GetUptimeMonitorsCountOutput, error) {
	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/count", apiRoutes.uptimeMonitors))
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

	return &GetUptimeMonitorsCountOutput{Count: body.Count}, nil
}
