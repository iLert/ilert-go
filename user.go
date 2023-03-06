package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// User definition https://api.ilert.com/api-docs/#tag/Users
type User struct {
	ID         int64  `json:"id,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
	Position   string `json:"position,omitempty"`
	Department string `json:"department,omitempty"`
	Language   string `json:"language,omitempty"`
	Role       string `json:"role,omitempty"`
	ShiftColor string `json:"shiftColor,omitempty"`
}

// Phone definition
type Phone struct {
	RegionCode string `json:"regionCode"`
	Number     string `json:"number"`
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
		return nil, errors.New("user input is required")
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

// GetUser gets the user with specified id or username. https://api.ilert.com/api-docs/#tag/Users/paths/~1users~1{user-id}/get
func (c *Client) GetUser(input *GetUserInput) (*GetUserOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("user id or username is required")
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

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetUsersOutput represents the output of a GetUsers operation.
type GetUsersOutput struct {
	_     struct{}
	Users []*User
}

// GetUsers lists existing users. https://api.ilert.com/api-docs/#tag/Users/paths/~1users/get
func (c *Client) GetUsers(input *GetUsersInput) (*GetUsersOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.users, q.Encode()))
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

// SearchUserInput represents the input of a SearchUser operation.
type SearchUserInput struct {
	_         struct{}
	UserEmail *string
}

// SearchUserOutput represents the output of a SearchUser operation.
type SearchUserOutput struct {
	_    struct{}
	User *User
}

// SearchUser gets the user with specified name.
func (c *Client) SearchUser(input *SearchUserInput) (*SearchUserOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserEmail == nil {
		return nil, errors.New("user email is required")
	}

	resp, err := c.httpClient.R().SetBody(User{Email: *input.UserEmail}).Post(fmt.Sprintf("%s/search-email", apiRoutes.users))
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

	return &SearchUserOutput{User: user}, nil
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
		return nil, errors.New("user input is required")
	}
	if input.UserID == nil && input.Username == nil {
		return nil, errors.New("user id or username is required")
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
		return nil, errors.New("user id or username is required")
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
