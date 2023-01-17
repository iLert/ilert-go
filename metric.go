package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Metric definition https://api.ilert.com/api-docs/#tag/Metrics
type Metric struct {
	ID                    int64                   `json:"id,omitempty"`
	Name                  string                  `json:"name"`
	Description           string                  `json:"description,omitempty"`
	AggregationType       string                  `json:"aggregationType"`
	DisplayType           string                  `json:"displayType"`
	InterpolateGaps       bool                    `json:"interpolateGaps,omitempty"`
	LockYAxisMax          float64                 `json:"lockYAxisMax,omitempty"`
	LockYAxisMin          float64                 `json:"lockYAxisMin,omitempty"`
	MouseOverDecimal      float64                 `json:"mouseOverDecimal,omitempty"`
	ShowValuesOnMouseOver bool                    `json:"showValuesOnMouseOver,omitempty"`
	Teams                 *TeamShort              `json:"teams,omitempty"`
	UnitLabel             string                  `json:"unitLabel,omitempty"`
	Metadata              *MetricProviderMetadata `json:"metadata,omitempty"`
	DataSource            *MetricDataSource       `json:"dataSource,omitempty"`
}

// MetricProviderMetadata defines provider metadata for the metric
type MetricProviderMetadata struct {
	Query string `json:"query,omitempty"` // used for Datadog, Prometheus
}

// MetricAggregationType defines aggregation type for the metric
var MetricAggregationType = struct {
	Average string
	Sum     string
	Minimum string
	Maximum string
	Last    string
}{
	Average: "AVG",
	Sum:     "SUM",
	Minimum: "MIN",
	Maximum: "MAX",
	Last:    "LAST",
}

// MetricAggregationType defines aggregation type list
var MetricAggregationTypeAll = []string{
	MetricAggregationType.Average,
	MetricAggregationType.Sum,
	MetricAggregationType.Minimum,
	MetricAggregationType.Maximum,
	MetricAggregationType.Last,
}

// MetricDisplayType defines display type for the metric
var MetricDisplayType = struct {
	Graph  string
	Single string
}{
	Graph:  "GRAPH",
	Single: "SINGLE",
}

// MetricDisplayType defines display type list
var MetricDisplayTypeAll = []string{
	MetricDisplayType.Graph,
	MetricDisplayType.Single,
}

// CreateMetricInput represents the input of a CreateMetric operation.
type CreateMetricInput struct {
	_      struct{}
	Metric *Metric
}

// CreateMetricOutput represents the output of a CreateMetric operation.
type CreateMetricOutput struct {
	_      struct{}
	Metric *Metric
}

// CreateMetric creates a new metric. https://api.ilert.com/api-docs/#tag/Metrics/paths/~1metrics/post
func (c *Client) CreateMetric(input *CreateMetricInput) (*CreateMetricOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Metric == nil {
		return nil, errors.New("metric input is required")
	}
	if input.Metric.Metadata != nil && input.Metric.DataSource == nil {
		return nil, errors.New("data source id is required when setting provider metadata")
	}
	if input.Metric.DataSource != nil && input.Metric.Metadata == nil {
		return nil, errors.New("provider metadata is required when setting metric data source")
	}
	resp, err := c.httpClient.R().SetBody(input.Metric).Post(apiRoutes.metrics)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	metric := &Metric{}
	err = json.Unmarshal(resp.Body(), metric)
	if err != nil {
		return nil, err
	}

	return &CreateMetricOutput{Metric: metric}, nil
}

// GetMetricsInput represents the input of a GetMetrics operation.
type GetMetricsInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	// Default: 0
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 10, Maximum: 25 or 100 without include
	MaxResults *int

	// describes optional properties that should be included in the response
	// possible values: "dataSource", "integrationKey"
	Include []*string
}

// GetMetricsOutput represents the output of a GetMetrics operation.
type GetMetricsOutput struct {
	_       struct{}
	Metrics []*Metric
}

// GetMetrics lists metric sources. https://api.ilert.com/api-docs/#tag/Metrics/paths/~1metrics/get
func (c *Client) GetMetrics(input *GetMetricsInput) (*GetMetricsOutput, error) {
	if input == nil {
		input = &GetMetricsInput{}
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

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.metrics, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metrics := make([]*Metric, 0)
	err = json.Unmarshal(resp.Body(), &metrics)
	if err != nil {
		return nil, err
	}

	return &GetMetricsOutput{Metrics: metrics}, nil
}

// GetMetricInput represents the input of a GetMetric operation.
type GetMetricInput struct {
	_        struct{}
	MetricID *int64
}

// GetMetricOutput represents the output of a GetMetric operation.
type GetMetricOutput struct {
	_      struct{}
	Metric *Metric
}

// GetMetric gets a metric by ID. https://api.ilert.com/api-docs/#tag/Metrics/paths/~1metrics~1{id}/get
func (c *Client) GetMetric(input *GetMetricInput) (*GetMetricOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricID == nil {
		return nil, errors.New("metric id is required")
	}

	var url = fmt.Sprintf("%s/%d", apiRoutes.metrics, *input.MetricID)

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metric := &Metric{}
	err = json.Unmarshal(resp.Body(), metric)
	if err != nil {
		return nil, err
	}

	return &GetMetricOutput{Metric: metric}, nil
}

// SearchMetricInput represents the input of a SearchMetric operation.
type SearchMetricInput struct {
	_          struct{}
	MetricName *string
}

// SearchMetricOutput represents the output of a SearchMetric operation.
type SearchMetricOutput struct {
	_      struct{}
	Metric *Metric
}

// SearchMetric gets the metric with specified name.
func (c *Client) SearchMetric(input *SearchMetricInput) (*SearchMetricOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricName == nil {
		return nil, errors.New("metric name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.metrics, *input.MetricName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metric := &Metric{}
	err = json.Unmarshal(resp.Body(), metric)
	if err != nil {
		return nil, err
	}

	return &SearchMetricOutput{Metric: metric}, nil
}

// UpdateMetricInput represents the input of a UpdateMetric operation.
type UpdateMetricInput struct {
	_        struct{}
	MetricID *int64
	Metric   *Metric
}

// UpdateMetricOutput represents the output of a UpdateMetric operation.
type UpdateMetricOutput struct {
	_      struct{}
	Metric *Metric
}

// UpdateMetric updates the specific metric. https://api.ilert.com/api-docs/#tag/Metrics/paths/~1metrics~1{id}/put
func (c *Client) UpdateMetric(input *UpdateMetricInput) (*UpdateMetricOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricID == nil {
		return nil, errors.New("metric id is required")
	}
	if input.Metric == nil {
		return nil, errors.New("metric input is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.metrics, *input.MetricID)

	resp, err := c.httpClient.R().SetBody(input.Metric).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metric := &Metric{}
	err = json.Unmarshal(resp.Body(), metric)
	if err != nil {
		return nil, err
	}

	return &UpdateMetricOutput{Metric: metric}, nil
}

// DeleteMetricInput represents the input of a DeleteMetric operation.
type DeleteMetricInput struct {
	_        struct{}
	MetricID *int64
}

// DeleteMetricOutput represents the output of a DeleteMetric operation.
type DeleteMetricOutput struct {
	_ struct{}
}

// DeleteMetric deletes the specified metric. https://api.ilert.com/api-docs/#tag/Metrics/paths/~1metrics~1{id}/delete
func (c *Client) DeleteMetric(input *DeleteMetricInput) (*DeleteMetricOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricID == nil {
		return nil, errors.New("metric id is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.metrics, *input.MetricID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteMetricOutput{}, nil
}
