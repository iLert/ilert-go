package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Schedule definition https://api.ilert.com/api-docs/#!/Schedules
type Schedule struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Timezone     string `json:"timezone"`
	StartsOn     string `json:"startsOn"` // Date time string in ISO format
	CurrentShift Shift  `json:"currentShift"`
	NextShift    Shift  `json:"nextShift"`
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
		return nil, errors.New("Input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("/api/v1/schedules/%d", *input.ScheduleID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	schedule := &Schedule{}
	err = json.Unmarshal(resp.Body(), schedule)
	if err != nil {
		return nil, err
	}

	output := &GetScheduleOutput{
		Schedule: schedule,
	}

	return output, nil
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
	resp, err := c.httpClient.R().Get("/api/v1/schedules")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	schedules := make([]*Schedule, 0)
	err = json.Unmarshal(resp.Body(), &schedules)
	if err != nil {
		return nil, err
	}

	output := &GetSchedulesOutput{Schedules: schedules}

	return output, nil
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
		return nil, errors.New("Input is required")
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

	resp, err := c.httpClient.R().Get(fmt.Sprintf("/api/v1/schedules/%d/shifts?%s", *input.ScheduleID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	shifts := make([]*Shift, 0)
	err = json.Unmarshal(resp.Body(), &shifts)
	if err != nil {
		return nil, err
	}

	output := &GetScheduleShiftsOutput{Shifts: shifts}

	return output, nil
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
		return nil, errors.New("Input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("/api/v1/schedules/%d/overrides", *input.ScheduleID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	overrides := make([]*Shift, 0)
	err = json.Unmarshal(resp.Body(), &overrides)
	if err != nil {
		return nil, err
	}

	output := &GetScheduleOverridesOutput{Overrides: overrides}

	return output, nil
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
		return nil, errors.New("Input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("Schedule id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("/api/v1/schedules/%d/user-on-call", *input.ScheduleID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200, 204); err != nil {
		return nil, err
	}

	output := &GetScheduleUserOnCallOutput{}
	if resp.StatusCode() == 204 {
		return output, nil
	}

	shift := &Shift{}
	err = json.Unmarshal(resp.Body(), shift)
	if err != nil {
		return nil, err
	}

	output.Shift = shift

	return output, nil
}
