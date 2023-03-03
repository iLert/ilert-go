package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UserNotificationPreference definition https://api.ilert.com/api-docs/#tag/Notification-Preferences
type UserDutyPreference struct {
	ID        int64       `json:"id,omitempty"`
	Method    string      `json:"method"`
	Contact   interface{} `json:"contact"` // is either UserEmailContact or UserPhoneNumberContact
	BeforeMin int64       `json:"beforeMin"`
	Type      string      `json:"type"`
}

// UserDutyPreferenceMethodAll defines user duty notification preference method list
var UserDutyPreferenceMethodAll = []string{
	UserPreferenceMethod.Email,
	UserPreferenceMethod.Sms,
	UserPreferenceMethod.Push,
	UserPreferenceMethod.WhatsApp,
	UserPreferenceMethod.Telegram,
}

// UserDutyPreferenceType defines user duty notification preference type
var UserDutyPreferenceType = struct {
	OnCall string
}{
	OnCall: "ON_CALL",
}

// UserDutyPreferenceTypeAll defines user duty notification preference type list
var UserDutyPreferenceTypeAll = []string{
	UserDutyPreferenceType.OnCall,
}

// CreateUserDutyPreferenceInput represents the input of a CreateUserDutyPreference operation.
type CreateUserDutyPreferenceInput struct {
	_                  struct{}
	UserID             *int64
	UserDutyPreference *UserDutyPreference
}

// CreateUserDutyPreferenceOutput represents the output of a CreateUserDutyPreference operation.
type CreateUserDutyPreferenceOutput struct {
	_                  struct{}
	UserDutyPreference *UserDutyPreference
}

// CreateUserDutyPreference creates a new duty notification preference for a user. Requires ADMIN privileges or user id equals your current user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1duties/post
func (c *Client) CreateUserDutyPreference(input *CreateUserDutyPreferenceInput) (*CreateUserDutyPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserDutyPreference == nil {
		return nil, errors.New("user duty notification preference input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/duties", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(input.UserDutyPreference).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserDutyPreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &CreateUserDutyPreferenceOutput{UserDutyPreference: preference}, nil
}

// GetUserDutyPreferenceInput represents the input of a GetUserDutyPreference operation.
type GetUserDutyPreferenceInput struct {
	_                    struct{}
	UserID               *int64
	UserDutyPreferenceID *int64
}

// GetUserDutyPreferenceOutput represents the output of a GetUserDutyPreference operation.
type GetUserDutyPreferenceOutput struct {
	_                  struct{}
	UserDutyPreference *UserDutyPreference
}

// GetUserDutyPreference gets an duty notification preference of a user by id. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1duties~1{id}/get
func (c *Client) GetUserDutyPreference(input *GetUserDutyPreferenceInput) (*GetUserDutyPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserDutyPreferenceID == nil {
		return nil, errors.New("user duty notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/duties/%d", apiRoutes.users, *input.UserID, *input.UserDutyPreferenceID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserDutyPreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &GetUserDutyPreferenceOutput{UserDutyPreference: preference}, nil
}

// GetUserDutyPreferencesInput represents the input of a GetUserDutyPreferences operation.
type GetUserDutyPreferencesInput struct {
	_      struct{}
	UserID *int64
}

// GetUserDutyPreferencesOutput represents the output of a GetUserDutyPreferences operation.
type GetUserDutyPreferencesOutput struct {
	_                   struct{}
	UserDutyPreferences []*UserDutyPreference
}

// GetUserDutyPreferences lists existing duty notification preferences of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1duties/get
func (c *Client) GetUserDutyPreferences(input *GetUserDutyPreferencesInput) (*GetUserDutyPreferencesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/duties", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preferences := make([]*UserDutyPreference, 0)
	err = json.Unmarshal(resp.Body(), &preferences)
	if err != nil {
		return nil, err
	}

	return &GetUserDutyPreferencesOutput{UserDutyPreferences: preferences}, nil
}

// UpdateUserDutyPreferenceInput represents the input of a UpdateUserDutyPreference operation.
type UpdateUserDutyPreferenceInput struct {
	_                    struct{}
	UserID               *int64
	UserDutyPreferenceID *int64
	UserDutyPreference   *UserDutyPreference
}

// UpdateUserDutyPreferenceOutput represents the output of a UpdateUserDutyPreference operation.
type UpdateUserDutyPreferenceOutput struct {
	_                  struct{}
	UserDutyPreference *UserDutyPreference
}

// UpdateUserDutyPreference updates an existing duty notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1duties~1{id}/put
func (c *Client) UpdateUserDutyPreference(input *UpdateUserDutyPreferenceInput) (*UpdateUserDutyPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserDutyPreferenceID == nil {
		return nil, errors.New("user duty notification preference id is required")
	}
	if input.UserDutyPreference == nil {
		return nil, errors.New("user input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/duties/%d", apiRoutes.users, *input.UserID, *input.UserDutyPreferenceID)
	resp, err := c.httpClient.R().SetBody(input.UserDutyPreference).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	preference := &UserDutyPreference{}
	err = json.Unmarshal(resp.Body(), preference)
	if err != nil {
		return nil, err
	}

	return &UpdateUserDutyPreferenceOutput{UserDutyPreference: preference}, nil
}

// DeleteUserDutyPreferenceInput represents the input of a DeleteUserDutyPreference operation.
type DeleteUserDutyPreferenceInput struct {
	_                    struct{}
	UserID               *int64
	UserDutyPreferenceID *int64
}

// DeleteUserDutyPreferenceOutput represents the output of a DeleteUserDutyPreference operation.
type DeleteUserDutyPreferenceOutput struct {
	_ struct{}
}

// DeleteUserDutyPreference deletes the specified duty notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1duties~1{id}/delete
func (c *Client) DeleteUserDutyPreference(input *DeleteUserDutyPreferenceInput) (*DeleteUserDutyPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserDutyPreferenceID == nil {
		return nil, errors.New("user duty notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/duties/%d", apiRoutes.users, *input.UserID, *input.UserDutyPreferenceID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserDutyPreferenceOutput{}, nil
}
