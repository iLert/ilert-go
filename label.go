package ilert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// LabelKey defines a distinct label key in use on an account
type LabelKey struct {
	Key string `json:"key"`
}

// LabelValue defines a distinct value used for a label key on an account
type LabelValue struct {
	Value string `json:"value"`
}

// AlertLabelPageSize is the page size the alert label endpoints enforce. Unlike the
// service and telemetry source label endpoints, GetAlertLabelKeys and GetAlertLabelValues
// reject a MaxResults that is not exactly this value, and a StartIndex that is not a
// multiple of it, with a 400.
const AlertLabelPageSize = 100

// LabelMaxResults is the largest MaxResults the service and telemetry source label
// endpoints accept. A higher value is rejected with a 400.
const LabelMaxResults = 1000

// GetAlertLabelKeysInput represents the input of a GetAlertLabelKeys operation.
type GetAlertLabelKeysInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a
	// list of entities, has to be a multiple of AlertLabelPageSize
	StartIndex *int

	// the maximum number of results when paging through a list of entities,
	// has to be AlertLabelPageSize when it is set
	MaxResults *int

	// filters the result by a prefix of the label key
	Query *string
}

// GetAlertLabelKeysOutput represents the output of a GetAlertLabelKeys operation.
type GetAlertLabelKeysOutput struct {
	_         struct{}
	LabelKeys []*LabelKey
}

// GetAlertLabelKeys lists the distinct label keys used by the alerts of this account.
//
// Unlike its service and telemetry source counterparts this endpoint pages in fixed steps
// of AlertLabelPageSize, and it serves the first three pages of an unfiltered listing from
// a cache with a five minute TTL, so a label key of a freshly created alert can take that
// long to appear. Setting Query bypasses the cache (verified against the API on 16.09.2026).
//
// Callers with the guest or stakeholder role get an empty list: the API answers them with an
// empty object rather than an empty list, which is read as no labels (verified against the
// API on 25.09.2026).
// https://docs.ilert.com/developer-docs/rest-api/api-reference/alerts
func (c *Client) GetAlertLabelKeys(input *GetAlertLabelKeysInput) (*GetAlertLabelKeysOutput, error) {
	if input == nil {
		input = &GetAlertLabelKeysInput{}
	}

	labelKeys, err := c.getLabelKeys(apiRoutes.alerts, input.StartIndex, input.MaxResults, input.Query)
	if err != nil {
		return nil, err
	}

	return &GetAlertLabelKeysOutput{LabelKeys: labelKeys}, nil
}

// GetAlertLabelValuesInput represents the input of a GetAlertLabelValues operation.
type GetAlertLabelValuesInput struct {
	_ struct{}

	// the label key whose values are listed
	LabelKey *string

	// an integer specifying the starting point (beginning with 0) when paging through a
	// list of entities, has to be a multiple of AlertLabelPageSize
	StartIndex *int

	// the maximum number of results when paging through a list of entities,
	// has to be AlertLabelPageSize when it is set
	MaxResults *int

	// filters the result by a prefix of the label value
	Query *string
}

// GetAlertLabelValuesOutput represents the output of a GetAlertLabelValues operation.
type GetAlertLabelValuesOutput struct {
	_           struct{}
	LabelValues []*LabelValue
}

// GetAlertLabelValues lists the distinct values used for a label key by the alerts of this
// account. Pages in fixed steps of AlertLabelPageSize, and caches the first page of an
// unfiltered listing for two minutes, so a value of a freshly created alert can take that
// long to appear. Setting Query bypasses the cache (verified against the API on 16.09.2026).
// Like GetAlertLabelKeys it returns an empty list to callers with the guest or stakeholder role.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/alerts
func (c *Client) GetAlertLabelValues(input *GetAlertLabelValuesInput) (*GetAlertLabelValuesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	labelValues, err := c.getLabelValues(apiRoutes.alerts, input.LabelKey, input.StartIndex, input.MaxResults, input.Query)
	if err != nil {
		return nil, err
	}

	return &GetAlertLabelValuesOutput{LabelValues: labelValues}, nil
}

// GetServiceLabelKeysInput represents the input of a GetServiceLabelKeys operation.
type GetServiceLabelKeysInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: LabelMaxResults
	MaxResults *int

	// filters the result by a prefix of the label key
	Query *string
}

// GetServiceLabelKeysOutput represents the output of a GetServiceLabelKeys operation.
type GetServiceLabelKeysOutput struct {
	_         struct{}
	LabelKeys []*LabelKey
}

// GetServiceLabelKeys lists the distinct label keys used by the services of this account.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) GetServiceLabelKeys(input *GetServiceLabelKeysInput) (*GetServiceLabelKeysOutput, error) {
	if input == nil {
		input = &GetServiceLabelKeysInput{}
	}

	labelKeys, err := c.getLabelKeys(apiRoutes.services, input.StartIndex, input.MaxResults, input.Query)
	if err != nil {
		return nil, err
	}

	return &GetServiceLabelKeysOutput{LabelKeys: labelKeys}, nil
}

// GetServiceLabelValuesInput represents the input of a GetServiceLabelValues operation.
type GetServiceLabelValuesInput struct {
	_ struct{}

	// the label key whose values are listed
	LabelKey *string

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: LabelMaxResults
	MaxResults *int

	// filters the result by a prefix of the label value
	Query *string
}

// GetServiceLabelValuesOutput represents the output of a GetServiceLabelValues operation.
type GetServiceLabelValuesOutput struct {
	_           struct{}
	LabelValues []*LabelValue
}

// GetServiceLabelValues lists the distinct values used for a label key by the services of
// this account. https://docs.ilert.com/developer-docs/rest-api/api-reference/services
func (c *Client) GetServiceLabelValues(input *GetServiceLabelValuesInput) (*GetServiceLabelValuesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	labelValues, err := c.getLabelValues(apiRoutes.services, input.LabelKey, input.StartIndex, input.MaxResults, input.Query)
	if err != nil {
		return nil, err
	}

	return &GetServiceLabelValuesOutput{LabelValues: labelValues}, nil
}

// GetTelemetrySourceLabelKeysInput represents the input of a GetTelemetrySourceLabelKeys operation.
type GetTelemetrySourceLabelKeysInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: LabelMaxResults
	MaxResults *int

	// filters the result by a prefix of the label key
	Query *string
}

// GetTelemetrySourceLabelKeysOutput represents the output of a GetTelemetrySourceLabelKeys operation.
type GetTelemetrySourceLabelKeysOutput struct {
	_         struct{}
	LabelKeys []*LabelKey
}

// GetTelemetrySourceLabelKeys lists the distinct label keys used by the telemetry sources of
// this account. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) GetTelemetrySourceLabelKeys(input *GetTelemetrySourceLabelKeysInput) (*GetTelemetrySourceLabelKeysOutput, error) {
	if input == nil {
		input = &GetTelemetrySourceLabelKeysInput{}
	}

	labelKeys, err := c.getLabelKeys(apiRoutes.telemetrySources, input.StartIndex, input.MaxResults, input.Query)
	if err != nil {
		return nil, err
	}

	return &GetTelemetrySourceLabelKeysOutput{LabelKeys: labelKeys}, nil
}

// GetTelemetrySourceLabelValuesInput represents the input of a GetTelemetrySourceLabelValues operation.
type GetTelemetrySourceLabelValuesInput struct {
	_ struct{}

	// the label key whose values are listed
	LabelKey *string

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: LabelMaxResults
	MaxResults *int

	// filters the result by a prefix of the label value
	Query *string
}

// GetTelemetrySourceLabelValuesOutput represents the output of a GetTelemetrySourceLabelValues operation.
type GetTelemetrySourceLabelValuesOutput struct {
	_           struct{}
	LabelValues []*LabelValue
}

// GetTelemetrySourceLabelValues lists the distinct values used for a label key by the telemetry
// sources of this account. https://docs.ilert.com/developer-docs/rest-api/api-reference/telemetry-sources
func (c *Client) GetTelemetrySourceLabelValues(input *GetTelemetrySourceLabelValuesInput) (*GetTelemetrySourceLabelValuesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	labelValues, err := c.getLabelValues(apiRoutes.telemetrySources, input.LabelKey, input.StartIndex, input.MaxResults, input.Query)
	if err != nil {
		return nil, err
	}

	return &GetTelemetrySourceLabelValuesOutput{LabelValues: labelValues}, nil
}

func (c *Client) getLabelKeys(route string, startIndex *int, maxResults *int, query *string) ([]*LabelKey, error) {
	q := labelQuery(startIndex, maxResults, query)

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/labels?%s", route, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	labelKeys := make([]*LabelKey, 0)
	if isEmptyJSONObject(resp.Body()) {
		return labelKeys, nil
	}
	err = json.Unmarshal(resp.Body(), &labelKeys)
	if err != nil {
		return nil, err
	}

	return labelKeys, nil
}

func (c *Client) getLabelValues(route string, labelKey *string, startIndex *int, maxResults *int, query *string) ([]*LabelValue, error) {
	if labelKey == nil {
		return nil, errors.New("label key is required")
	}

	q := labelQuery(startIndex, maxResults, query)

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/labels/%s/values?%s", route, url.PathEscape(*labelKey), q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	labelValues := make([]*LabelValue, 0)
	if isEmptyJSONObject(resp.Body()) {
		return labelValues, nil
	}
	err = json.Unmarshal(resp.Body(), &labelValues)
	if err != nil {
		return nil, err
	}

	return labelValues, nil
}

// isEmptyJSONObject reports whether a label response is an empty object. The alert label
// endpoints answer callers with the guest or stakeholder role with "{}" where every other
// label endpoint answers with an empty list, so decoding it as a list would fail.
func isEmptyJSONObject(body []byte) bool {
	return bytes.Equal(bytes.TrimSpace(body), []byte("{}"))
}

// labelQuery builds the query of a label endpoint. Every parameter is only sent when the
// caller set it: the alert label endpoints accept no page size other than their fixed one,
// so a page size the SDK invented would turn every unparameterized call into a 400.
func labelQuery(startIndex *int, maxResults *int, query *string) url.Values {
	q := url.Values{}
	if startIndex != nil {
		q.Add("start-index", strconv.Itoa(*startIndex))
	}
	if maxResults != nil {
		q.Add("max-results", strconv.Itoa(*maxResults))
	}
	if query != nil {
		q.Add("query", *query)
	}

	return q
}
