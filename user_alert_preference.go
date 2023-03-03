package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UserNotificationPreference definition https://api.ilert.com/api-docs/#tag/Notification-Preferences
type UserAlertPreference struct {
	ID       int64       `json:"id,omitempty"`
	Method   string      `json:"method"`
	Contact  interface{} `json:"contact"` // is either UserEmailContact or UserPhoneNumberContact
	DelayMin int64       `json:"delayMin"`
	Type     string      `json:"type"`
}

// UserPreferenceMethod defines user notification preference method
var UserPreferenceMethod = struct {
	Email    string
	Sms      string
	Push     string
	Voice    string
	WhatsApp string
	Telegram string
}{
	Email:    "EMAIL",
	Sms:      "SMS",
	Push:     "PUSH",
	Voice:    "VOICE",
	WhatsApp: "WHATSAPP",
	Telegram: "TELEGRAM",
}

// UserAlertPreferenceMethodAll defines user alert notification preference method list
var UserAlertPreferenceMethodAll = []string{
	UserPreferenceMethod.Email,
	UserPreferenceMethod.Sms,
	UserPreferenceMethod.Push,
	UserPreferenceMethod.Voice,
	UserPreferenceMethod.WhatsApp,
	UserPreferenceMethod.Telegram,
}

// UserAlertPreferenceType defines user alert notification preference type
var UserAlertPreferenceType = struct {
	HighPriority string
	LowPriority  string
}{
	HighPriority: "HIGH_PRIORITY",
	LowPriority:  "LOW_PRIORITY",
}

// UserAlertPreferenceTypeAll defines user alert notification preference type list
var UserAlertPreferenceTypeAll = []string{
	UserAlertPreferenceType.HighPriority,
	UserAlertPreferenceType.LowPriority,
}

// CreateUserAlertPreferenceInput represents the input of a CreateUserAlertPreference operation.
type CreateUserAlertPreferenceInput struct {
	_                   struct{}
	UserID              *int64
	UserAlertPreference *UserAlertPreference
}

// CreateUserAlertPreferenceOutput represents the output of a CreateUserAlertPreference operation.
type CreateUserAlertPreferenceOutput struct {
	_                   struct{}
	UserAlertPreference *UserAlertPreference
}

// CreateUserAlertPreference creates a new alert notification preference for a user. Requires ADMIN privileges or user id equals your current user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1alerts/post
func (c *Client) CreateUserAlertPreference(input *CreateUserAlertPreferenceInput) (*CreateUserAlertPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserAlertPreference == nil {
		return nil, errors.New("user alert notification preference input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/alerts", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(input.UserAlertPreference).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserAlertPreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &CreateUserAlertPreferenceOutput{UserAlertPreference: preference}, nil
}

// GetUserAlertPreferenceInput represents the input of a GetUserAlertPreference operation.
type GetUserAlertPreferenceInput struct {
	_                     struct{}
	UserID                *int64
	UserAlertPreferenceID *int64
}

// GetUserAlertPreferenceOutput represents the output of a GetUserAlertPreference operation.
type GetUserAlertPreferenceOutput struct {
	_                   struct{}
	UserAlertPreference *UserAlertPreference
}

// GetUserAlertPreference gets an alert notification preference of a user by id. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1alerts~1{id}/get
func (c *Client) GetUserAlertPreference(input *GetUserAlertPreferenceInput) (*GetUserAlertPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserAlertPreferenceID == nil {
		return nil, errors.New("user alert notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/alerts/%d", apiRoutes.users, *input.UserID, *input.UserAlertPreferenceID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserAlertPreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &GetUserAlertPreferenceOutput{UserAlertPreference: preference}, nil
}

// GetUserAlertPreferencesInput represents the input of a GetUserAlertPreferences operation.
type GetUserAlertPreferencesInput struct {
	_      struct{}
	UserID *int64
}

// GetUserAlertPreferencesOutput represents the output of a GetUserAlertPreferences operation.
type GetUserAlertPreferencesOutput struct {
	_                    struct{}
	UserAlertPreferences []*UserAlertPreference
}

// GetUserAlertPreferences lists existing alert notification preferences of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1alerts/get
func (c *Client) GetUserAlertPreferences(input *GetUserAlertPreferencesInput) (*GetUserAlertPreferencesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/alerts", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preferences := make([]*UserAlertPreference, 0)
	err = json.Unmarshal(resp.Body(), &preferences)
	if err != nil {
		return nil, err
	}

	return &GetUserAlertPreferencesOutput{UserAlertPreferences: preferences}, nil
}

// UpdateUserAlertPreferenceInput represents the input of a UpdateUserAlertPreference operation.
type UpdateUserAlertPreferenceInput struct {
	_                     struct{}
	UserID                *int64
	UserAlertPreferenceID *int64
	UserAlertPreference   *UserAlertPreference
}

// UpdateUserAlertPreferenceOutput represents the output of a UpdateUserAlertPreference operation.
type UpdateUserAlertPreferenceOutput struct {
	_                   struct{}
	UserAlertPreference *UserAlertPreference
}

// UpdateUserAlertPreference updates an existing alert notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1alerts~1{id}/put
func (c *Client) UpdateUserAlertPreference(input *UpdateUserAlertPreferenceInput) (*UpdateUserAlertPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserAlertPreferenceID == nil {
		return nil, errors.New("user alert notification preference id is required")
	}
	if input.UserAlertPreference == nil {
		return nil, errors.New("user input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/alerts/%d", apiRoutes.users, *input.UserID, *input.UserAlertPreferenceID)
	resp, err := c.httpClient.R().SetBody(input.UserAlertPreference).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserAlertPreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &UpdateUserAlertPreferenceOutput{UserAlertPreference: preference}, nil
}

// DeleteUserAlertPreferenceInput represents the input of a DeleteUserAlertPreference operation.
type DeleteUserAlertPreferenceInput struct {
	_                     struct{}
	UserID                *int64
	UserAlertPreferenceID *int64
}

// DeleteUserAlertPreferenceOutput represents the output of a DeleteUserAlertPreference operation.
type DeleteUserAlertPreferenceOutput struct {
	_ struct{}
}

// DeleteUserAlertPreference deletes the specified alert notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1alerts~1{id}/delete
func (c *Client) DeleteUserAlertPreference(input *DeleteUserAlertPreferenceInput) (*DeleteUserAlertPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserAlertPreferenceID == nil {
		return nil, errors.New("user alert notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/alerts/%d", apiRoutes.users, *input.UserID, *input.UserAlertPreferenceID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserAlertPreferenceOutput{}, nil
}
