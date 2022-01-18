package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UptimeMonitor definition https://api.ilert.com/api-docs/#tag/Uptime-Monitors
type UptimeMonitor struct {
	ID                              int64                    `json:"id,omitempty"`
	Name                            string                   `json:"name"`
	Region                          string                   `json:"region"`
	CheckType                       string                   `json:"checkType"`
	CheckParams                     UptimeMonitorCheckParams `json:"checkParams,omitempty"`
	IntervalSec                     int                      `json:"intervalSec,omitempty"`                     // default: 300
	TimeoutMs                       int                      `json:"timeoutMs,omitempty"`                       // default: 30000
	CreateIncidentAfterFailedChecks int                      `json:"createIncidentAfterFailedChecks,omitempty"` // default: 1
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
	Down:    "Down",
	Warning: "warn",
	Paused:  "paused",
	Unknown: "unknown",
}

// UptimeMonitorRegions defines uptime monitor regions
var UptimeMonitorRegions = struct {
	EU string
	US string
}{
	EU: "EU",
	US: "US",
}

// UptimeMonitorCheckTypes defines uptime monitor check types
var UptimeMonitorCheckTypes = struct {
	HTTP string
	Ping string
	TCP  string
	UDP  string
}{
	HTTP: "http",
	Ping: "ping",
	TCP:  "tcp",
	UDP:  "udp",
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
}

// GetUptimeMonitorsOutput represents the output of a GetUptimeMonitors operation.
type GetUptimeMonitorsOutput struct {
	_              struct{}
	UptimeMonitors []*UptimeMonitor
}

// GetUptimeMonitors gets list uptime monitors. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors/get
func (c *Client) GetUptimeMonitors(input *GetUptimeMonitorsInput) (*GetUptimeMonitorsOutput, error) {
	resp, err := c.httpClient.R().Get(apiRoutes.uptimeMonitors)
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

// DeleteUptimeMonitor deletes the specified alert source. https://api.ilert.com/api-docs/#tag/Uptime-Monitors/paths/~1uptime-monitors~1{id}/delete
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
