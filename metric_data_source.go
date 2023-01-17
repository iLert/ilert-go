package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// MetricDataSource definition https://api.ilert.com/api-docs/#tag/Metric-Data-Sources
type MetricDataSource struct {
	ID       int64                     `json:"id,omitempty"`
	Name     string                    `json:"name"`
	Type     string                    `json:"type"`
	Teams    *TeamShort                `json:"teams,omitempty"`
	Metadata *MetricDataSourceMetadata `json:"metadata"`
}

// MetricDataSourceMetadata defines provider metadata for the metric data source
type MetricDataSourceMetadata struct {
	Region         string `json:"region,omitempty"`         // used for Datadog
	ApiKey         string `json:"apiKey,omitempty"`         // used for Datadog
	ApplicationKey string `json:"applicationKey,omitempty"` // used for Datadog
	AuthType       string `json:"authType,omitempty"`       // used for Prometheus
	BasicUser      string `json:"basicUser,omitempty"`      // used for Prometheus
	BasicPass      string `json:"basicPass,omitempty"`      // used for Prometheus
	HeaderKey      string `json:"headerKey,omitempty"`      // used for Prometheus
	HeaderValue    string `json:"headerValue,omitempty"`    // used for Prometheus
}

// MetricDataSourceType defines provider type of the metric data source
var MetricDataSourceType = struct {
	Datadog    string
	Prometheus string
}{
	Datadog:    "DATADOG",
	Prometheus: "PROMETHEUS",
}

// MetricDataSourceType defines provider type list
var MetricDataSourceAll = []string{
	MetricDataSourceType.Datadog,
	MetricDataSourceType.Prometheus,
}

// CreateMetricDataSourceInput represents the input of a CreateMetricDataSource operation.
type CreateMetricDataSourceInput struct {
	_                struct{}
	MetricDataSource *MetricDataSource
}

// CreateMetricDataSourceOutput represents the output of a CreateMetricDataSource operation.
type CreateMetricDataSourceOutput struct {
	_                struct{}
	MetricDataSource *MetricDataSource
}

// CreateMetricDataSource creates a new metric data source. https://api.ilert.com/api-docs/#tag/Metric-Data-Sources/paths/~1metric-data-sources/post
func (c *Client) CreateMetricDataSource(input *CreateMetricDataSourceInput) (*CreateMetricDataSourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricDataSource == nil {
		return nil, errors.New("metric data source input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.MetricDataSource).Post(apiRoutes.metricDataSources)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	metricdatasource := &MetricDataSource{}
	err = json.Unmarshal(resp.Body(), metricdatasource)
	if err != nil {
		return nil, err
	}

	return &CreateMetricDataSourceOutput{MetricDataSource: metricdatasource}, nil
}

// GetMetricDataSourcesInput represents the input of a GetMetricDataSources operation.
type GetMetricDataSourcesInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	// Default: 0
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 10, Maximum: 25 or 100 without include
	MaxResults *int
}

// GetMetricDataSourcesOutput represents the output of a GetMetricDataSources operation.
type GetMetricDataSourcesOutput struct {
	_                 struct{}
	MetricDataSources []*MetricDataSource
}

// GetMetricDataSources lists metricdatasource sources. https://api.ilert.com/api-docs/#tag/Metric-Data-Sources/paths/~1metric-data-sources/get
func (c *Client) GetMetricDataSources(input *GetMetricDataSourcesInput) (*GetMetricDataSourcesOutput, error) {
	if input == nil {
		input = &GetMetricDataSourcesInput{}
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

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.metricDataSources, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metricdatasources := make([]*MetricDataSource, 0)
	err = json.Unmarshal(resp.Body(), &metricdatasources)
	if err != nil {
		return nil, err
	}

	return &GetMetricDataSourcesOutput{MetricDataSources: metricdatasources}, nil
}

// GetMetricDataSourceInput represents the input of a GetMetricDataSource operation.
type GetMetricDataSourceInput struct {
	_                  struct{}
	MetricDataSourceID *int64
}

// GetMetricDataSourceOutput represents the output of a GetMetricDataSource operation.
type GetMetricDataSourceOutput struct {
	_                struct{}
	MetricDataSource *MetricDataSource
}

// GetMetricDataSource gets a metric data source by ID. https://api.ilert.com/api-docs/#tag/Metric-Data-Sources/paths/~1metric-data-sources~1{id}/get
func (c *Client) GetMetricDataSource(input *GetMetricDataSourceInput) (*GetMetricDataSourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricDataSourceID == nil {
		return nil, errors.New("metric data source id is required")
	}

	var url = fmt.Sprintf("%s/%d", apiRoutes.metricDataSources, *input.MetricDataSourceID)

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metricDataSource := &MetricDataSource{}
	err = json.Unmarshal(resp.Body(), metricDataSource)
	if err != nil {
		return nil, err
	}

	return &GetMetricDataSourceOutput{MetricDataSource: metricDataSource}, nil
}

// SearchMetricDataSourceInput represents the input of a SearchMetricDataSource operation.
type SearchMetricDataSourceInput struct {
	_                    struct{}
	MetricDataSourceName *string
}

// SearchMetricDataSourceOutput represents the output of a SearchMetricDataSource operation.
type SearchMetricDataSourceOutput struct {
	_                struct{}
	MetricDataSource *MetricDataSource
}

// SearchMetricDataSource gets the metric data source with specified name.
func (c *Client) SearchMetricDataSource(input *SearchMetricDataSourceInput) (*SearchMetricDataSourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricDataSourceName == nil {
		return nil, errors.New("metric data source name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.metricDataSources, *input.MetricDataSourceName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metricDataSource := &MetricDataSource{}
	err = json.Unmarshal(resp.Body(), metricDataSource)
	if err != nil {
		return nil, err
	}

	return &SearchMetricDataSourceOutput{MetricDataSource: metricDataSource}, nil
}

// UpdateMetricDataSourceInput represents the input of a UpdateMetricDataSource operation.
type UpdateMetricDataSourceInput struct {
	_                  struct{}
	MetricDataSourceID *int64
	MetricDataSource   *MetricDataSource
}

// UpdateMetricDataSourceOutput represents the output of a UpdateMetricDataSource operation.
type UpdateMetricDataSourceOutput struct {
	_                struct{}
	MetricDataSource *MetricDataSource
}

// UpdateMetricDataSource updates the specific metric data source. https://api.ilert.com/api-docs/#tag/Metric-Data-Sources/paths/~1metric-data-sources~1{id}/put
func (c *Client) UpdateMetricDataSource(input *UpdateMetricDataSourceInput) (*UpdateMetricDataSourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricDataSourceID == nil {
		return nil, errors.New("metric data source id is required")
	}
	if input.MetricDataSource == nil {
		return nil, errors.New("metric data source input is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.metricDataSources, *input.MetricDataSourceID)

	resp, err := c.httpClient.R().SetBody(input.MetricDataSource).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	metricDataSource := &MetricDataSource{}
	err = json.Unmarshal(resp.Body(), metricDataSource)
	if err != nil {
		return nil, err
	}

	return &UpdateMetricDataSourceOutput{MetricDataSource: metricDataSource}, nil
}

// DeleteMetricDataSourceInput represents the input of a DeleteMetricDataSource operation.
type DeleteMetricDataSourceInput struct {
	_                  struct{}
	MetricDataSourceID *int64
}

// DeleteMetricDataSourceOutput represents the output of a DeleteMetricDataSource operation.
type DeleteMetricDataSourceOutput struct {
	_ struct{}
}

// DeleteMetricDataSource deletes the specified metric data source. https://api.ilert.com/api-docs/#tag/Metric-Data-Sources/paths/~1metric-data-sources~1{id}/delete
func (c *Client) DeleteMetricDataSource(input *DeleteMetricDataSourceInput) (*DeleteMetricDataSourceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.MetricDataSourceID == nil {
		return nil, errors.New("metric data source id is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.metricDataSources, *input.MetricDataSourceID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteMetricDataSourceOutput{}, nil
}
