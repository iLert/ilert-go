package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// SupportHour definition https://api.ilert.com/api-docs/#tag/SupportHours
//
// A write has to carry the weekly coverage in SupportDays or SupportWindows, the API rejects a
// support hour with neither with a 400.
type SupportHour struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Teams       []TeamShort            `json:"teams"`
	Timezone    string                 `json:"timezone,omitempty"`
	SupportDays *SupportDays           `json:"supportDays"`
	Exceptions  []SupportHourException `json:"exceptions,omitempty"`

	// the complete weekly coverage, see SupportWindow. Returned whenever SupportDays, one window
	// per day, cannot describe the coverage, and nil otherwise. SupportDays is lossy then: each day
	// reports its first window, and a window running to midnight ends at 23:59. Sending back
	// exactly the SupportDays that was read keeps the stored coverage, a different one is rejected
	// with a 400.
	//
	// When sent it is the whole coverage: a nil pointer disappears from the payload, a pointer to an
	// empty slice marshals to "supportWindows":[] and clears the coverage. Leave SupportDays nil when
	// sending it: the API rejects the two with a 400 when they describe different coverage, which a
	// SupportDays read before changing the windows does. A write takes at most 84 windows.
	SupportWindows *[]SupportWindow `json:"supportWindows,omitempty"`
}

// SupportWindow is one window of the weekly coverage of a support hour, covering [From, To) in
// the timezone of the support hour. A To at or before From wraps past Sunday midnight, so FRIDAY
// 17:00 to MONDAY 09:00 is one window, and From equal to To is the whole week. Times are whole
// minutes (HH:mm). The API merges overlapping and touching windows and returns them in ascending
// order of From, the whole week as MONDAY 00:00 to MONDAY 00:00. Same shape as LayerRestriction.
//
// Only a support hour of its own has windows: the inline SupportHours of an alert source stays one
// window per day, and the API rejects windows there.
type SupportWindow struct {
	From *TimeOfWeek `json:"from"`
	To   *TimeOfWeek `json:"to"`
}

// SupportHourException definition
type SupportHourException struct {
	Name          string `json:"name,omitempty"`
	Start         string `json:"start"`                   // Date time string in ISO format
	End           string `json:"end"`                     // Date time string in ISO format
	SupportStatus string `json:"supportStatus,omitempty"` // one of SupportStatus
}

// SupportStatus defines support status
var SupportStatus = struct {
	None    string
	During  string
	Outside string
}{
	None:    "NONE",
	During:  "DURING",
	Outside: "OUTSIDE",
}

// SupportStatusAll defines support status list
var SupportStatusAll = []string{
	SupportStatus.None,
	SupportStatus.During,
	SupportStatus.Outside,
}

// CreateSupportHourInput represents the input of a CreateSupportHour operation.
type CreateSupportHourInput struct {
	_           struct{}
	SupportHour *SupportHour
}

// CreateSupportHourOutput represents the output of a CreateSupportHour operation.
type CreateSupportHourOutput struct {
	_           struct{}
	SupportHour *SupportHour
}

// CreateSupportHour creates a new support hours resource. https://api.ilert.com/api-docs/#tag/Support-Hours/paths/~1support-hours/post
func (c *Client) CreateSupportHour(input *CreateSupportHourInput) (*CreateSupportHourOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.SupportHour == nil {
		return nil, errors.New("support hour input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.SupportHour).Post(apiRoutes.supportHours)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	supportHour := &SupportHour{}
	err = json.Unmarshal(resp.Body(), supportHour)
	if err != nil {
		return nil, err
	}

	return &CreateSupportHourOutput{SupportHour: supportHour}, nil
}

// GetSupportHourInput represents the input of a GetSupportHour operation.
type GetSupportHourInput struct {
	_             struct{}
	SupportHourID *int64
}

// GetSupportHourOutput represents the output of a GetSupportHour operation.
type GetSupportHourOutput struct {
	_           struct{}
	SupportHour *SupportHour
}

// GetSupportHour gets the support hours resource with specified id. https://api.ilert.com/api-docs/#tag/Support-Hours/paths/~1support-hours~1{id}/get
func (c *Client) GetSupportHour(input *GetSupportHourInput) (*GetSupportHourOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.SupportHourID == nil {
		return nil, errors.New("support hour id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.supportHours, *input.SupportHourID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	supportHour := &SupportHour{}
	err = json.Unmarshal(resp.Body(), supportHour)
	if err != nil {
		return nil, err
	}

	return &GetSupportHourOutput{SupportHour: supportHour}, nil
}

// GetSupportHoursInput represents the input of a GetSupportHours operation.
type GetSupportHoursInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetSupportHoursOutput represents the output of a GetSupportHours operation.
type GetSupportHoursOutput struct {
	_            struct{}
	SupportHours []*SupportHour
}

// GetSupportHours lists existing support hours resources. https://api.ilert.com/api-docs/#tag/Support-Hours/paths/~1support-hours/get
func (c *Client) GetSupportHours(input *GetSupportHoursInput) (*GetSupportHoursOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.supportHours, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	supportHours := make([]*SupportHour, 0)
	err = json.Unmarshal(resp.Body(), &supportHours)
	if err != nil {
		return nil, err
	}

	return &GetSupportHoursOutput{SupportHours: supportHours}, nil
}

// SearchSupportHourInput represents the input of a SearchSupportHour operation.
type SearchSupportHourInput struct {
	_               struct{}
	SupportHourName *string
}

// SearchSupportHourOutput represents the output of a SearchSupportHour operation.
type SearchSupportHourOutput struct {
	_           struct{}
	SupportHour *SupportHour
}

// SearchSupportHour gets the support hours resource with specified name.
func (c *Client) SearchSupportHour(input *SearchSupportHourInput) (*SearchSupportHourOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.SupportHourName == nil {
		return nil, errors.New("support hour name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.supportHours, *input.SupportHourName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	supportHour := &SupportHour{}
	err = json.Unmarshal(resp.Body(), supportHour)
	if err != nil {
		return nil, err
	}

	return &SearchSupportHourOutput{SupportHour: supportHour}, nil
}

// UpdateSupportHourInput represents the input of a UpdateSupportHour operation.
type UpdateSupportHourInput struct {
	_             struct{}
	SupportHourID *int64
	SupportHour   *SupportHour
}

// UpdateSupportHourOutput represents the output of a UpdateSupportHour operation.
type UpdateSupportHourOutput struct {
	_           struct{}
	SupportHour *SupportHour
}

// UpdateSupportHour updates an existing support hours resource. https://api.ilert.com/api-docs/#tag/Support-Hours/paths/~1support-hours~1{id}/put
func (c *Client) UpdateSupportHour(input *UpdateSupportHourInput) (*UpdateSupportHourOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.SupportHour == nil {
		return nil, errors.New("support hour input is required")
	}
	if input.SupportHourID == nil {
		return nil, errors.New("support hour id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.SupportHour).Put(fmt.Sprintf("%s/%d", apiRoutes.supportHours, *input.SupportHourID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	supportHour := &SupportHour{}
	err = json.Unmarshal(resp.Body(), supportHour)
	if err != nil {
		return nil, err
	}

	return &UpdateSupportHourOutput{SupportHour: supportHour}, nil
}

// DeleteSupportHourInput represents the input of a DeleteSupportHour operation.
type DeleteSupportHourInput struct {
	_             struct{}
	SupportHourID *int64
}

// DeleteSupportHourOutput represents the output of a DeleteSupportHour operation.
type DeleteSupportHourOutput struct {
	_ struct{}
}

// DeleteSupportHour deletes the specified support hours resource. https://api.ilert.com/api-docs/#tag/Support-Hours/paths/~1support-hours~1{id}/delete
func (c *Client) DeleteSupportHour(input *DeleteSupportHourInput) (*DeleteSupportHourOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.SupportHourID == nil {
		return nil, errors.New("support hour id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.supportHours, *input.SupportHourID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteSupportHourOutput{}, nil
}
