package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AlertSource definition
type AlertSource struct {
	ID                     int64                  `json:"id"`
	Name                   string                 `json:"name"`
	IconURL                string                 `json:"iconUrl"`
	LightIconURL           string                 `json:"lightIconUrl"`
	DarkIconURL            string                 `json:"darkIconUrl"`
	IntegrationType        string                 `json:"integrationType"`
	IntegrationKey         string                 `json:"integrationKey"`
	IncidentCreation       string                 `json:"incidentCreation"`
	EmailFiltered          bool                   `json:"emailFiltered"`
	EmailResolveFiltered   bool                   `json:"emailResolveFiltered"`
	Active                 bool                   `json:"active"`
	Status                 string                 `json:"status"`
	EmailPredicates        []EmailPredicate       `json:"emailPredicates"`
	EmailResolvePredicates []EmailPredicate       `json:"emailResolvePredicates"`
	ResolveKeyExtractor    *EmailPredicate        `json:"resolveKeyExtractor"`
	FilterOperator         string                 `json:"filterOperator"`
	ResolveFilterOperator  string                 `json:"resolveFilterOperator"`
	IncidentPriorityRule   string                 `json:"incidentPriorityRule"`
	SupportHours           *SupportHours          `json:"supportHours"`
	EscalationPolicy       *EscalationPolicy      `json:"escalationPolicy"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
	AutotaskMetadata       *AutotaskMetadata      `json:"autotaskMetadata,omitempty"`
	Heartbeat              *Heartbeat             `json:"heartbeat,omitempty"`
}

// EmailPredicate definition
type EmailPredicate struct {
	Field    string `json:"field"`
	Criteria string `json:"criteria"`
	Value    string `json:"value"`
}

// SupportHours definition
type SupportHours struct {
	Timezone    string      `json:"timezone"`
	SupportDays SupportDays `json:"supportDays"`
}

// SupportDays definition
type SupportDays struct {
	MONDAY    *SupportDay `json:"MONDAY"`
	TUESDAY   *SupportDay `json:"TUESDAY"`
	WEDNESDAY *SupportDay `json:"WEDNESDAY"`
	THURSDAY  *SupportDay `json:"THURSDAY"`
	FRIDAY    *SupportDay `json:"FRIDAY"`
	SATURDAY  *SupportDay `json:"SATURDAY"`
	SUNDAY    *SupportDay `json:"SUNDAY"`
}

// SupportDay definition
type SupportDay struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// AutotaskMetadata definition
type AutotaskMetadata struct {
	Username  string `json:"userName"`
	Secret    string `json:"secret"`
	WebServer string `json:"webServer"`
}

// Heartbeat definition
type Heartbeat struct {
	Summary     string `json:"summary"`
	IntervalSec int    `json:"intervalSec"`
	Status      string `json:"status"`
}

// AlertSourceStatuses defines alert source statuses
var AlertSourceStatuses = struct {
	Pending       string
	AllAccepted   string
	AllResolved   string
	InMaintenance string
	Disabled      string
}{
	Pending:       "PENDING",
	AllAccepted:   "ALL_ACCEPTED",
	AllResolved:   "ALL_RESOLVED",
	InMaintenance: "IN_MAINTENANCE",
	Disabled:      "DISABLED",
}

// AlertSourceIncidentCreations defines alert source incident creations
var AlertSourceIncidentCreations = struct {
	OneIncidentPerEmail        string
	OneIncidentPerEmailSubject string
	OnePendingIncidentAllowed  string
	OneOpenIncidentAllowed     string
	OpenResolveOnExtraction    string
}{
	OneIncidentPerEmail:        "ONE_INCIDENT_PER_EMAIL",
	OneIncidentPerEmailSubject: "ONE_INCIDENT_PER_EMAIL_SUBJECT",
	OnePendingIncidentAllowed:  "ONE_PENDING_INCIDENT_ALLOWED",
	OneOpenIncidentAllowed:     "ONE_OPEN_INCIDENT_ALLOWED",
	OpenResolveOnExtraction:    "OPEN_RESOLVE_ON_EXTRACTION",
}

// AlertSourceIntegrationTypes defines alert source integration types
var AlertSourceIntegrationTypes = struct {
	API                        string
	AppDynamics                string
	Autotask                   string
	AWSBudget                  string
	AWSPersonalHealthDashboard string
	checkmk                    string
	AmazonCloudWatch           string
	CallRoutingNumber          string
	Datadog                    string
	Dynatrace                  string
	Email                      string
	Github                     string
	Grafana                    string
	Heartbeat                  string
	Icinga                     string
	Instana                    string
	KentixAlarmManager         string
	UptimeMonitor              string
	Nagios                     string
	NewRelic                   string
	Pingdom                    string
	Prometheus                 string
	PRTGNetworkMonitor         string
	SMS                        string
	Solarwinds                 string
	GoogleStackdriver          string
	StatusCake                 string
	TOPdesk                    string
	UPTIMEROBOT                string
	Zabbix                     string
}{
	API:                        "API",
	AppDynamics:                "APPDYNAMICS",
	Autotask:                   "AUTOTASK",
	AWSBudget:                  "AWSBUDGET",
	AWSPersonalHealthDashboard: "AWSPHD",
	checkmk:                    "CHECKMK",
	AmazonCloudWatch:           "CLOUDWATCH",
	CallRoutingNumber:          "CRN",
	Datadog:                    "DATADOG",
	Dynatrace:                  "DYNATRACE",
	Email:                      "EMAIL",
	Github:                     "GITHUB",
	Grafana:                    "GRAFANA",
	Heartbeat:                  "HEARTBEAT",
	Icinga:                     "ICINGA",
	Instana:                    "INSTANA",
	KentixAlarmManager:         "KENTIXAM",
	UptimeMonitor:              "MONITOR",
	Nagios:                     "NAGIOS",
	NewRelic:                   "NEWRELIC",
	Pingdom:                    "PINGDOM",
	Prometheus:                 "PROMETHEUS",
	PRTGNetworkMonitor:         "PRTG",
	SMS:                        "SMS",
	Solarwinds:                 "SOLARWINDS",
	GoogleStackdriver:          "STACKDRIVER",
	StatusCake:                 "STATUSCAKE",
	TOPdesk:                    "TOPDESK",
	UPTIMEROBOT:                "UPTIMEROBOT",
	Zabbix:                     "ZABBIX",
}

// CreateAlertSourceInput represents the input of a CreateAlertSource operation.
type CreateAlertSourceInput struct {
	_           struct{}
	AlertSource *AlertSource
}

// CreateAlertSourceOutput represents the output of a CreateAlertSource operation.
type CreateAlertSourceOutput struct {
	_           struct{}
	AlertSource *AlertSource
}

// CreateAlertSource creates a new alert source. https://api.ilert.com/api-docs/#tag/Alert-Sources/paths/~1alert-sources/post
func (c *Client) CreateAlertSource(input *CreateAlertSourceInput) (*CreateAlertSourceOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.AlertSource == nil {
		return nil, errors.New("Alert source input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.AlertSource).Post("/api/v1/alert-sources")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 201); err != nil {
		return nil, err
	}

	alertSource := &AlertSource{}
	err = json.Unmarshal(resp.Body(), alertSource)
	if err != nil {
		return nil, err
	}

	output := &CreateAlertSourceOutput{AlertSource: alertSource}

	return output, nil
}

// GetAlertSourceInput represents the input of a GetAlertSource operation.
type GetAlertSourceInput struct {
	_             struct{}
	AlertSourceID *int64
}

// GetAlertSourceOutput represents the output of a GetAlertSource operation.
type GetAlertSourceOutput struct {
	_           struct{}
	AlertSource *AlertSource
}

// GetAlertSource gets the alert source with specified id. https://api.ilert.com/api-docs/#tag/Alert-Sources/paths/~1alert-sources~1{id}/get
func (c *Client) GetAlertSource(input *GetAlertSourceInput) (*GetAlertSourceOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.AlertSourceID == nil {
		return nil, errors.New("AlertSource id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("/api/v1/alert-sources/%d", *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	alertSource := &AlertSource{}
	err = json.Unmarshal(resp.Body(), alertSource)
	if err != nil {
		return nil, err
	}

	output := &GetAlertSourceOutput{
		AlertSource: alertSource,
	}

	return output, nil
}

// GetAlertSourcesInput represents the input of a GetAlertSources operation.
type GetAlertSourcesInput struct {
	_ struct{}
}

// GetAlertSourcesOutput represents the output of a GetAlertSources operation.
type GetAlertSourcesOutput struct {
	_            struct{}
	AlertSources []*AlertSource
}

// GetAlertSources lists alert sources. https://api.ilert.com/api-docs/#tag/Alert-Sources/paths/~1alert-sources/get
func (c *Client) GetAlertSources(input *GetAlertSourcesInput) (*GetAlertSourcesOutput, error) {
	resp, err := c.httpClient.R().Get("/api/v1/alert-sources")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	alertSources := make([]*AlertSource, 0)
	err = json.Unmarshal(resp.Body(), &alertSources)
	if err != nil {
		return nil, err
	}

	output := &GetAlertSourcesOutput{AlertSources: alertSources}

	return output, nil
}

// UpdateAlertSourceInput represents the input of a UpdateAlertSource operation.
type UpdateAlertSourceInput struct {
	_             struct{}
	AlertSourceID *int64
	AlertSource   *AlertSource
}

// UpdateAlertSourceOutput represents the output of a UpdateAlertSource operation.
type UpdateAlertSourceOutput struct {
	_           struct{}
	AlertSource *AlertSource
}

// UpdateAlertSource updates an existing alert source. https://api.ilert.com/api-docs/#tag/Alert-Sources/paths/~1alert-sources~1{id}/put
func (c *Client) UpdateAlertSource(input *UpdateAlertSourceInput) (*UpdateAlertSourceOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.AlertSource == nil {
		return nil, errors.New("AlertSource input is required")
	}
	if input.AlertSourceID == nil {
		return nil, errors.New("Alert source id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.AlertSource).Put(fmt.Sprintf("/api/v1/alert-sources/%d", *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	alertSource := &AlertSource{}
	err = json.Unmarshal(resp.Body(), alertSource)
	if err != nil {
		return nil, err
	}

	output := &UpdateAlertSourceOutput{AlertSource: alertSource}

	return output, nil
}

// DeleteAlertSourceInput represents the input of a DeleteAlertSource operation.
type DeleteAlertSourceInput struct {
	_             struct{}
	AlertSourceID *int64
}

// DeleteAlertSourceOutput represents the output of a DeleteAlertSource operation.
type DeleteAlertSourceOutput struct {
	_ struct{}
}

// DeleteAlertSource deletes the specified alert source. https://api.ilert.com/api-docs/#tag/Alert-Sources/paths/~1alert-sources~1{id}/delete
func (c *Client) DeleteAlertSource(input *DeleteAlertSourceInput) (*DeleteAlertSourceOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.AlertSourceID == nil {
		return nil, errors.New("AlertSource id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("/api/v1/alert-sources/%d", *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 204); err != nil {
		return nil, err
	}

	output := &DeleteAlertSourceOutput{}
	return output, nil
}
