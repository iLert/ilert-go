package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UserEmailContact definition https://api.ilert.com/api-docs/#tag/Contacts
type UserEmailContact struct {
	ID     int64  `json:"id,omitempty"`
	Target string `json:"target"`
	Status string `json:"status,omitempty"`
}

// UserContactStatus defines user contact status
var UserContactStatus = struct {
	Ok          string
	Locked      string
	Blacklisted string
}{
	Ok:          "OK",
	Locked:      "LOCKED",
	Blacklisted: "BLACKLISTED",
}

// UserContactStatusAll defines user contact status list
var UserContactStatusAll = []string{
	UserContactStatus.Ok,
	UserContactStatus.Locked,
	UserContactStatus.Blacklisted,
}

// CreateUserEmailContactInput represents the input of a CreateUserEmailContact operation.
type CreateUserEmailContactInput struct {
	_                struct{}
	UserID           *int64
	UserEmailContact *UserEmailContact
}

// CreateUserEmailContactOutput represents the output of a CreateUserEmailContact operation.
type CreateUserEmailContactOutput struct {
	_                struct{}
	UserEmailContact *UserEmailContact
}

// CreateUserEmailContact creates a new email contact for a user. Requires ADMIN privileges or user id equals your current user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1emails/post
func (c *Client) CreateUserEmailContact(input *CreateUserEmailContactInput) (*CreateUserEmailContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserEmailContact == nil {
		return nil, errors.New("user email contact input is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/emails", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(input.UserEmailContact).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserEmailContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &CreateUserEmailContactOutput{UserEmailContact: contact}, nil
}

// GetUserEmailContactInput represents the input of a GetUserEmailContact operation.
type GetUserEmailContactInput struct {
	_                  struct{}
	UserID             *int64
	UserEmailContactID *int64
}

// GetUserEmailContactOutput represents the output of a GetUserEmailContact operation.
type GetUserEmailContactOutput struct {
	_                struct{}
	UserEmailContact *UserEmailContact
}

// GetUserEmailContact gets an email contact of a user by id. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1emails~1{id}/get
func (c *Client) GetUserEmailContact(input *GetUserEmailContactInput) (*GetUserEmailContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserEmailContactID == nil {
		return nil, errors.New("user email contact id is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/emails/%d", apiRoutes.users, *input.UserID, *input.UserEmailContactID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserEmailContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &GetUserEmailContactOutput{UserEmailContact: contact}, nil
}

// GetUserEmailContactsInput represents the input of a GetUserEmailContacts operation.
type GetUserEmailContactsInput struct {
	_      struct{}
	UserID *int64
}

// GetUserEmailContactsOutput represents the output of a GetUserEmailContacts operation.
type GetUserEmailContactsOutput struct {
	_                 struct{}
	UserEmailContacts []*UserEmailContact
}

// GetUserEmailContacts lists existing email contacts of a user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1emails/get
func (c *Client) GetUserEmailContacts(input *GetUserEmailContactsInput) (*GetUserEmailContactsOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/emails", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contacts := make([]*UserEmailContact, 0)
	err = json.Unmarshal(resp.Body(), &contacts)
	if err != nil {
		return nil, err
	}

	return &GetUserEmailContactsOutput{UserEmailContacts: contacts}, nil
}

// SearchUserEmailContactInput represents the input of a SearchUserEmailContact operation.
type SearchUserEmailContactInput struct {
	_                      struct{}
	UserID                 *int64
	UserEmailContactTarget *string
}

// SearchUserEmailContactOutput represents the output of a SearchUserEmailContact operation.
type SearchUserEmailContactOutput struct {
	_                struct{}
	UserEmailContact *UserEmailContact
}

// SearchUserEmailContact gets the email contact with specified target of a user.
func (c *Client) SearchUserEmailContact(input *SearchUserEmailContactInput) (*SearchUserEmailContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserEmailContactTarget == nil {
		return nil, errors.New("user email contact target is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/emails/search-target", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(UserEmailContact{Target: *input.UserEmailContactTarget}).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserEmailContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &SearchUserEmailContactOutput{UserEmailContact: contact}, nil
}

// UpdateUserEmailContactInput represents the input of a UpdateUserEmailContact operation.
type UpdateUserEmailContactInput struct {
	_                  struct{}
	UserID             *int64
	UserEmailContactID *int64
	UserEmailContact   *UserEmailContact
}

// UpdateUserEmailContactOutput represents the output of a UpdateUserEmailContact operation.
type UpdateUserEmailContactOutput struct {
	_                struct{}
	UserEmailContact *UserEmailContact
}

// UpdateUserEmailContact updates an existing email contact of a user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1emails~1{id}/put
func (c *Client) UpdateUserEmailContact(input *UpdateUserEmailContactInput) (*UpdateUserEmailContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserEmailContactID == nil {
		return nil, errors.New("user email contact id is required")
	}
	if input.UserEmailContact == nil {
		return nil, errors.New("user input is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/emails/%d", apiRoutes.users, *input.UserID, *input.UserEmailContactID)
	resp, err := c.httpClient.R().SetBody(input.UserEmailContact).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserEmailContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &UpdateUserEmailContactOutput{UserEmailContact: contact}, nil
}

// DeleteUserEmailContactInput represents the input of a DeleteUserEmailContact operation.
type DeleteUserEmailContactInput struct {
	_                  struct{}
	UserID             *int64
	UserEmailContactID *int64
}

// DeleteUserEmailContactOutput represents the output of a DeleteUserEmailContact operation.
type DeleteUserEmailContactOutput struct {
	_ struct{}
}

// DeleteUserEmailContact deletes the specified email contact of a user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1emails~1{id}/delete
func (c *Client) DeleteUserEmailContact(input *DeleteUserEmailContactInput) (*DeleteUserEmailContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserEmailContactID == nil {
		return nil, errors.New("user email contact id is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/emails/%d", apiRoutes.users, *input.UserID, *input.UserEmailContactID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserEmailContactOutput{}, nil
}
