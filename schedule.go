package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Schedule definition https://api.ilert.com/api-docs/#tag/Schedules
type Schedule struct {
	ID           int64       `json:"id,omitempty"`
	Name         string      `json:"name"`
	Timezone     string      `json:"timezone,omitempty"`
	StartsOn     string      `json:"startsOn,omitempty"` // Date time string in ISO format
	CurrentShift Shift       `json:"currentShift,omitempty"`
	NextShift    Shift       `json:"nextShift,omitempty"`
	Teams        []TeamShort `json:"teams,omitempty"`
}

// Shift definition
type Shift struct {
	User  User   `json:"user"`
	Start string `json:"start"` // Date time string in ISO format
	End   string `json:"end"`   // Date time string in ISO format
}

// GetScheduleInput represents the input of a GetSchedule operation.
type GetScheduleInput struct {
	_          struct{}
	ScheduleID *int64
}

// GetScheduleOutput represents the output of a GetSchedule operation.
type GetScheduleOutput struct {
	_        struct{}
	Schedule *Schedule
}

// GetSchedule gets the on-call schedule with the specified id. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}/get
func (c *Client) GetSchedule(input *GetScheduleInput) (*GetScheduleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.schedules, *input.ScheduleID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	schedule := &Schedule{}
	err = json.Unmarshal(resp.Body(), schedule)
	if err != nil {
		return nil, err
	}

	return &GetScheduleOutput{Schedule: schedule}, nil
}

// GetSchedulesInput represents the input of a GetSchedules operation.
type GetSchedulesInput struct {
	_ struct{}
}

// GetSchedulesOutput represents the output of a GetSchedules operation.
type GetSchedulesOutput struct {
	_         struct{}
	Schedules []*Schedule
}

// GetSchedules gets list on-call schedules. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules/get
func (c *Client) GetSchedules(input *GetSchedulesInput) (*GetSchedulesOutput, error) {
	resp, err := c.httpClient.R().Get(apiRoutes.schedules)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	schedules := make([]*Schedule, 0)
	err = json.Unmarshal(resp.Body(), &schedules)
	if err != nil {
		return nil, err
	}

	return &GetSchedulesOutput{Schedules: schedules}, nil
}

// GetScheduleShiftsInput represents the input of a GetScheduleShifts operation.
type GetScheduleShiftsInput struct {
	_                struct{}
	ScheduleID       *int64
	From             *string // Date time string in ISO format
	Until            *string // Date time string in ISO format
	ExcludeOverrides *bool
}

// GetScheduleShiftsOutput represents the output of a GetScheduleShifts operation.
type GetScheduleShiftsOutput struct {
	_      struct{}
	Shifts []*Shift
}

// GetScheduleShifts gets shifts for the specified schedule and date range. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}~1shifts/get
func (c *Client) GetScheduleShifts(input *GetScheduleShiftsInput) (*GetScheduleShiftsOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	q := url.Values{}
	if input.From != nil {
		q.Add("from", *input.From)
	}
	if input.Until != nil {
		q.Add("until", *input.From)
	}
	if input.ExcludeOverrides != nil {
		q.Add("exclude-overrides", strconv.FormatBool(*input.ExcludeOverrides))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/shifts?%s", apiRoutes.schedules, *input.ScheduleID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	shifts := make([]*Shift, 0)
	err = json.Unmarshal(resp.Body(), &shifts)
	if err != nil {
		return nil, err
	}

	return &GetScheduleShiftsOutput{Shifts: shifts}, nil
}

// GetScheduleOverridesInput represents the input of a GetScheduleOverrides operation.
type GetScheduleOverridesInput struct {
	_          struct{}
	ScheduleID *int64
}

// GetScheduleOverridesOutput represents the output of a GetScheduleOverrides operation.
type GetScheduleOverridesOutput struct {
	_         struct{}
	Overrides []*Shift
}

// GetScheduleOverrides gets overrides for the specified schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}~1overrides/get
func (c *Client) GetScheduleOverrides(input *GetScheduleOverridesInput) (*GetScheduleOverridesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/overrides", apiRoutes.schedules, *input.ScheduleID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	overrides := make([]*Shift, 0)
	err = json.Unmarshal(resp.Body(), &overrides)
	if err != nil {
		return nil, err
	}

	return &GetScheduleOverridesOutput{Overrides: overrides}, nil
}

// GetScheduleUserOnCallInput represents the input of a GetScheduleUserOnCall operation.
type GetScheduleUserOnCallInput struct {
	_          struct{}
	ScheduleID *int64
}

// GetScheduleUserOnCallOutput represents the output of a GetScheduleUserOnCall operation.
type GetScheduleUserOnCallOutput struct {
	_     struct{}
	Shift *Shift
}

// GetScheduleUserOnCall gets overrides for the specified schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}~1user-on-call/get
func (c *Client) GetScheduleUserOnCall(input *GetScheduleUserOnCallInput) (*GetScheduleUserOnCallOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/user-on-call", apiRoutes.schedules, *input.ScheduleID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200, 204); apiErr != nil {
		return nil, apiErr
	}

	if resp.StatusCode() == 204 {
		return &GetScheduleUserOnCallOutput{}, nil
	}

	shift := &Shift{}
	err = json.Unmarshal(resp.Body(), shift)
	if err != nil {
		return nil, err
	}

	return &GetScheduleUserOnCallOutput{Shift: shift}, nil
}
