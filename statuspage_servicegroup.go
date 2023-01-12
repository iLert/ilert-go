package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// StatusPage ServiceGroup definition https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups/post
type StatusPageServiceGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CreateStatusPageServiceGroupInput represents the input of a CreateStatusPageServiceGroup operation.
type CreateStatusPageServiceGroupInput struct {
	_                      struct{}
	StatusPageServiceGroup *StatusPageServiceGroup
	StatusPageID           *int64
}

// CreateStatusPageServiceGroupOutput represents the output of a CreateStatusPageServiceGroup operation.
type CreateStatusPageServiceGroupOutput struct {
	_                      struct{}
	StatusPageServiceGroup *StatusPageServiceGroup
}

// CreateStatusPageServiceGroup creates a new StatusPageServiceGroup. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups/post
func (c *Client) CreateStatusPageServiceGroup(input *CreateStatusPageServiceGroupInput) (*CreateStatusPageServiceGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageServiceGroup == nil {
		return nil, errors.New("status page service group input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups", apiRoutes.statusPages, *input.StatusPageID)
	resp, err := c.httpClient.R().SetBody(input.StatusPageServiceGroup).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	StatusPageServiceGroup := &StatusPageServiceGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageServiceGroup)
	if err != nil {
		return nil, err
	}

	return &CreateStatusPageServiceGroupOutput{StatusPageServiceGroup: StatusPageServiceGroup}, nil
}

// GetStatusPageServiceGroupInput represents the input of a GetStatusPageServiceGroup operation.
type GetStatusPageServiceGroupInput struct {
	_                        struct{}
	StatusPageServiceGroupID *int64
	StatusPageID             *int64
}

// GetStatusPageServiceGroupOutput represents the output of a GetStatusPageServiceGroup operation.
type GetStatusPageServiceGroupOutput struct {
	_                      struct{}
	StatusPageServiceGroup *StatusPageServiceGroup
}

// GetStatusPageServiceGroup gets the StatusPageServiceGroup with specified id. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups~1{group-id}/get
func (c *Client) GetStatusPageServiceGroup(input *GetStatusPageServiceGroupInput) (*GetStatusPageServiceGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageServiceGroupID == nil {
		return nil, errors.New("status page service group id is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/%d", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageServiceGroupID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageServiceGroup := &StatusPageServiceGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageServiceGroup)
	if err != nil {
		return nil, err
	}

	return &GetStatusPageServiceGroupOutput{StatusPageServiceGroup: StatusPageServiceGroup}, nil
}

// GetStatusPageServiceGroupsInput represents the input of a GetStatusPageServiceGroups operation.
type GetStatusPageServiceGroupsInput struct {
	_            struct{}
	StatusPageID *int64
}

// GetStatusPageServiceGroupsOutput represents the output of a GetStatusPageServiceGroups operation.
type GetStatusPageServiceGroupsOutput struct {
	_                       struct{}
	StatusPageServiceGroups []*StatusPageServiceGroup
}

// GetStatusPageServiceGroups gets list of StatusPageServiceGroups. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups/get
func (c *Client) GetStatusPageServiceGroups(input *GetStatusPageServiceGroupsInput) (*GetStatusPageServiceGroupsOutput, error) {
	url := fmt.Sprintf("%s/%d/groups", apiRoutes.statusPages, *input.StatusPageID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageServiceGroups := make([]*StatusPageServiceGroup, 0)
	err = json.Unmarshal(resp.Body(), &StatusPageServiceGroups)
	if err != nil {
		return nil, err
	}

	return &GetStatusPageServiceGroupsOutput{StatusPageServiceGroups: StatusPageServiceGroups}, nil
}

// SearchStatusPageServiceGroupInput represents the input of a SearchStatusPageServiceGroup operation.
type SearchStatusPageServiceGroupInput struct {
	_                          struct{}
	StatusPageServiceGroupName *string
	StatusPageID               *int64
}

// SearchStatusPageServiceGroupOutput represents the output of a SearchStatusPageServiceGroup operation.
type SearchStatusPageServiceGroupOutput struct {
	_                      struct{}
	StatusPageServiceGroup *StatusPageServiceGroup
}

// SearchStatusPageServiceGroup gets the StatusPageServiceGroup with specified name.
func (c *Client) SearchStatusPageServiceGroup(input *SearchStatusPageServiceGroupInput) (*SearchStatusPageServiceGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageServiceGroupName == nil {
		return nil, errors.New("status page service group name is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/name/%s", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageServiceGroupName)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageServiceGroup := &StatusPageServiceGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageServiceGroup)
	if err != nil {
		return nil, err
	}

	return &SearchStatusPageServiceGroupOutput{StatusPageServiceGroup: StatusPageServiceGroup}, nil
}

// UpdateStatusPageServiceGroupInput represents the input of a UpdateStatusPageServiceGroup operation.
type UpdateStatusPageServiceGroupInput struct {
	_                        struct{}
	StatusPageServiceGroup   *StatusPageServiceGroup
	StatusPageServiceGroupID *int64
	StatusPageID             *int64
}

// UpdateStatusPageServiceGroupOutput represents the output of a UpdateStatusPageServiceGroup operation.
type UpdateStatusPageServiceGroupOutput struct {
	_                      struct{}
	StatusPageServiceGroup *StatusPageServiceGroup
}

// UpdateStatusPageServiceGroup updates an existing StatusPageServiceGroup. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups~1{group-id}/put
func (c *Client) UpdateStatusPageServiceGroup(input *UpdateStatusPageServiceGroupInput) (*UpdateStatusPageServiceGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageServiceGroup == nil {
		return nil, errors.New("status page service group input is required")
	}
	if input.StatusPageServiceGroupID == nil {
		return nil, errors.New("status page service group id is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/%d", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageServiceGroupID)
	resp, err := c.httpClient.R().SetBody(input.StatusPageServiceGroup).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageServiceGroup := &StatusPageServiceGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageServiceGroup)
	if err != nil {
		return nil, err
	}

	return &UpdateStatusPageServiceGroupOutput{StatusPageServiceGroup: StatusPageServiceGroup}, nil
}

// DeleteStatusPageServiceGroupInput represents the input of a DeleteStatusPageServiceGroup operation.
type DeleteStatusPageServiceGroupInput struct {
	_                        struct{}
	StatusPageServiceGroupID *int64
	StatusPageID             *int64
}

// DeleteStatusPageServiceGroupOutput represents the output of a DeleteStatusPageServiceGroup operation.
type DeleteStatusPageServiceGroupOutput struct {
	_ struct{}
}

// DeleteStatusPageServiceGroup deletes the specified StatusPageServiceGroup. https://api.ilert.com/api-docs/#tag/StatusPageServiceGroups/paths/~1StatusPageServiceGroups~1{id}/delete
func (c *Client) DeleteStatusPageServiceGroup(input *DeleteStatusPageServiceGroupInput) (*DeleteStatusPageServiceGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageServiceGroupID == nil {
		return nil, errors.New("status page service group id is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/%d", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageServiceGroupID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteStatusPageServiceGroupOutput{}, nil
}
