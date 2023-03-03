package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UserNotificationPreference definition https://api.ilert.com/api-docs/#tag/Notification-Preferences
type UserUpdatePreference struct {
	ID      int64       `json:"id,omitempty"`
	Method  string      `json:"method"`
	Contact interface{} `json:"contact"` // is either UserEmailContact or UserPhoneNumberContact
	Type    string      `json:"type"`
}

// UserUpdatePreferenceMethodAll defines user update notification preference method list
var UserUpdatePreferenceMethodAll = []string{
	UserPreferenceMethod.Email,
	UserPreferenceMethod.Sms,
	UserPreferenceMethod.Push,
	UserPreferenceMethod.WhatsApp,
	UserPreferenceMethod.Telegram,
}

// UserUpdatePreferenceType defines user update notification preference type
var UserUpdatePreferenceType = struct {
	AlertAccepted  string
	AlertResolved  string
	AlertEscalated string
}{
	AlertAccepted:  "ALERT_ACCEPTED",
	AlertResolved:  "ALERT_RESOLVED",
	AlertEscalated: "ALERT_ESCALATED",
}

// UserUpdatePreferenceTypeAll defines user update notification preference type list
var UserUpdatePreferenceTypeAll = []string{
	UserUpdatePreferenceType.AlertAccepted,
	UserUpdatePreferenceType.AlertResolved,
	UserUpdatePreferenceType.AlertEscalated,
}

// CreateUserUpdatePreferenceInput represents the input of a CreateUserUpdatePreference operation.
type CreateUserUpdatePreferenceInput struct {
	_                    struct{}
	UserID               *int64
	UserUpdatePreference *UserUpdatePreference
}

// CreateUserUpdatePreferenceOutput represents the output of a CreateUserUpdatePreference operation.
type CreateUserUpdatePreferenceOutput struct {
	_                    struct{}
	UserUpdatePreference *UserUpdatePreference
}

// CreateUserUpdatePreference creates a new update notification preference for a user. Requires ADMIN privileges or user id equals your current user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1updates/post
func (c *Client) CreateUserUpdatePreference(input *CreateUserUpdatePreferenceInput) (*CreateUserUpdatePreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserUpdatePreference == nil {
		return nil, errors.New("user update notification preference input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/updates", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(input.UserUpdatePreference).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserUpdatePreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &CreateUserUpdatePreferenceOutput{UserUpdatePreference: preference}, nil
}

// GetUserUpdatePreferenceInput represents the input of a GetUserUpdatePreference operation.
type GetUserUpdatePreferenceInput struct {
	_                      struct{}
	UserID                 *int64
	UserUpdatePreferenceID *int64
}

// GetUserUpdatePreferenceOutput represents the output of a GetUserUpdatePreference operation.
type GetUserUpdatePreferenceOutput struct {
	_                    struct{}
	UserUpdatePreference *UserUpdatePreference
}

// GetUserUpdatePreference gets an update notification preference of a user by id. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1updates~1{id}/get
func (c *Client) GetUserUpdatePreference(input *GetUserUpdatePreferenceInput) (*GetUserUpdatePreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserUpdatePreferenceID == nil {
		return nil, errors.New("user update notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/updates/%d", apiRoutes.users, *input.UserID, *input.UserUpdatePreferenceID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserUpdatePreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &GetUserUpdatePreferenceOutput{UserUpdatePreference: preference}, nil
}

// GetUserUpdatePreferencesInput represents the input of a GetUserUpdatePreferences operation.
type GetUserUpdatePreferencesInput struct {
	_      struct{}
	UserID *int64
}

// GetUserUpdatePreferencesOutput represents the output of a GetUserUpdatePreferences operation.
type GetUserUpdatePreferencesOutput struct {
	_                     struct{}
	UserUpdatePreferences []*UserUpdatePreference
}

// GetUserUpdatePreferences lists existing update notification preferences of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1updates/get
func (c *Client) GetUserUpdatePreferences(input *GetUserUpdatePreferencesInput) (*GetUserUpdatePreferencesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/updates", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preferences := make([]*UserUpdatePreference, 0)
	err = json.Unmarshal(resp.Body(), &preferences)
	if err != nil {
		return nil, err
	}

	return &GetUserUpdatePreferencesOutput{UserUpdatePreferences: preferences}, nil
}

// UpdateUserUpdatePreferenceInput represents the input of a UpdateUserUpdatePreference operation.
type UpdateUserUpdatePreferenceInput struct {
	_                      struct{}
	UserID                 *int64
	UserUpdatePreferenceID *int64
	UserUpdatePreference   *UserUpdatePreference
}

// UpdateUserUpdatePreferenceOutput represents the output of a UpdateUserUpdatePreference operation.
type UpdateUserUpdatePreferenceOutput struct {
	_                    struct{}
	UserUpdatePreference *UserUpdatePreference
}

// UpdateUserUpdatePreference updates an existing update notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1updates~1{id}/put
func (c *Client) UpdateUserUpdatePreference(input *UpdateUserUpdatePreferenceInput) (*UpdateUserUpdatePreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserUpdatePreferenceID == nil {
		return nil, errors.New("user update notification preference id is required")
	}
	if input.UserUpdatePreference == nil {
		return nil, errors.New("user input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/updates/%d", apiRoutes.users, *input.UserID, *input.UserUpdatePreferenceID)
	resp, err := c.httpClient.R().SetBody(input.UserUpdatePreference).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserUpdatePreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &UpdateUserUpdatePreferenceOutput{UserUpdatePreference: preference}, nil
}

// DeleteUserUpdatePreferenceInput represents the input of a DeleteUserUpdatePreference operation.
type DeleteUserUpdatePreferenceInput struct {
	_                      struct{}
	UserID                 *int64
	UserUpdatePreferenceID *int64
}

// DeleteUserUpdatePreferenceOutput represents the output of a DeleteUserUpdatePreference operation.
type DeleteUserUpdatePreferenceOutput struct {
	_ struct{}
}

// DeleteUserUpdatePreference deletes the specified update notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1updates~1{id}/delete
func (c *Client) DeleteUserUpdatePreference(input *DeleteUserUpdatePreferenceInput) (*DeleteUserUpdatePreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserUpdatePreferenceID == nil {
		return nil, errors.New("user update notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/updates/%d", apiRoutes.users, *input.UserID, *input.UserUpdatePreferenceID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserUpdatePreferenceOutput{}, nil
}
