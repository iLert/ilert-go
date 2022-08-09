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
	ID                   int64           `json:"id,omitempty"`
	Name                 string          `json:"name"`
	Timezone             string          `json:"timezone,omitempty"`
	Type                 string          `json:"type,omitempty"`
	StartsOn             string          `json:"startsOn,omitempty"` // Date time string in ISO format, @deprecated
	ScheduleLayers       []ScheduleLayer `json:"scheduleLayers,omitempty"`
	Shifts               []Shift         `json:"shifts,omitempty"`
	ShowGaps             bool            `json:"showGaps,omitempty"`
	DefaultShiftDuration string          `json:"defaultShiftDuration,omitempty"` // P7D
	CurrentShift         *Shift          `json:"currentShift,omitempty"`
	NextShift            *Shift          `json:"nextShift,omitempty"`
	Teams                []TeamShort     `json:"teams,omitempty"`
}

// Shift definition
type Shift struct {
	User  User   `json:"user"`
	Start string `json:"start"` // Date time string in ISO format
	End   string `json:"end"`   // Date time string in ISO format
}

// Schedule layer definition
type ScheduleLayer struct {
	Name            string             `json:"name"`
	StartsOn        string             `json:"startsOn"`         // Date time string in ISO format
	EndsOn          string             `json:"endsOn,omitempty"` // Date time string in ISO format
	Users           []User             `json:"users"`
	Rotation        string             `json:"rotation"` // P7D
	RestrictionType string             `json:"restrictionType,omitempty"`
	Restrictions    []LayerRestriction `json:"restrictions,omitempty"`
}

type LayerRestriction struct {
	From *TimeOfWeek `json:"from"`
	To   *TimeOfWeek `json:"to"`
}

type TimeOfWeek struct {
	DayOfWeek string `json:"dayOfWeek"`
	Time      string `json:"time"` // Time string in format <15:00>
}

var ScheduleType = struct {
	Static    string
	Recurring string
}{
	Static:    "STATIC",
	Recurring: "RECURRING",
}

var RestrictionType = struct {
	TimeOfWeek string
	TimeOfDay  string
}{
	TimeOfWeek: "TIME_OF_WEEK",
	TimeOfDay:  "TIME_OF_DAY",
}

var DayOfWeek = struct {
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
}{
	Monday:    "MONDAY",
	Tuesday:   "TUESDAY",
	Wednesday: "WEDNESDAY",
	Thursday:  "THURSDAY",
	Friday:    "FRIDAY",
	Saturday:  "SATURDAY",
	Sunday:    "SUNDAY",
}

// CreateScheduleInput represents the input of a CreateSchedule operation.
type CreateScheduleInput struct {
	_           struct{}
	Schedule    *Schedule
	AbortOnGaps *bool
}

// CreateScheduleOutput represents the output of a CreateSchedule operation.
type CreateScheduleOutput struct {
	_        struct{}
	Schedule *Schedule
}

// CreateSchedule creates a new schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules/post
func (c *Client) CreateSchedule(input *CreateScheduleInput) (*CreateScheduleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Schedule == nil {
		return nil, errors.New("schedule input is required")
	}

	q := url.Values{}

	if input.AbortOnGaps != nil {
		q.Add("abort-on-gaps", strconv.FormatBool(*input.AbortOnGaps))
	}

	resp, err := c.httpClient.R().SetBody(input.Schedule).Post(fmt.Sprintf("%s?%s", apiRoutes.schedules, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	schedule := &Schedule{}
	err = json.Unmarshal(resp.Body(), schedule)
	if err != nil {
		return nil, err
	}

	return &CreateScheduleOutput{Schedule: schedule}, nil
}

// GetScheduleInput represents the input of a GetSchedule operation.
type GetScheduleInput struct {
	_          struct{}
	ScheduleID *int64

	// describes optional properties that should be included in the response
	Include []*string
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

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.schedules, *input.ScheduleID, q.Encode()))
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

	// describes optional properties that should be included in the response
	Include []*string
}

// GetSchedulesOutput represents the output of a GetSchedules operation.
type GetSchedulesOutput struct {
	_         struct{}
	Schedules []*Schedule
}

// GetSchedules gets list on-call schedules. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules/get
func (c *Client) GetSchedules(input *GetSchedulesInput) (*GetSchedulesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.schedules, q.Encode()))
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
		q.Add("until", *input.Until)
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

// GetScheduleUserOnCall gets the current user on call for specified schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}~1user-on-call/get
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

// UpdateScheduleInput represents the input of a UpdateSchedule operation.
type UpdateScheduleInput struct {
	_           struct{}
	ScheduleID  *int64
	Schedule    *Schedule
	AbortOnGaps *bool
}

// UpdateScheduleOutput represents the output of a UpdateSchedule operation.
type UpdateScheduleOutput struct {
	_        struct{}
	Schedule *Schedule
}

// UpdateSchedule updates the specific Schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}/put
func (c *Client) UpdateSchedule(input *UpdateScheduleInput) (*UpdateScheduleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("schedule id is required")
	}
	if input.Schedule == nil {
		return nil, errors.New("schedule input is required")
	}

	q := url.Values{}

	if input.AbortOnGaps != nil {
		q.Add("abort-on-gaps", strconv.FormatBool(*input.AbortOnGaps))
	}

	url := fmt.Sprintf("%s/%d?%s", apiRoutes.schedules, *input.ScheduleID, q.Encode())

	resp, err := c.httpClient.R().SetBody(input.Schedule).Put(url)
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

	return &UpdateScheduleOutput{Schedule: schedule}, nil
}

// AddScheduleShiftOverrideInput represents the input of a AddScheduleShiftOverride operation.
type AddScheduleShiftOverrideInput struct {
	_          struct{}
	ScheduleID *int64
	Shift      *Shift
}

// AddScheduleShiftOverrideOutput represents the output of a AddScheduleShiftOverride operation.
type AddScheduleShiftOverrideOutput struct {
	_        struct{}
	Schedule *Schedule
}

// AddScheduleShiftOverride adds an override to a shift on the schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}~1overrides/put
func (c *Client) AddScheduleShiftOverride(input *AddScheduleShiftOverrideInput) (*AddScheduleShiftOverrideOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("schedule id is required")
	}
	if input.Shift == nil {
		return nil, errors.New("shift input is required")
	}

	url := fmt.Sprintf("%s/%d/overrides", apiRoutes.schedules, *input.ScheduleID)

	resp, err := c.httpClient.R().SetBody(input.Shift).Post(url)
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

	return &AddScheduleShiftOverrideOutput{Schedule: schedule}, nil
}

// DeleteScheduleInput represents the input of a DeleteSchedule operation.
type DeleteScheduleInput struct {
	_          struct{}
	ScheduleID *int64
}

// DeleteScheduleOutput represents the output of a DeleteSchedule operation.
type DeleteScheduleOutput struct {
	_ struct{}
}

// DeleteSchedule deletes the specified Schedule. https://api.ilert.com/api-docs/#tag/Schedules/paths/~1schedules~1{id}/delete
func (c *Client) DeleteSchedule(input *DeleteScheduleInput) (*DeleteScheduleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ScheduleID == nil {
		return nil, errors.New("schedule id is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.schedules, *input.ScheduleID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteScheduleOutput{}, nil
}
