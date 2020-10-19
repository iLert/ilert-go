package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// User definition https://api.ilert.com/api-docs/#!/Users
type User struct {
	ID                                        int64                     `json:"id"`
	Username                                  string                    `json:"username"`
	FirstName                                 string                    `json:"firstName"`
	LastName                                  string                    `json:"lastName"`
	Email                                     string                    `json:"email"`
	Mobile                                    *Phone                    `json:"mobile"`
	Landline                                  *Phone                    `json:"landline"`
	Position                                  string                    `json:"position"`
	Department                                string                    `json:"department"`
	Timezone                                  string                    `json:"timezone"`
	Language                                  string                    `json:"language"`
	Role                                      string                    `json:"role"`
	NotificationPreferences                   []NotificationPreferences `json:"notificationPreferences"`
	LowNotificationPreferences                []NotificationPreferences `json:"lowPriorityNotificationPreferences"`
	OnCallNotificationPreferences             []NotificationPreferences `json:"onCallNotificationPreferences"`
	SubscribedIncidentUpdateStates            []string                  `json:"subscribedIncidentUpdateStates"`
	SubscribedIncidentUpdateNotificationTypes []string                  `json:"subscribedIncidentUpdateNotificationTypes"`
}

// Phone definition
type Phone struct {
	RegionCode string `json:"regionCode"`
	Number     string `json:"number"`
}

// NotificationPreferences definition
type NotificationPreferences struct {
	Delay  int    `json:"delay"`
	Method string `json:"method"` // e.g. EMAIL
}

// UserRole defines user roles
var UserRole = struct {
	User        string
	Admin       string
	Stakeholder string
}{
	User:        "USER",
	Admin:       "ADMIN",
	Stakeholder: "STAKEHOLDER",
}

// UserIncidentUpdateStates defines user incident update states
var UserIncidentUpdateStates = struct {
	Accepted  string
	Escalated string
	Resolved  string
}{
	Accepted:  "ACCEPTED",
	Escalated: "ESCALATED",
	Resolved:  "RESOLVED",
}

// UserIncidentUpdateNotificationTypes defines user incident update notification types
var UserIncidentUpdateNotificationTypes = struct {
	Email         string
	PushAndroid   string
	PushIPhone    string
	SMS           string
	VoiceMobile   string
	VoiceLandline string
}{
	Email:         "EMAIL",
	PushAndroid:   "ANDROID",
	PushIPhone:    "IPHONE",
	SMS:           "SMS",
	VoiceMobile:   "VOICE_MOBILE",
	VoiceLandline: "VOICE_LANDLINE",
}

// UserLanguage defines user language
var UserLanguage = struct {
	English string
	German  string
}{
	English: "en",
	German:  "de",
}

// CreateUserInput represents the input of a CreateUser operation.
type CreateUserInput struct {
	_    struct{}
	User *User
}

// CreateUserOutput represents the output of a CreateUser operation.
type CreateUserOutput struct {
	_    struct{}
	User *User
}

// CreateUser creates a new user. Requires ADMIN privileges. https://api.ilert.com/api-docs/#tag/Users/paths/~1users/post
func (c *Client) CreateUser(input *CreateUserInput) (*CreateUserOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.User == nil {
		return nil, errors.New("User input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.User).Post("/api/v1/users")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 201); err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	output := &CreateUserOutput{User: user}

	return output, nil
}

// GetUserInput represents the input of a GetUser operation.
type GetUserInput struct {
	_        struct{}
	UserID   *int64
	Username *string
}

// GetUserOutput represents the output of a GetUser operation.
type GetUserOutput struct {
	_    struct{}
	User *User
}

// GetCurrentUser gets the currently authenticated user. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1current/get
func (c *Client) GetCurrentUser() (*GetUserOutput, error) {
	input := &GetUserInput{Username: String("current")}
	return c.GetUser(input)
}

// GetUser gets information about a user including contact methods and notification preferences. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1{user-id}/get
func (c *Client) GetUser(input *GetUserInput) (*GetUserOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("User id or username is required")
	}
	var url string
	if input.UserID != nil {
		url = fmt.Sprintf("/api/v1/users/%d", *input.UserID)
	} else {
		url = fmt.Sprintf("/api/v1/users/%s", *input.Username)
	}
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	output := &GetUserOutput{
		User: user,
	}

	return output, nil
}

// GetUsersInput represents the input of a GetUsers operation.
type GetUsersInput struct {
	_ struct{}
}

// GetUsersOutput represents the output of a GetUsers operation.
type GetUsersOutput struct {
	_     struct{}
	Users []*User
}

// GetUsers lists existing users. https://api.ilert.com/api-docs/#tag/Users/paths/~1users/get
func (c *Client) GetUsers(input *GetUsersInput) (*GetUsersOutput, error) {
	resp, err := c.httpClient.R().Get("/api/v1/users")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	err = json.Unmarshal(resp.Body(), &users)
	if err != nil {
		return nil, err
	}

	output := &GetUsersOutput{Users: users}

	return output, nil
}

// UpdateUserInput represents the input of a UpdateUser operation.
type UpdateUserInput struct {
	_        struct{}
	UserID   *int64
	Username *string
	User     *User
}

// UpdateUserOutput represents the output of a UpdateUser operation.
type UpdateUserOutput struct {
	_    struct{}
	User *User
}

// UpdateCurrentUser updates the currently authenticated user. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1current/put
func (c *Client) UpdateCurrentUser(input *UpdateUserInput) (*UpdateUserOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	input.Username = String("current")
	return c.UpdateUser(input)
}

// UpdateUser updates an existing user. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1{user-id}/put
func (c *Client) UpdateUser(input *UpdateUserInput) (*UpdateUserOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.User == nil {
		return nil, errors.New("User input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("User id or username is required")
	}
	var url string
	if input.UserID != nil {
		url = fmt.Sprintf("/api/v1/users/%d", *input.UserID)
	} else {
		url = fmt.Sprintf("/api/v1/users/%s", *input.Username)
	}
	resp, err := c.httpClient.R().SetBody(input.User).Put(url)
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	output := &UpdateUserOutput{User: user}

	return output, nil
}

// DeleteUserInput represents the input of a DeleteUser operation.
type DeleteUserInput struct {
	_        struct{}
	UserID   *int64
	Username *string
}

// DeleteUserOutput represents the output of a DeleteUser operation.
type DeleteUserOutput struct {
	_ struct{}
}

// DeleteUser deletes the specified user. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1{user-id}/delete
func (c *Client) DeleteUser(input *DeleteUserInput) (*DeleteUserOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("User id or username is required")
	}
	var url string
	if input.UserID != nil {
		url = fmt.Sprintf("/api/v1/users/%d", *input.UserID)
	} else {
		url = fmt.Sprintf("/api/v1/users/%s", *input.Username)
	}
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 204); err != nil {
		return nil, err
	}

	output := &DeleteUserOutput{}
	return output, nil
}
