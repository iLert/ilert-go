package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// User definition https://api.ilert.com/api-docs/#!/Users
type User struct {
	ID                                        int64                          `json:"id,omitempty"`
	Username                                  string                         `json:"username,omitempty"`
	FirstName                                 string                         `json:"firstName,omitempty"`
	LastName                                  string                         `json:"lastName,omitempty"`
	Email                                     string                         `json:"email,omitempty"`
	Mobile                                    *Phone                         `json:"mobile,omitempty"`
	Landline                                  *Phone                         `json:"landline,omitempty"`
	Position                                  string                         `json:"position,omitempty"`
	Department                                string                         `json:"department,omitempty"`
	Timezone                                  string                         `json:"timezone,omitempty"`
	Language                                  string                         `json:"language,omitempty"`
	Role                                      string                         `json:"role,omitempty"`
	NotificationPreferences                   []NotificationPreference       `json:"notificationPreferences,omitempty"`
	LowNotificationPreferences                []NotificationPreference       `json:"lowPriorityNotificationPreferences,omitempty"`
	OnCallNotificationPreferences             []OnCallNotificationPreference `json:"onCallNotificationPreferences,omitempty"`
	SubscribedAlertUpdateStates               []string                       `json:"subscribedAlertUpdateStates,omitempty"`
	SubscribedIncidentUpdateStates            []string                       `json:"subscribedIncidentUpdateStates,omitempty"` // @deprecated
	SubscribedAlertUpdateNotificationTypes    []string                       `json:"subscribedAlertUpdateNotificationTypes,omitempty"`
	SubscribedIncidentUpdateNotificationTypes []string                       `json:"subscribedIncidentUpdateNotificationTypes,omitempty"` // @deprecated
}

// Phone definition
type Phone struct {
	RegionCode string `json:"regionCode"`
	Number     string `json:"number"`
}

// NotificationPreference definition
type NotificationPreference struct {
	Delay  int    `json:"delay"`
	Method string `json:"method"` // e.g. EMAIL
}

// OnCallNotificationPreference definition
type OnCallNotificationPreference struct {
	BeforeMin int    `json:"beforeMin"`
	Method    string `json:"method"` // e.g. EMAIL
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

// UserRoleAll defines user roles list
var UserRoleAll = []string{
	UserRole.User,
	UserRole.Admin,
	UserRole.Stakeholder,
}

// UserAlertUpdateStates defines user alert update states
var UserAlertUpdateStates = struct {
	Accepted  string
	Escalated string
	Resolved  string
}{
	Accepted:  "ACCEPTED",
	Escalated: "ESCALATED",
	Resolved:  "RESOLVED",
}

// UserAlertUpdateStatesAll defines user alert update states list
var UserAlertUpdateStatesAll = []string{
	UserAlertUpdateStates.Accepted,
	UserAlertUpdateStates.Escalated,
	UserAlertUpdateStates.Resolved,
}

// UserAlertUpdateNotificationTypes defines user alert update notification types
var UserAlertUpdateNotificationTypes = struct {
	Email         string
	PushAndroid   string
	PushIPhone    string
	SMS           string
	VoiceMobile   string
	VoiceLandline string
	WhatsApp      string
}{
	Email:         "EMAIL",
	PushAndroid:   "ANDROID",
	PushIPhone:    "IPHONE",
	SMS:           "SMS",
	VoiceMobile:   "VOICE_MOBILE",
	VoiceLandline: "VOICE_LANDLINE",
	WhatsApp:      "WHATSAPP",
}

// UserAlertUpdateNotificationTypesAll defines user alert update notification types list
var UserAlertUpdateNotificationTypesAll = []string{
	UserAlertUpdateNotificationTypes.Email,
	UserAlertUpdateNotificationTypes.PushAndroid,
	UserAlertUpdateNotificationTypes.PushIPhone,
	UserAlertUpdateNotificationTypes.SMS,
	UserAlertUpdateNotificationTypes.VoiceMobile,
	UserAlertUpdateNotificationTypes.VoiceLandline,
}

// UserLanguage defines user language
var UserLanguage = struct {
	English string
	German  string
}{
	English: "en",
	German:  "de",
}

// UserLanguageAll defines user language list
var UserLanguageAll = []string{
	UserLanguage.English,
	UserLanguage.German,
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

func (u *User) GetUsername() string {
	return u.Username
}

// CreateUser creates a new user. Requires ADMIN privileges. https://api.ilert.com/api-docs/#tag/Users/paths/~1users/post
func (c *Client) CreateUser(input *CreateUserInput) (*CreateUserOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.User == nil {
		return nil, errors.New("User input is required")
	}

	if len(input.User.SubscribedAlertUpdateNotificationTypes) > 0 && len(input.User.SubscribedIncidentUpdateNotificationTypes) > 0 {
		input.User.SubscribedIncidentUpdateNotificationTypes = nil
	}
	if len(input.User.SubscribedAlertUpdateNotificationTypes) == 0 {
		input.User.SubscribedAlertUpdateNotificationTypes = input.User.SubscribedIncidentUpdateNotificationTypes
		input.User.SubscribedIncidentUpdateNotificationTypes = nil
	}

	if len(input.User.SubscribedAlertUpdateStates) > 0 && len(input.User.SubscribedIncidentUpdateStates) > 0 {
		input.User.SubscribedIncidentUpdateStates = nil
	}
	if len(input.User.SubscribedAlertUpdateStates) == 0 {
		input.User.SubscribedAlertUpdateStates = input.User.SubscribedIncidentUpdateStates
		input.User.SubscribedIncidentUpdateStates = nil
	}

	resp, err := c.httpClient.R().SetBody(input.User).Post(apiRoutes.users)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	user := &User{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{User: user}, nil
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
		return nil, errors.New("input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("User id or username is required")
	}
	var url string
	if input.UserID != nil {
		url = fmt.Sprintf("%s/%d", apiRoutes.users, *input.UserID)
	} else {
		url = fmt.Sprintf("%s/%s", apiRoutes.users, *input.Username)
	}
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	user := &User{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{User: user}, nil
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
	resp, err := c.httpClient.R().Get(apiRoutes.users)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	users := make([]*User, 0)
	err = json.Unmarshal(resp.Body(), &users)
	if err != nil {
		return nil, err
	}

	return &GetUsersOutput{Users: users}, nil
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
		return nil, errors.New("input is required")
	}
	input.Username = String("current")
	return c.UpdateUser(input)
}

// UpdateUser updates an existing user. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1{user-id}/put
func (c *Client) UpdateUser(input *UpdateUserInput) (*UpdateUserOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.User == nil {
		return nil, errors.New("User input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("User id or username is required")
	}

	if len(input.User.SubscribedAlertUpdateNotificationTypes) > 0 && len(input.User.SubscribedIncidentUpdateNotificationTypes) > 0 {
		input.User.SubscribedIncidentUpdateNotificationTypes = nil
	}
	if len(input.User.SubscribedAlertUpdateNotificationTypes) == 0 {
		input.User.SubscribedAlertUpdateNotificationTypes = input.User.SubscribedIncidentUpdateNotificationTypes
		input.User.SubscribedIncidentUpdateNotificationTypes = nil
	}

	if len(input.User.SubscribedAlertUpdateStates) > 0 && len(input.User.SubscribedIncidentUpdateStates) > 0 {
		input.User.SubscribedIncidentUpdateStates = nil
	}
	if len(input.User.SubscribedAlertUpdateStates) == 0 {
		input.User.SubscribedAlertUpdateStates = input.User.SubscribedIncidentUpdateStates
		input.User.SubscribedIncidentUpdateStates = nil
	}

	var url string
	if input.UserID != nil {
		url = fmt.Sprintf("%s/%d", apiRoutes.users, *input.UserID)
	} else {
		url = fmt.Sprintf("%s/%s", apiRoutes.users, *input.Username)
	}
	resp, err := c.httpClient.R().SetBody(input.User).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	user := &User{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{User: user}, nil
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
		return nil, errors.New("input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("User id or username is required")
	}
	var url string
	if input.UserID != nil {
		url = fmt.Sprintf("%s/%d", apiRoutes.users, *input.UserID)
	} else {
		url = fmt.Sprintf("%s/%s", apiRoutes.users, *input.Username)
	}
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserOutput{}, nil
}
