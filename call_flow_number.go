package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// CallFlowNumber definition https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-numbers
type CallFlowNumber struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`

	// whether the number is currently assigned to a call flow, read-only
	// possible values: "AVAILABLE", "USED"
	State string `json:"state,omitempty"`

	// the call flow this number is assigned to. Only returned by the list endpoint, only
	// when "assignedTo" is requested through Include, and only on a number whose state
	// is "USED". GetCallFlowNumber and SearchCallFlowNumber never return it.
	AssignedTo *CallFlowNumberAssignedTo `json:"assignedTo,omitempty"`
}

// CallFlowNumberAssignedTo references the call flow a number is assigned to
type CallFlowNumberAssignedTo struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
}

// CallFlowNumberState defines the states a call flow number is returned in
var CallFlowNumberState = struct {
	Available string
	Used      string
}{
	Available: "AVAILABLE",
	Used:      "USED",
}

// CallFlowNumberStateAll defines the call flow number state list
var CallFlowNumberStateAll = []string{
	CallFlowNumberState.Available,
	CallFlowNumberState.Used,
}

// CallFlowNumberStateFilter defines the values accepted by the state filter of a
// call flow number list. "ANY" is a filter value only, it is never returned on a number.
var CallFlowNumberStateFilter = struct {
	Any       string
	Available string
	Used      string
}{
	Any:       "ANY",
	Available: CallFlowNumberState.Available,
	Used:      CallFlowNumberState.Used,
}

// CallFlowNumberStateFilterAll defines the call flow number state filter list
var CallFlowNumberStateFilterAll = []string{
	CallFlowNumberStateFilter.Any,
	CallFlowNumberStateFilter.Available,
	CallFlowNumberStateFilter.Used,
}

// CallFlowNumberInclude defines the optional properties of a call flow number
var CallFlowNumberInclude = struct {
	AssignedTo string
}{
	AssignedTo: "assignedTo",
}

// CallFlowNumberIncludeAll defines the optional properties of a call flow number list
var CallFlowNumberIncludeAll = []string{
	CallFlowNumberInclude.AssignedTo,
}

// GetCallFlowNumbersInput represents the input of a GetCallFlowNumbers operation.
type GetCallFlowNumbersInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int

	// filters the result by availability state, "ANY" when left out
	// possible values: "ANY", "AVAILABLE", "USED"
	State *string

	// filters the result by a case-insensitive substring of the call flow number name.
	// The API applies it before paging.
	Query *string

	// describes optional properties that should be included in the response
	// possible values: "assignedTo"
	Include []*string
}

// GetCallFlowNumbersOutput represents the output of a GetCallFlowNumbers operation.
type GetCallFlowNumbersOutput struct {
	_               struct{}
	CallFlowNumbers []*CallFlowNumber
}

// GetCallFlowNumbers lists existing call flow numbers.
//
// The API reads at most 250 numbers per state bucket before it filters and pages, so an
// account with more than that does not see all of them through this endpoint.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-numbers
func (c *Client) GetCallFlowNumbers(input *GetCallFlowNumbersInput) (*GetCallFlowNumbersOutput, error) {
	if input == nil {
		input = &GetCallFlowNumbersInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}
	if input.State != nil {
		q.Add("state", *input.State)
	}
	if input.Query != nil {
		q.Add("query", *input.Query)
	}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.callFlowNumbers, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlowNumbers := make([]*CallFlowNumber, 0)
	err = json.Unmarshal(resp.Body(), &callFlowNumbers)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowNumbersOutput{CallFlowNumbers: callFlowNumbers}, nil
}

// GetCallFlowNumberInput represents the input of a GetCallFlowNumber operation.
type GetCallFlowNumberInput struct {
	_                struct{}
	CallFlowNumberID *int64
}

// GetCallFlowNumberOutput represents the output of a GetCallFlowNumber operation.
type GetCallFlowNumberOutput struct {
	_              struct{}
	CallFlowNumber *CallFlowNumber
}

// GetCallFlowNumber gets the call flow number with the specified id. https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-numbers
func (c *Client) GetCallFlowNumber(input *GetCallFlowNumberInput) (*GetCallFlowNumberOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowNumberID == nil {
		return nil, errors.New("call flow number id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.callFlowNumbers, *input.CallFlowNumberID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlowNumber := &CallFlowNumber{}
	err = json.Unmarshal(resp.Body(), callFlowNumber)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowNumberOutput{CallFlowNumber: callFlowNumber}, nil
}

// SearchCallFlowNumberInput represents the input of a SearchCallFlowNumber operation.
type SearchCallFlowNumberInput struct {
	_                  struct{}
	CallFlowNumberName *string
}

// SearchCallFlowNumberOutput represents the output of a SearchCallFlowNumber operation.
type SearchCallFlowNumberOutput struct {
	_              struct{}
	CallFlowNumber *CallFlowNumber
}

// SearchCallFlowNumber gets the call flow number with the specified name. The API matches
// the value against the phone number first and falls back to the name, so a number is found
// by either. A number matched by its phone number comes back without State.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-numbers
func (c *Client) SearchCallFlowNumber(input *SearchCallFlowNumberInput) (*SearchCallFlowNumberOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowNumberName == nil {
		return nil, errors.New("call flow number name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.callFlowNumbers, url.PathEscape(*input.CallFlowNumberName)))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlowNumber := &CallFlowNumber{}
	err = json.Unmarshal(resp.Body(), callFlowNumber)
	if err != nil {
		return nil, err
	}

	return &SearchCallFlowNumberOutput{CallFlowNumber: callFlowNumber}, nil
}
