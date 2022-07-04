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
	AlertCreation          string                 `json:"alertCreation,omitempty"`
	IncidentCreation       string                 `json:"incidentCreation,omitempty"` // @deprecated
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
	AlertPriorityRule      string                 `json:"alertPriorityRule,omitempty"`
	IncidentPriorityRule   string                 `json:"incidentPriorityRule,omitempty"` // @deprecated
	SupportHours           *SupportHours          `json:"supportHours,omitempty"`
	EscalationPolicy       *EscalationPolicy      `json:"escalationPolicy,omitempty"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
	AutotaskMetadata       *AutotaskMetadata      `json:"autotaskMetadata,omitempty"`
	Heartbeat              *Heartbeat             `json:"heartbeat,omitempty"`
	Teams                  []TeamShort            `json:"teams,omitempty"`
}

// EmailPredicate definition
type EmailPredicate struct {
	Field    string `json:"field"`
	Criteria string `json:"criteria"`
	Value    string `json:"value"`
}

// SupportHours definition
type SupportHours struct {
	Timezone           string      `json:"timezone"`
	AutoRaiseAlerts    bool        `json:"autoRaiseAlerts,omitempty"`    // Raise priority of all pending alerts for this alert source to 'high' when support hours begin
	AutoRaiseIncidents bool        `json:"autoRaiseIncidents,omitempty"` // @deprecated
	SupportDays        SupportDays `json:"supportDays"`
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

// AlertSourceAlertCreations defines alert source alert creations
var AlertSourceAlertCreations = struct {
	OneAlertPerEmail        string
	OneAlertPerEmailSubject string
	OnePendingAlertAllowed  string
	OneOpenAlertAllowed     string
	OpenResolveOnExtraction string
}{
	OneAlertPerEmail:        "ONE_INCIDENT_PER_EMAIL",
	OneAlertPerEmailSubject: "ONE_INCIDENT_PER_EMAIL_SUBJECT",
	OnePendingAlertAllowed:  "ONE_PENDING_INCIDENT_ALLOWED",
	OneOpenAlertAllowed:     "ONE_OPEN_INCIDENT_ALLOWED",
	OpenResolveOnExtraction: "OPEN_RESOLVE_ON_EXTRACTION",
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
	CheckMK                    string
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
	Consul                     string
	Zammad                     string
	SignalFx                   string
	Splunk                     string
	Kubernetes                 string
	Sematext                   string
	Sentry                     string
	Sumologic                  string
	Raygun                     string
	MXToolBox                  string
	ESWatcher                  string
	AmazonSNS                  string
	Kapacitor                  string
	CortexXSOAR                string
	Sysdig                     string
	ServerDensity              string
	Zapier                     string
	ServiceNow                 string
	SearchGuard                string
	AzureAlerts                string
	TerraformCloud             string
	Zendesk                    string
	Auvik                      string
	Sensu                      string
	NCentral                   string
	JumpCloud                  string
	Salesforce                 string
	GuardDuty                  string
	StatusHub                  string
	IXON                       string
	ApiFortress                string
	FreshService               string
	AppSignal                  string
	LightStep                  string
	IBMCloudFunctions          string
	CrowdStrike                string
	Humio                      string
	OhDear                     string
	MongodbAtlas               string
	Gitlab                     string
}{
	AmazonCloudWatch:           "CLOUDWATCH",
	API:                        "API",
	AppDynamics:                "APPDYNAMICS",
	Autotask:                   "AUTOTASK",
	AWSBudget:                  "AWSBUDGET",
	AWSPersonalHealthDashboard: "AWSPHD",
	CallRoutingNumber:          "CRN",
	CheckMK:                    "CHECKMK",
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
	Consul:                     "CONSUL",
	Zammad:                     "ZAMMAD",
	SignalFx:                   "SIGNALFX",
	Splunk:                     "SPLUNK",
	Kubernetes:                 "KUBERNETES",
	Sematext:                   "SEMATEXT",
	Sentry:                     "SENTRY",
	Sumologic:                  "SUMOLOGIC",
	Raygun:                     "RAYGUN",
	MXToolBox:                  "MXTOOLBOX",
	ESWatcher:                  "ESWATCHER",
	AmazonSNS:                  "AMAZONSNS",
	Kapacitor:                  "KAPACITOR",
	CortexXSOAR:                "CORTEXXSOAR",
	Sysdig:                     "SYSDIG",
	ServerDensity:              "SERVERDENSITY",
	Zapier:                     "ZAPIER",
	ServiceNow:                 "SERVICENOW",
	SearchGuard:                "SEARCHGUARD",
	AzureAlerts:                "AZUREALERTS",
	TerraformCloud:             "TERRAFORMCLOUD",
	Zendesk:                    "ZENDESK",
	Auvik:                      "AUVIK",
	Sensu:                      "SENSU",
	NCentral:                   "NCENTRAL",
	JumpCloud:                  "JUMPCLOUD",
	Salesforce:                 "SALESFORCE",
	GuardDuty:                  "GUARDDUTY",
	StatusHub:                  "STATUSHUB",
	IXON:                       "IXON",
	ApiFortress:                "APIFORTRESS",
	FreshService:               "FRESHSERVICE",
	AppSignal:                  "APPSIGNAL",
	LightStep:                  "LIGHTSTEP",
	IBMCloudFunctions:          "IBMCLOUDFUNCTIONS",
	CrowdStrike:                "CROWDSTRIKE",
	Humio:                      "HUMIO",
	OhDear:                     "OHDEAR",
	MongodbAtlas:               "MONGODBATLAS",
	Gitlab:                     "GITLAB",
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
	AlertSourceIntegrationTypes.CheckMK,
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
	AlertSourceIntegrationTypes.Consul,
	AlertSourceIntegrationTypes.Zammad,
	AlertSourceIntegrationTypes.SignalFx,
	AlertSourceIntegrationTypes.Splunk,
	AlertSourceIntegrationTypes.Kubernetes,
	AlertSourceIntegrationTypes.Sematext,
	AlertSourceIntegrationTypes.Sentry,
	AlertSourceIntegrationTypes.Sumologic,
	AlertSourceIntegrationTypes.Raygun,
	AlertSourceIntegrationTypes.MXToolBox,
	AlertSourceIntegrationTypes.ESWatcher,
	AlertSourceIntegrationTypes.AmazonSNS,
	AlertSourceIntegrationTypes.Kapacitor,
	AlertSourceIntegrationTypes.CortexXSOAR,
	AlertSourceIntegrationTypes.Sysdig,
	AlertSourceIntegrationTypes.ServerDensity,
	AlertSourceIntegrationTypes.Zapier,
	AlertSourceIntegrationTypes.ServiceNow,
	AlertSourceIntegrationTypes.SearchGuard,
	AlertSourceIntegrationTypes.AzureAlerts,
	AlertSourceIntegrationTypes.TerraformCloud,
	AlertSourceIntegrationTypes.Zendesk,
	AlertSourceIntegrationTypes.Auvik,
	AlertSourceIntegrationTypes.Sensu,
	AlertSourceIntegrationTypes.NCentral,
	AlertSourceIntegrationTypes.JumpCloud,
	AlertSourceIntegrationTypes.Salesforce,
	AlertSourceIntegrationTypes.GuardDuty,
	AlertSourceIntegrationTypes.StatusHub,
	AlertSourceIntegrationTypes.IXON,
	AlertSourceIntegrationTypes.ApiFortress,
	AlertSourceIntegrationTypes.FreshService,
	AlertSourceIntegrationTypes.AppSignal,
	AlertSourceIntegrationTypes.LightStep,
	AlertSourceIntegrationTypes.IBMCloudFunctions,
	AlertSourceIntegrationTypes.CrowdStrike,
	AlertSourceIntegrationTypes.Humio,
	AlertSourceIntegrationTypes.OhDear,
	AlertSourceIntegrationTypes.MongodbAtlas,
	AlertSourceIntegrationTypes.Gitlab,
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
		return nil, errors.New("input is required")
	}
	if input.AlertSource == nil {
		return nil, errors.New("alert source input is required")
	}

	if input.AlertSource.AlertCreation != "" && input.AlertSource.IncidentCreation != "" {
		input.AlertSource.IncidentCreation = ""
	}
	if input.AlertSource.AlertCreation == "" {
		input.AlertSource.AlertCreation = input.AlertSource.IncidentCreation
		input.AlertSource.IncidentCreation = ""
	}

	if input.AlertSource.AlertPriorityRule != "" && input.AlertSource.IncidentPriorityRule != "" {
		input.AlertSource.IncidentPriorityRule = ""
	}
	if input.AlertSource.AlertPriorityRule == "" {
		input.AlertSource.AlertPriorityRule = input.AlertSource.IncidentPriorityRule
		input.AlertSource.IncidentPriorityRule = ""
	}

	if input.AlertSource.SupportHours.AutoRaiseIncidents {
		input.AlertSource.SupportHours.AutoRaiseAlerts = true
		input.AlertSource.SupportHours.AutoRaiseIncidents = false
	}

	resp, err := c.httpClient.R().SetBody(input.AlertSource).Post(apiRoutes.alertSources)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	alertSource := &AlertSource{}
	err = json.Unmarshal(resp.Body(), alertSource)
	if err != nil {
		return nil, err
	}

	return &CreateAlertSourceOutput{AlertSource: alertSource}, nil
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
		return nil, errors.New("input is required")
	}
	if input.AlertSourceID == nil {
		return nil, errors.New("AlertSource id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.alertSources, *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertSource := &AlertSource{}
	err = json.Unmarshal(resp.Body(), alertSource)
	if err != nil {
		return nil, err
	}

	return &GetAlertSourceOutput{AlertSource: alertSource}, nil
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
	resp, err := c.httpClient.R().Get(apiRoutes.alertSources)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertSources := make([]*AlertSource, 0)
	err = json.Unmarshal(resp.Body(), &alertSources)
	if err != nil {
		return nil, err
	}

	return &GetAlertSourcesOutput{AlertSources: alertSources}, nil
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
		return nil, errors.New("input is required")
	}
	if input.AlertSource == nil {
		return nil, errors.New("AlertSource input is required")
	}
	if input.AlertSourceID == nil {
		return nil, errors.New("alert source id is required")
	}

	if input.AlertSource.AlertCreation != "" && input.AlertSource.IncidentCreation != "" {
		input.AlertSource.IncidentCreation = ""
	}
	if input.AlertSource.AlertCreation == "" {
		input.AlertSource.AlertCreation = input.AlertSource.IncidentCreation
		input.AlertSource.IncidentCreation = ""
	}

	if input.AlertSource.AlertPriorityRule != "" && input.AlertSource.IncidentPriorityRule != "" {
		input.AlertSource.IncidentPriorityRule = ""
	}
	if input.AlertSource.AlertPriorityRule == "" {
		input.AlertSource.AlertPriorityRule = input.AlertSource.IncidentPriorityRule
		input.AlertSource.IncidentPriorityRule = ""
	}

	if input.AlertSource.SupportHours.AutoRaiseIncidents {
		input.AlertSource.SupportHours.AutoRaiseAlerts = true
		input.AlertSource.SupportHours.AutoRaiseIncidents = false
	}

	resp, err := c.httpClient.R().SetBody(input.AlertSource).Put(fmt.Sprintf("%s/%d", apiRoutes.alertSources, *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	alertSource := &AlertSource{}
	err = json.Unmarshal(resp.Body(), alertSource)
	if err != nil {
		return nil, err
	}

	return &UpdateAlertSourceOutput{AlertSource: alertSource}, nil
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
		return nil, errors.New("input is required")
	}
	if input.AlertSourceID == nil {
		return nil, errors.New("AlertSource id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.alertSources, *input.AlertSourceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteAlertSourceOutput{}, nil
}
