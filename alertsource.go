package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AlertSource definition
type AlertSource struct {
	ID                     int64                  `json:"id,omitempty"`
	Name                   string                 `json:"name"`
	IconURL                string                 `json:"iconUrl,omitempty"`
	LightIconURL           string                 `json:"lightIconUrl,omitempty"`
	DarkIconURL            string                 `json:"darkIconUrl,omitempty"`
	IntegrationType        string                 `json:"integrationType"`
	IntegrationKey         string                 `json:"integrationKey,omitempty"`
	IntegrationURL         string                 `json:"integrationUrl,omitempty"`
	IncidentCreation       string                 `json:"incidentCreation,omitempty"`
	EmailFiltered          bool                   `json:"emailFiltered,omitempty"`
	EmailResolveFiltered   bool                   `json:"emailResolveFiltered,omitempty"`
	Active                 bool                   `json:"active,omitempty"`
	Status                 string                 `json:"status,omitempty"`
	AutoResolutionTimeout  string                 `json:"autoResolutionTimeout,omitempty"` // e.g. PT4H
	EmailPredicates        []EmailPredicate       `json:"emailPredicates,omitempty"`
	EmailResolvePredicates []EmailPredicate       `json:"emailResolvePredicates,omitempty"`
	ResolveKeyExtractor    *EmailPredicate        `json:"resolveKeyExtractor,omitempty"`
	FilterOperator         string                 `json:"filterOperator,omitempty"`
	ResolveFilterOperator  string                 `json:"resolveFilterOperator,omitempty"`
	IncidentPriorityRule   string                 `json:"incidentPriorityRule,omitempty"`
	SupportHours           *SupportHours          `json:"supportHours,omitempty"`
	EscalationPolicy       *EscalationPolicy      `json:"escalationPolicy,omitempty"`
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
	AmazonCloudWatch           string
	API                        string
	AppDynamics                string
	Autotask                   string
	AWSBudget                  string
	AWSPersonalHealthDashboard string
	CallRoutingNumber          string
	checkmk                    string
	Datadog                    string
	Dynatrace                  string
	Email                      string
	Github                     string
	GoogleStackdriver          string
	Grafana                    string
	Heartbeat                  string
	Icinga                     string
	Instana                    string
	Jira                       string
	KentixAlarmManager         string
	Nagios                     string
	NewRelic                   string
	Pingdom                    string
	Prometheus                 string
	PRTGNetworkMonitor         string
	SMS                        string
	Solarwinds                 string
	StatusCake                 string
	TOPdesk                    string
	UptimeMonitor              string
	UPTIMEROBOT                string
	Zabbix                     string
}{
	AmazonCloudWatch:           "CLOUDWATCH",
	API:                        "API",
	AppDynamics:                "APPDYNAMICS",
	Autotask:                   "AUTOTASK",
	AWSBudget:                  "AWSBUDGET",
	AWSPersonalHealthDashboard: "AWSPHD",
	CallRoutingNumber:          "CRN",
	checkmk:                    "CHECKMK",
	Datadog:                    "DATADOG",
	Dynatrace:                  "DYNATRACE",
	Email:                      "EMAIL",
	Github:                     "GITHUB",
	GoogleStackdriver:          "STACKDRIVER",
	Grafana:                    "GRAFANA",
	Heartbeat:                  "HEARTBEAT",
	Icinga:                     "ICINGA",
	Instana:                    "INSTANA",
	Jira:                       "JIRA",
	KentixAlarmManager:         "KENTIXAM",
	Nagios:                     "NAGIOS",
	NewRelic:                   "NEWRELIC",
	Pingdom:                    "PINGDOM",
	Prometheus:                 "PROMETHEUS",
	PRTGNetworkMonitor:         "PRTG",
	SMS:                        "SMS",
	Solarwinds:                 "SOLARWINDS",
	StatusCake:                 "STATUSCAKE",
	TOPdesk:                    "TOPDESK",
	UptimeMonitor:              "MONITOR",
	UPTIMEROBOT:                "UPTIMEROBOT",
	Zabbix:                     "ZABBIX",
}

// AlertSourceIntegrationTypesAll defines all alert source integration types
var AlertSourceIntegrationTypesAll = []string{
	AlertSourceIntegrationTypes.AmazonCloudWatch,
	AlertSourceIntegrationTypes.API,
	AlertSourceIntegrationTypes.AppDynamics,
	AlertSourceIntegrationTypes.Autotask,
	AlertSourceIntegrationTypes.AWSBudget,
	AlertSourceIntegrationTypes.AWSPersonalHealthDashboard,
	AlertSourceIntegrationTypes.CallRoutingNumber,
	AlertSourceIntegrationTypes.checkmk,
	AlertSourceIntegrationTypes.Datadog,
	AlertSourceIntegrationTypes.Dynatrace,
	AlertSourceIntegrationTypes.Email,
	AlertSourceIntegrationTypes.Github,
	AlertSourceIntegrationTypes.GoogleStackdriver,
	AlertSourceIntegrationTypes.Grafana,
	AlertSourceIntegrationTypes.Heartbeat,
	AlertSourceIntegrationTypes.Icinga,
	AlertSourceIntegrationTypes.Instana,
	AlertSourceIntegrationTypes.Jira,
	AlertSourceIntegrationTypes.KentixAlarmManager,
	AlertSourceIntegrationTypes.Nagios,
	AlertSourceIntegrationTypes.NewRelic,
	AlertSourceIntegrationTypes.Pingdom,
	AlertSourceIntegrationTypes.Prometheus,
	AlertSourceIntegrationTypes.PRTGNetworkMonitor,
	AlertSourceIntegrationTypes.SMS,
	AlertSourceIntegrationTypes.Solarwinds,
	AlertSourceIntegrationTypes.StatusCake,
	AlertSourceIntegrationTypes.TOPdesk,
	AlertSourceIntegrationTypes.UptimeMonitor,
	AlertSourceIntegrationTypes.UPTIMEROBOT,
	AlertSourceIntegrationTypes.Zabbix,
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
	resp, err := c.httpClient.R().SetBody(input.AlertSource).Post(fmt.Sprintf("%s", apiRoutes.alertSources))
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

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.alertSources, *input.AlertSourceID))
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
	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s", apiRoutes.alertSources))
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

	resp, err := c.httpClient.R().SetBody(input.AlertSource).Put(fmt.Sprintf("%s/%d", apiRoutes.alertSources, *input.AlertSourceID))
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

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.alertSources, *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 204); err != nil {
		return nil, err
	}

	output := &DeleteAlertSourceOutput{}
	return output, nil
}
