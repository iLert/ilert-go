package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// TelemetrySource definition https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
type TelemetrySource struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// the telemetry protocol of this source, required on create and immutable afterwards
	// possible values: "OTEL"
	Type string `json:"type"`

	// whether ilert is currently receiving telemetry from this source, read-only
	// possible values: "PENDING", "RECEIVING", "STOPPED", "ERROR"
	Status string `json:"status,omitempty"`

	// the credential used to send telemetry to ilert, read-only. Omitted from list responses,
	// and from the detail response for users without update permission on the telemetry source.
	IntegrationKey string `json:"integrationKey,omitempty"`

	// prepended to the name of every service ilert creates from the telemetry of this source.
	// Leading and trailing whitespace is removed, no separator is inserted. Max 64 characters.
	ServiceNamePrefix string `json:"serviceNamePrefix,omitempty"`

	// appended to the name of every service ilert creates from the telemetry of this source.
	// Leading and trailing whitespace is removed, no separator is inserted. Max 64 characters.
	ServiceNameSuffix string `json:"serviceNameSuffix,omitempty"`

	// free-form key-value labels assigned to this telemetry source. A nil map disappears
	// from the payload and leaves the labels untouched, a non-nil empty map marshals to
	// "labels":{} and clears them.
	Labels *map[string]string `json:"labels,omitempty"`

	// only returned when "teams" is requested through Include. A nil slice disappears from
	// the payload and leaves the teams untouched, a non-nil empty slice marshals to
	// "teams":[] and clears them, matching AlertAction, CallFlow and EventFlow.
	Teams *[]TeamShort `json:"teams,omitempty"`

	ManagedBy *ManagedBy `json:"managedBy,omitempty"`
	CreatedAt string     `json:"createdAt,omitempty"`
	UpdatedAt string     `json:"updatedAt,omitempty"`
}

// ManagedBy is present when an entity is managed outside of the ilert web app,
// for example through Terraform, Pulumi or the API. Read-only.
type ManagedBy struct {
	// possible values: "API", "TERRAFORM", "PULUMI", "ILERT"
	Type string `json:"type,omitempty"`

	// identifier of the system managing this entity
	Source string `json:"source,omitempty"`
}

// ManagedByType defines the systems that can manage an entity outside of the ilert web app
var ManagedByType = struct {
	Api       string
	Terraform string
	Pulumi    string
	Ilert     string
}{
	Api:       "API",
	Terraform: "TERRAFORM",
	Pulumi:    "PULUMI",
	Ilert:     "ILERT",
}

// ManagedByTypeAll defines the managed by type list
var ManagedByTypeAll = []string{
	ManagedByType.Api,
	ManagedByType.Terraform,
	ManagedByType.Pulumi,
	ManagedByType.Ilert,
}

// TelemetrySourceStats defines the ingestion statistics of a telemetry source
type TelemetrySourceStats struct {
	// when the most recent telemetry was received from this source
	LastTraceAt string `json:"lastTraceAt,omitempty"`
}

// TelemetrySourceType defines telemetry source types
var TelemetrySourceType = struct {
	Otel string
}{
	Otel: "OTEL",
}

// TelemetrySourceTypeAll defines telemetry source types list
var TelemetrySourceTypeAll = []string{
	TelemetrySourceType.Otel,
}

// TelemetrySourceStatus defines telemetry source statuses
var TelemetrySourceStatus = struct {
	Pending   string
	Receiving string
	Stopped   string
	Error     string
}{
	Pending:   "PENDING",
	Receiving: "RECEIVING",
	Stopped:   "STOPPED",
	Error:     "ERROR",
}

// TelemetrySourceStatusAll defines telemetry source statuses list
var TelemetrySourceStatusAll = []string{
	TelemetrySourceStatus.Pending,
	TelemetrySourceStatus.Receiving,
	TelemetrySourceStatus.Stopped,
	TelemetrySourceStatus.Error,
}

// TelemetrySourceInclude defines the optional properties of a telemetry source
var TelemetrySourceInclude = struct {
	Teams string
}{
	Teams: "teams",
}

// TelemetrySourceIncludeAll defines the optional properties of a telemetry source list
var TelemetrySourceIncludeAll = []string{
	TelemetrySourceInclude.Teams,
}

// CreateTelemetrySourceInput represents the input of a CreateTelemetrySource operation.
type CreateTelemetrySourceInput struct {
	_               struct{}
	TelemetrySource *TelemetrySource

	// describes optional properties that should be included in the response
	// possible values: "teams"
	Include []*string
}

// CreateTelemetrySourceOutput represents the output of a CreateTelemetrySource operation.
type CreateTelemetrySourceOutput struct {
	_               struct{}
	TelemetrySource *TelemetrySource
}

// CreateTelemetrySource creates a new telemetry source. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) CreateTelemetrySource(input *CreateTelemetrySourceInput) (*CreateTelemetrySourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySource == nil {
		return nil, errors.New("telemetry source input is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.TelemetrySource).Post(fmt.Sprintf("%s?%s", apiRoutes.telemetrySources, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	telemetrySource := &TelemetrySource{}
	err = json.Unmarshal(resp.Body(), telemetrySource)
	if err != nil {
		return nil, err
	}

	return &CreateTelemetrySourceOutput{TelemetrySource: telemetrySource}, nil
}

// GetTelemetrySourceInput represents the input of a GetTelemetrySource operation.
type GetTelemetrySourceInput struct {
	_                 struct{}
	TelemetrySourceID *int64

	// describes optional properties that should be included in the response
	// possible values: "teams"
	Include []*string
}

// GetTelemetrySourceOutput represents the output of a GetTelemetrySource operation.
type GetTelemetrySourceOutput struct {
	_               struct{}
	TelemetrySource *TelemetrySource
}

// GetTelemetrySource gets the telemetry source with the specified id. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) GetTelemetrySource(input *GetTelemetrySourceInput) (*GetTelemetrySourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySourceID == nil {
		return nil, errors.New("telemetry source id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.telemetrySources, *input.TelemetrySourceID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	telemetrySource := &TelemetrySource{}
	err = json.Unmarshal(resp.Body(), telemetrySource)
	if err != nil {
		return nil, err
	}

	return &GetTelemetrySourceOutput{TelemetrySource: telemetrySource}, nil
}

// GetTelemetrySourcesInput represents the input of a GetTelemetrySources operation.
type GetTelemetrySourcesInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int

	// describes optional properties that should be included in the response
	// possible values: "teams"
	Include []*string

	// filters the result by a search string matched against the telemetry source name
	Query *string

	// filters the result by labels, each entry in the form "key:value"
	Labels []*string
}

// GetTelemetrySourcesOutput represents the output of a GetTelemetrySources operation.
type GetTelemetrySourcesOutput struct {
	_                struct{}
	TelemetrySources []*TelemetrySource
}

// GetTelemetrySources lists existing telemetry sources. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) GetTelemetrySources(input *GetTelemetrySourcesInput) (*GetTelemetrySourcesOutput, error) {
	if input == nil {
		input = &GetTelemetrySourcesInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}
	if input.Query != nil {
		q.Add("query", *input.Query)
	}
	for _, include := range input.Include {
		q.Add("include", *include)
	}
	for _, label := range input.Labels {
		q.Add("labels", *label)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.telemetrySources, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	telemetrySources := make([]*TelemetrySource, 0)
	err = json.Unmarshal(resp.Body(), &telemetrySources)
	if err != nil {
		return nil, err
	}

	return &GetTelemetrySourcesOutput{TelemetrySources: telemetrySources}, nil
}

// SearchTelemetrySourceInput represents the input of a SearchTelemetrySource operation.
type SearchTelemetrySourceInput struct {
	_                   struct{}
	TelemetrySourceName *string

	// describes optional properties that should be included in the response
	// possible values: "teams"
	Include []*string
}

// SearchTelemetrySourceOutput represents the output of a SearchTelemetrySource operation.
type SearchTelemetrySourceOutput struct {
	_               struct{}
	TelemetrySource *TelemetrySource
}

// SearchTelemetrySource gets the telemetry source with the specified name. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) SearchTelemetrySource(input *SearchTelemetrySourceInput) (*SearchTelemetrySourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySourceName == nil {
		return nil, errors.New("telemetry source name is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s?%s", apiRoutes.telemetrySources, url.PathEscape(*input.TelemetrySourceName), q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	telemetrySource := &TelemetrySource{}
	err = json.Unmarshal(resp.Body(), telemetrySource)
	if err != nil {
		return nil, err
	}

	return &SearchTelemetrySourceOutput{TelemetrySource: telemetrySource}, nil
}

// UpdateTelemetrySourceInput represents the input of an UpdateTelemetrySource operation.
type UpdateTelemetrySourceInput struct {
	_                 struct{}
	TelemetrySourceID *int64
	TelemetrySource   *TelemetrySource

	// describes optional properties that should be included in the response
	// possible values: "teams"
	Include []*string
}

// UpdateTelemetrySourceOutput represents the output of an UpdateTelemetrySource operation.
type UpdateTelemetrySourceOutput struct {
	_               struct{}
	TelemetrySource *TelemetrySource
}

// UpdateTelemetrySource updates an existing telemetry source. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) UpdateTelemetrySource(input *UpdateTelemetrySourceInput) (*UpdateTelemetrySourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySource == nil {
		return nil, errors.New("telemetry source input is required")
	}
	if input.TelemetrySourceID == nil {
		return nil, errors.New("telemetry source id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.TelemetrySource).Put(fmt.Sprintf("%s/%d?%s", apiRoutes.telemetrySources, *input.TelemetrySourceID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	telemetrySource := &TelemetrySource{}
	err = json.Unmarshal(resp.Body(), telemetrySource)
	if err != nil {
		return nil, err
	}

	return &UpdateTelemetrySourceOutput{TelemetrySource: telemetrySource}, nil
}

// DeleteTelemetrySourceInput represents the input of a DeleteTelemetrySource operation.
type DeleteTelemetrySourceInput struct {
	_                 struct{}
	TelemetrySourceID *int64
}

// DeleteTelemetrySourceOutput represents the output of a DeleteTelemetrySource operation.
type DeleteTelemetrySourceOutput struct {
	_ struct{}
}

// DeleteTelemetrySource deletes the specified telemetry source. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) DeleteTelemetrySource(input *DeleteTelemetrySourceInput) (*DeleteTelemetrySourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySourceID == nil {
		return nil, errors.New("telemetry source id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.telemetrySources, *input.TelemetrySourceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteTelemetrySourceOutput{}, nil
}

// RotateTelemetrySourceIntegrationKeyInput represents the input of a RotateTelemetrySourceIntegrationKey operation.
type RotateTelemetrySourceIntegrationKeyInput struct {
	_                 struct{}
	TelemetrySourceID *int64
}

// RotateTelemetrySourceIntegrationKeyOutput represents the output of a RotateTelemetrySourceIntegrationKey operation.
type RotateTelemetrySourceIntegrationKeyOutput struct {
	_ struct{}
}

// RotateTelemetrySourceIntegrationKey rotates the integration key of a telemetry source. The new key is not
// returned by this endpoint, read the telemetry source again to obtain it.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) RotateTelemetrySourceIntegrationKey(input *RotateTelemetrySourceIntegrationKeyInput) (*RotateTelemetrySourceIntegrationKeyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySourceID == nil {
		return nil, errors.New("telemetry source id is required")
	}

	resp, err := c.httpClient.R().Put(fmt.Sprintf("%s/%d/rotate-integration-key", apiRoutes.telemetrySources, *input.TelemetrySourceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &RotateTelemetrySourceIntegrationKeyOutput{}, nil
}

// GetTelemetrySourceStatsInput represents the input of a GetTelemetrySourceStats operation.
type GetTelemetrySourceStatsInput struct {
	_                 struct{}
	TelemetrySourceID *int64
}

// GetTelemetrySourceStatsOutput represents the output of a GetTelemetrySourceStats operation.
type GetTelemetrySourceStatsOutput struct {
	_                    struct{}
	TelemetrySourceStats *TelemetrySourceStats
}

// GetTelemetrySourceStats gets the ingestion statistics of a telemetry source. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) GetTelemetrySourceStats(input *GetTelemetrySourceStatsInput) (*GetTelemetrySourceStatsOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelemetrySourceID == nil {
		return nil, errors.New("telemetry source id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/stats", apiRoutes.telemetrySources, *input.TelemetrySourceID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	stats := &TelemetrySourceStats{}
	err = json.Unmarshal(resp.Body(), stats)
	if err != nil {
		return nil, err
	}

	return &GetTelemetrySourceStatsOutput{TelemetrySourceStats: stats}, nil
}
