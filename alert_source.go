package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
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
	IncidentCreation       string                 `json:"incidentCreation,omitempty"`     // @deprecated
	EmailFiltered          bool                   `json:"emailFiltered,omitempty"`        // @deprecated
	EmailResolveFiltered   bool                   `json:"emailResolveFiltered,omitempty"` // @deprecated
	Active                 bool                   `json:"active,omitempty"`
	Status                 string                 `json:"status,omitempty"`
	AutoResolutionTimeout  string                 `json:"autoResolutionTimeout,omitempty"`  // e.g. PT4H
	EmailPredicates        []EmailPredicate       `json:"emailPredicates,omitempty"`        // @deprecated
	EmailResolvePredicates []EmailPredicate       `json:"emailResolvePredicates,omitempty"` // @deprecated
	ResolveKeyExtractor    *EmailPredicate        `json:"resolveKeyExtractor,omitempty"`    // @deprecated
	FilterOperator         string                 `json:"filterOperator,omitempty"`         // @deprecated
	ResolveFilterOperator  string                 `json:"resolveFilterOperator,omitempty"`  // @deprecated
	AlertPriorityRule      string                 `json:"alertPriorityRule,omitempty"`
	IncidentPriorityRule   string                 `json:"incidentPriorityRule,omitempty"` // @deprecated
	SupportHours           interface{}            `json:"supportHours,omitempty"`
	EscalationPolicy       *EscalationPolicy      `json:"escalationPolicy,omitempty"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`         // @deprecated
	AutotaskMetadata       *AutotaskMetadata      `json:"autotaskMetadata,omitempty"` // @deprecated
	Heartbeat              *Heartbeat             `json:"heartbeat,omitempty"`        // @deprecated
	Teams                  []TeamShort            `json:"teams,omitempty"`
	SummaryTemplate        *Template              `json:"summaryTemplate,omitempty"`
	DetailsTemplate        *Template              `json:"detailsTemplate,omitempty"`
	RoutingTemplate        *Template              `json:"routingTemplate,omitempty"`
	LinkTemplates          []LinkTemplate         `json:"linkTemplates,omitempty"`
	PriorityTemplate       *PriorityTemplate      `json:"priorityTemplate,omitempty"`
	AlertGroupingWindow    string                 `json:"alertGroupingWindow,omitempty"` // e.g. PT4H
	ScoreThreshold         float64                `json:"scoreThreshold,omitempty"`
	EventFilter            string                 `json:"eventFilter,omitempty"`
}

// @deprecated EmailPredicate definition
type EmailPredicate struct {
	Field    string `json:"field"`
	Criteria string `json:"criteria"`
	Value    string `json:"value"`
}

// @deprecated SupportHours definition
type SupportHours struct {
	Timezone           string      `json:"timezone"`
	AutoRaiseAlerts    bool        `json:"autoRaiseAlerts,omitempty"`    // Raise priority of all pending alerts for this alert source to 'high' when support hours begin
	AutoRaiseIncidents bool        `json:"autoRaiseIncidents,omitempty"` // @deprecated
	SupportDays        SupportDays `json:"supportDays"`
}

func (s *SupportHours) RemoveLegacyFields() {
	if s.AutoRaiseIncidents {
		s.AutoRaiseAlerts = true
		s.AutoRaiseIncidents = false
	}
}

// SupportHoursReference definition
type SupportHoursReference struct {
	ID int64 `json:"id"`
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

// @deprecated Heartbeat definition
type Heartbeat struct {
	Summary     string `json:"summary"`
	IntervalSec int    `json:"intervalSec"`
	Status      string `json:"status"`
}

// Template definition
type Template struct {
	TextTemplate string `json:"textTemplate,omitempty"`
}

// LinkTemplate definition
type LinkTemplate struct {
	Text         string    `json:"text,omitempty"`
	HrefTemplate *Template `json:"hrefTemplate,omitempty"`
}

// PriorityTemplate definition
type PriorityTemplate struct {
	ValueTemplate *Template `json:"valueTemplate,omitempty"`
	Mappings      []Mapping `json:"mappings,omitempty"`
}

// Mapping definition
type Mapping struct {
	Value    string `json:"value"`
	Priority string `json:"priority"`
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

// AlertSourceStatusesAll defines alert source statuses list
var AlertSourceStatusesAll = []string{
	AlertSourceStatuses.Pending,
	AlertSourceStatuses.AllAccepted,
	AlertSourceStatuses.AllResolved,
	AlertSourceStatuses.InMaintenance,
	AlertSourceStatuses.Disabled,
}

// AlertSourceAlertCreations defines alert source alert creations
var AlertSourceAlertCreations = struct {
	// @deprecated
	OneIncidentPerEmail        string
	OneIncidentPerEmailSubject string
	OnePendingIncidentAllowed  string
	OneOpenIncidentAllowed     string

	OneAlertPerEmail         string
	OneAlertPerEmailSubject  string
	OnePendingAlertAllowed   string
	OneOpenAlertAllowed      string
	OpenResolveOnExtraction  string
	OneAlertGroupedPerWindow string
	IntelligentGrouping      string
}{
	// @deprecated
	OneIncidentPerEmail:        "ONE_INCIDENT_PER_EMAIL",
	OneIncidentPerEmailSubject: "ONE_INCIDENT_PER_EMAIL_SUBJECT",
	OnePendingIncidentAllowed:  "ONE_PENDING_INCIDENT_ALLOWED",
	OneOpenIncidentAllowed:     "ONE_OPEN_INCIDENT_ALLOWED",

	OneAlertPerEmail:         "ONE_ALERT_PER_EMAIL",
	OneAlertPerEmailSubject:  "ONE_ALERT_PER_EMAIL_SUBJECT",
	OnePendingAlertAllowed:   "ONE_PENDING_ALERT_ALLOWED",
	OneOpenAlertAllowed:      "ONE_OPEN_ALERT_ALLOWED",
	OpenResolveOnExtraction:  "OPEN_RESOLVE_ON_EXTRACTION",
	OneAlertGroupedPerWindow: "ONE_ALERT_GROUPED_PER_WINDOW",
	IntelligentGrouping:      "INTELLIGENT_GROUPING",
}

// AlertSourceAlertCreationsAll defines alert source alert creations list
var AlertSourceAlertCreationsAll = []string{
	// @deprecated
	AlertSourceAlertCreations.OneIncidentPerEmail,
	AlertSourceAlertCreations.OneIncidentPerEmailSubject,
	AlertSourceAlertCreations.OnePendingIncidentAllowed,
	AlertSourceAlertCreations.OneOpenIncidentAllowed,

	AlertSourceAlertCreations.OneAlertPerEmail,
	AlertSourceAlertCreations.OneAlertPerEmailSubject,
	AlertSourceAlertCreations.OnePendingAlertAllowed,
	AlertSourceAlertCreations.OneOpenAlertAllowed,
	AlertSourceAlertCreations.OpenResolveOnExtraction,
	AlertSourceAlertCreations.OneAlertGroupedPerWindow,
	AlertSourceAlertCreations.IntelligentGrouping,
}

// AlertSourceAlertGroupingWindows defines alert source alert grouping windows
var AlertSourceAlertGroupingWindows = struct {
	TwoMinutes      string
	FiveMinutes     string
	FifteenMinutes  string
	ThirtyMinutes   string
	OneHour         string
	TwoHours        string
	FourHours       string
	EightHours      string
	TwelveHours     string
	TwentyFourHours string
}{
	TwoMinutes:      "PT2M",
	FiveMinutes:     "PT5M",
	FifteenMinutes:  "PT15M",
	ThirtyMinutes:   "PT30M",
	OneHour:         "PT1H",
	TwoHours:        "PT2H",
	FourHours:       "PT4H",
	EightHours:      "PT8H",
	TwelveHours:     "PT12H",
	TwentyFourHours: "PT24H",
}

// AlertSourceAlertGroupingWindowsAll defines alert source alert grouping windows list
var AlertSourceAlertGroupingWindowsAll = []string{
	AlertSourceAlertGroupingWindows.TwoMinutes,
	AlertSourceAlertGroupingWindows.FiveMinutes,
	AlertSourceAlertGroupingWindows.FifteenMinutes,
	AlertSourceAlertGroupingWindows.ThirtyMinutes,
	AlertSourceAlertGroupingWindows.OneHour,
	AlertSourceAlertGroupingWindows.TwoHours,
	AlertSourceAlertGroupingWindows.FourHours,
	AlertSourceAlertGroupingWindows.EightHours,
	AlertSourceAlertGroupingWindows.TwelveHours,
	AlertSourceAlertGroupingWindows.TwentyFourHours,
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
	Checkly                    string
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
	Checkly:                    "CHECKLY",
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
	AlertSourceIntegrationTypes.Checkly,
}

// CreateAlertSourceInput represents the input of a CreateAlertSource operation.
type CreateAlertSourceInput struct {
	_           struct{}
	AlertSource *AlertSource

	// describes optional properties that should be included in the response
	// possible values: "summaryTemplate", "detailsTemplate", "routingTemplate", "textTemplate", "linkTemplates", "priorityTemplate", "eventFilter"
	Include []*string
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

	if v, ok := input.AlertSource.SupportHours.(SupportHours); ok {
		v.RemoveLegacyFields()
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.AlertSource).Post(fmt.Sprintf("%s?%s", apiRoutes.alertSources, q.Encode()))
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

	// describes optional properties that should be included in the response
	// possible values: "summaryTemplate", "detailsTemplate", "routingTemplate", "textTemplate", "linkTemplates", "priorityTemplate", "eventFilter"
	Include []*string
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
		return nil, errors.New("alert source id is required")
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.alertSources, *input.AlertSourceID, q.Encode()))
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

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 50
	MaxResults *int
}

// GetAlertSourcesOutput represents the output of a GetAlertSources operation.
type GetAlertSourcesOutput struct {
	_            struct{}
	AlertSources []*AlertSource
}

// GetAlertSources lists existing alert sources. https://api.ilert.com/api-docs/#tag/Alert-Sources/paths/~1alert-sources/get
func (c *Client) GetAlertSources(input *GetAlertSourcesInput) (*GetAlertSourcesOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.alertSources, q.Encode()))
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

// SearchAlertSourceInput represents the input of a SearchAlertSource operation.
type SearchAlertSourceInput struct {
	_               struct{}
	AlertSourceName *string
}

// SearchAlertSourceOutput represents the output of a SearchAlertSource operation.
type SearchAlertSourceOutput struct {
	_           struct{}
	AlertSource *AlertSource
}

// SearchAlertSource gets the alert source with specified name.
func (c *Client) SearchAlertSource(input *SearchAlertSourceInput) (*SearchAlertSourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AlertSourceName == nil {
		return nil, errors.New("alert source name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.alertSources, *input.AlertSourceName))
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

	return &SearchAlertSourceOutput{AlertSource: alertSource}, nil
}

// UpdateAlertSourceInput represents the input of a UpdateAlertSource operation.
type UpdateAlertSourceInput struct {
	_             struct{}
	AlertSourceID *int64
	AlertSource   *AlertSource

	// describes optional properties that should be included in the response
	// possible values: "summaryTemplate", "detailsTemplate", "routingTemplate", "textTemplate", "linkTemplates", "priorityTemplate", "eventFilter"
	Include []*string
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
		return nil, errors.New("alert source input is required")
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

	if v, ok := input.AlertSource.SupportHours.(SupportHours); ok {
		v.RemoveLegacyFields()
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.AlertSource).Put(fmt.Sprintf("%s/%d?%s", apiRoutes.alertSources, *input.AlertSourceID, q.Encode()))
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
		return nil, errors.New("alert source id is required")
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
