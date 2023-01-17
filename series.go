package ilert

import (
	"errors"
	"fmt"
)

// SingleSeries definition https://api.ilert.com/api-docs/#tag/Series
type SingleSeries struct {
	Timestamp int64   `json:"timestamp,omitempty"`
	Value     float64 `json:"value"`
}

// MultipleSeries definition https://api.ilert.com/api-docs/#tag/Series
type MultipleSeries struct {
	Series []SingleSeries `json:"series"`
}

// CreateSingleSeriesInput represents the input of a CreateSingleSeries operation.
type CreateSingleSeriesInput struct {
	_         struct{}
	Series    *SingleSeries
	MetricKey *string
}

// CreateSingleSeries ingests a series for a metric. https://api.ilert.com/api-docs/#tag/Series/paths/~1series~1{key}/post
func (c *Client) CreateSingleSeries(input *CreateSingleSeriesInput) error {
	if input == nil {
		return errors.New("input is required")
	}
	if input.MetricKey == nil {
		return errors.New("metric integration key is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Series).Post(fmt.Sprintf("%s/%s", apiRoutes.series, *input.MetricKey))
	if err != nil {
		return err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return apiErr
	}

	return nil
}

// CreateMultipleSeriesInput represents the input of a CreateMultipleSeries operation.
type CreateMultipleSeriesInput struct {
	_         struct{}
	Series    *MultipleSeries
	MetricKey *string
}

// CreateMultipleSeries ingests multiple series for a metric. https://api.ilert.com/api-docs/#tag/Series/paths/~1series~1{key}/post
func (c *Client) CreateMultipleSeries(input *CreateMultipleSeriesInput) error {
	if input == nil {
		return errors.New("input is required")
	}
	if input.MetricKey == nil {
		return errors.New("metric integration key is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Series).Post(fmt.Sprintf("%s/%s", apiRoutes.series, *input.MetricKey))
	if err != nil {
		return err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return apiErr
	}

	return nil
}
