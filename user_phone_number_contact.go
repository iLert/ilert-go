package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UserPhoneNumberContact definition https://api.ilert.com/api-docs/#tag/Contacts
type UserPhoneNumberContact struct {
	ID         int64  `json:"id,omitempty"`
	RegionCode string `json:"regionCode"`
	Target     string `json:"target"`
	Status     string `json:"status,omitempty"`
}

// CreateUserPhoneNumberContactInput represents the input of a CreateUserPhoneNumberContact operation.
type CreateUserPhoneNumberContactInput struct {
	_                      struct{}
	UserID                 *int64
	UserPhoneNumberContact *UserPhoneNumberContact
}

// CreateUserPhoneNumberContactOutput represents the output of a CreateUserPhoneNumberContact operation.
type CreateUserPhoneNumberContactOutput struct {
	_                      struct{}
	UserPhoneNumberContact *UserPhoneNumberContact
}

// CreateUserPhoneNumberContact creates a new phone number contact for a user. Requires ADMIN privileges or user id equals your current user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1phone-numbers/post
func (c *Client) CreateUserPhoneNumberContact(input *CreateUserPhoneNumberContactInput) (*CreateUserPhoneNumberContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserPhoneNumberContact == nil {
		return nil, errors.New("user phone number contact input is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/phone-numbers", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(input.UserPhoneNumberContact).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserPhoneNumberContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &CreateUserPhoneNumberContactOutput{UserPhoneNumberContact: contact}, nil
}

// GetUserPhoneNumberContactInput represents the input of a GetUserPhoneNumberContact operation.
type GetUserPhoneNumberContactInput struct {
	_                        struct{}
	UserID                   *int64
	UserPhoneNumberContactID *int64
}

// GetUserPhoneNumberContactOutput represents the output of a GetUserPhoneNumberContact operation.
type GetUserPhoneNumberContactOutput struct {
	_                      struct{}
	UserPhoneNumberContact *UserPhoneNumberContact
}

// GetUserPhoneNumberContact gets a phone number contact of a user by id. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1phone-numbers~1{id}/get
func (c *Client) GetUserPhoneNumberContact(input *GetUserPhoneNumberContactInput) (*GetUserPhoneNumberContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserPhoneNumberContactID == nil {
		return nil, errors.New("user phone number contact id is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/phone-numbers/%d", apiRoutes.users, *input.UserID, *input.UserPhoneNumberContactID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserPhoneNumberContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &GetUserPhoneNumberContactOutput{UserPhoneNumberContact: contact}, nil
}

// GetUserPhoneNumberContactsInput represents the input of a GetUserPhoneNumberContacts operation.
type GetUserPhoneNumberContactsInput struct {
	_      struct{}
	UserID *int64
}

// GetUserPhoneNumberContactsOutput represents the output of a GetUserPhoneNumberContacts operation.
type GetUserPhoneNumberContactsOutput struct {
	_                       struct{}
	UserPhoneNumberContacts []*UserPhoneNumberContact
}

// GetUserPhoneNumberContacts lists existing phone number contacts of a user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1phone-numbers/get
func (c *Client) GetUserPhoneNumberContacts(input *GetUserPhoneNumberContactsInput) (*GetUserPhoneNumberContactsOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/phone-numbers", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contacts := make([]*UserPhoneNumberContact, 0)
	err = json.Unmarshal(resp.Body(), &contacts)
	if err != nil {
		return nil, err
	}

	return &GetUserPhoneNumberContactsOutput{UserPhoneNumberContacts: contacts}, nil
}

// SearchUserPhoneNumberContactInput represents the input of a SearchUserPhoneNumberContact operation.
type SearchUserPhoneNumberContactInput struct {
	_                            struct{}
	UserID                       *int64
	UserPhoneNumberContactTarget *string
}

// SearchUserPhoneNumberContactOutput represents the output of a SearchUserPhoneNumberContact operation.
type SearchUserPhoneNumberContactOutput struct {
	_                      struct{}
	UserPhoneNumberContact *UserPhoneNumberContact
}

// SearchUserPhoneNumberContact gets the phone number contact with specified target of a user.
func (c *Client) SearchUserPhoneNumberContact(input *SearchUserPhoneNumberContactInput) (*SearchUserPhoneNumberContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserPhoneNumberContactTarget == nil {
		return nil, errors.New("user phone number contact target is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/phone-numbers/name", apiRoutes.users, *input.UserID)
	resp, err := c.httpClient.R().SetBody(UserPhoneNumberContact{Target: *input.UserPhoneNumberContactTarget}).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	contact := &UserPhoneNumberContact{}
	err = json.Unmarshal(resp.Body(), contact)
	if err != nil {
		return nil, err
	}

	return &SearchUserPhoneNumberContactOutput{UserPhoneNumberContact: contact}, nil
}

// UpdateUserPhoneNumberContactInput represents the input of a UpdateUserPhoneNumberContact operation.
type UpdateUserPhoneNumberContactInput struct {
	_                        struct{}
	UserID                   *int64
	UserPhoneNumberContactID *int64
	UserPhoneNumberContact   *UserPhoneNumberContact
}

// UpdateUserPhoneNumberContactOutput represents the output of a UpdateUserPhoneNumberContact operation.
type UpdateUserPhoneNumberContactOutput struct {
	_                      struct{}
	UserPhoneNumberContact *UserPhoneNumberContact
}

// UpdateUserPhoneNumberContact updates an existing phone number contact of a user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1phone-numbers~1{id}/put
func (c *Client) UpdateUserPhoneNumberContact(input *UpdateUserPhoneNumberContactInput) (*UpdateUserPhoneNumberContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserPhoneNumberContactID == nil {
		return nil, errors.New("user phone number contact id is required")
	}
	if input.UserPhoneNumberContact == nil {
		return nil, errors.New("user input is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/phone-numbers/%d", apiRoutes.users, *input.UserID, *input.UserPhoneNumberContactID)
	resp, err := c.httpClient.R().SetBody(input.UserPhoneNumberContact).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	user := &UserPhoneNumberContact{}
	err = json.Unmarshal(resp.Body(), user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserPhoneNumberContactOutput{UserPhoneNumberContact: user}, nil
}

// DeleteUserPhoneNumberContactInput represents the input of a DeleteUserPhoneNumberContact operation.
type DeleteUserPhoneNumberContactInput struct {
	_                        struct{}
	UserID                   *int64
	UserPhoneNumberContactID *int64
}

// DeleteUserPhoneNumberContactOutput represents the output of a DeleteUserPhoneNumberContact operation.
type DeleteUserPhoneNumberContactOutput struct {
	_ struct{}
}

// DeleteUserPhoneNumberContact deletes the specified phone number contact of a user. https://api.ilert.com/api-docs/#tag/Contacts/paths/~1users~1{user-id}~1contacts~1phone-numbers~1{id}/delete
func (c *Client) DeleteUserPhoneNumberContact(input *DeleteUserPhoneNumberContactInput) (*DeleteUserPhoneNumberContactOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.UserID == nil {
		return nil, errors.New("user id is required")
	}
	if input.UserPhoneNumberContactID == nil {
		return nil, errors.New("user phone number contact id is required")
	}

	url := fmt.Sprintf("%s/%d/contacts/phone-numbers/%d", apiRoutes.users, *input.UserID, *input.UserPhoneNumberContactID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteUserPhoneNumberContactOutput{}, nil
}
