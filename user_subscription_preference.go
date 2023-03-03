package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UserNotificationPreference definition https://api.ilert.com/api-docs/#tag/Notification-Preferences
type UserSubscriptionPreference struct {
	ID        int64       `json:"id,omitempty"`
	Method    string      `json:"method"`
	Contact   interface{} `json:"contact"` // is either UserEmailContact or UserPhoneNumberContact
	BeforeMin int64       `json:"beforeMin"`
	Type      string      `json:"type"`
}

// UserSubscriptionPreferenceMethodAll defines user subscription notification preference method list
var UserSubscriptionPreferenceMethodAll = []string{
	UserPreferenceMethod.Email,
	UserPreferenceMethod.Sms,
	UserPreferenceMethod.Push,
	UserPreferenceMethod.WhatsApp,
	UserPreferenceMethod.Telegram,
}

// UserSubscriptionPreferenceType defines user subscription notification preference type
var UserSubscriptionPreferenceType = struct {
	OnCall string
}{
	OnCall: "ON_CALL",
}

// UserSubscriptionPreferenceTypeAll defines user subscription notification preference type list
var UserSubscriptionPreferenceTypeAll = []string{
	UserSubscriptionPreferenceType.OnCall,
}

// CreateUserSubscriptionPreferenceInput represents the input of a CreateUserSubscriptionPreference operation.
type CreateUserSubscriptionPreferenceInput struct {
	_                          struct{}
	UserID                     *int64
	UserSubscriptionPreference *UserSubscriptionPreference
}

// CreateUserSubscriptionPreferenceOutput represents the output of a CreateUserSubscriptionPreference operation.
type CreateUserSubscriptionPreferenceOutput struct {
	_                          struct{}
	UserSubscriptionPreference *UserSubscriptionPreference
}

// CreateUserSubscriptionPreference creates a new subscription notification preference for a user. Requires ADMIN privileges or user id equals your current user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1subscriptions/post
func (c *Client) CreateUserSubscriptionPreference(input *CreateUserSubscriptionPreferenceInput) (*CreateUserSubscriptionPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserSubscriptionPreference == nil {
		return nil, errors.New("user subscription notification preference input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/subscriptions", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(input.UserSubscriptionPreference).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserSubscriptionPreference{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &CreateUserSubscriptionPreferenceOutput{UserSubscriptionPreference: contact}, nil
}

// GetUserSubscriptionPreferenceInput represents the input of a GetUserSubscriptionPreference operation.
type GetUserSubscriptionPreferenceInput struct {
	_                            struct{}
	UserID                       *int64
	UserSubscriptionPreferenceID *int64
}

// GetUserSubscriptionPreferenceOutput represents the output of a GetUserSubscriptionPreference operation.
type GetUserSubscriptionPreferenceOutput struct {
	_                          struct{}
	UserSubscriptionPreference *UserSubscriptionPreference
}

// GetUserSubscriptionPreference gets an subscription notification preference of a user by id. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1subscriptions~1{id}/get
func (c *Client) GetUserSubscriptionPreference(input *GetUserSubscriptionPreferenceInput) (*GetUserSubscriptionPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserSubscriptionPreferenceID == nil {
		return nil, errors.New("user subscription notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/subscriptions/%d", apiRoutes.users, *input.UserID, *input.UserSubscriptionPreferenceID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserSubscriptionPreference{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &GetUserSubscriptionPreferenceOutput{UserSubscriptionPreference: contact}, nil
}

// GetUserSubscriptionPreferencesInput represents the input of a GetUserSubscriptionPreferences operation.
type GetUserSubscriptionPreferencesInput struct {
	_      struct{}
	UserID *int64
}

// GetUserSubscriptionPreferencesOutput represents the output of a GetUserSubscriptionPreferences operation.
type GetUserSubscriptionPreferencesOutput struct {
	_                           struct{}
	UserSubscriptionPreferences []*UserSubscriptionPreference
}

// GetUserSubscriptionPreferences lists existing subscription notification preferences of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1subscriptions/get
func (c *Client) GetUserSubscriptionPreferences(input *GetUserSubscriptionPreferencesInput) (*GetUserSubscriptionPreferencesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/subscriptions", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contacts := make([]*UserSubscriptionPreference, 0)
	err = json.Unmarshal(resp.Body(), &contacts)
	if err != nil {
		return nil, err
	}

	return &GetUserSubscriptionPreferencesOutput{UserSubscriptionPreferences: contacts}, nil
}

// UpdateUserSubscriptionPreferenceInput represents the input of a UpdateUserSubscriptionPreference operation.
type UpdateUserSubscriptionPreferenceInput struct {
	_                            struct{}
	UserID                       *int64
	UserSubscriptionPreferenceID *int64
	UserSubscriptionPreference   *UserSubscriptionPreference
}

// UpdateUserSubscriptionPreferenceOutput represents the output of a UpdateUserSubscriptionPreference operation.
type UpdateUserSubscriptionPreferenceOutput struct {
	_                          struct{}
	UserSubscriptionPreference *UserSubscriptionPreference
}

// UpdateUserSubscriptionPreference updates an existing subscription notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1subscriptions~1{id}/put
func (c *Client) UpdateUserSubscriptionPreference(input *UpdateUserSubscriptionPreferenceInput) (*UpdateUserSubscriptionPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserSubscriptionPreferenceID == nil {
		return nil, errors.New("user subscription notification preference id is required")
	}
	if input.UserSubscriptionPreference == nil {
		return nil, errors.New("user input is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/subscriptions/%d", apiRoutes.users, *input.UserID, *input.UserSubscriptionPreferenceID)
	resp, err := c.httpClient.R().SetBody(input.UserSubscriptionPreference).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	user := &UserSubscriptionPreference{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserSubscriptionPreferenceOutput{UserSubscriptionPreference: user}, nil
}

// DeleteUserSubscriptionPreferenceInput represents the input of a DeleteUserSubscriptionPreference operation.
type DeleteUserSubscriptionPreferenceInput struct {
	_                            struct{}
	UserID                       *int64
	UserSubscriptionPreferenceID *int64
}

// DeleteUserSubscriptionPreferenceOutput represents the output of a DeleteUserSubscriptionPreference operation.
type DeleteUserSubscriptionPreferenceOutput struct {
	_ struct{}
}

// DeleteUserSubscriptionPreference deletes the specified subscription notification preference of a user. https://api.ilert.com/api-docs/#tag/Notification-Preferences/paths/~1users~1{user-id}~1notification-preferences~1subscriptions~1{id}/delete
func (c *Client) DeleteUserSubscriptionPreference(input *DeleteUserSubscriptionPreferenceInput) (*DeleteUserSubscriptionPreferenceOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserSubscriptionPreferenceID == nil {
		return nil, errors.New("user subscription notification preference id is required")
	}

	url := fmt.Sprintf("%s/%d/notification-preferences/subscriptions/%d", apiRoutes.users, *input.UserID, *input.UserSubscriptionPreferenceID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserSubscriptionPreferenceOutput{}, nil
}
