package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// StatusPageGroup definition https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups/post
type StatusPageGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CreateStatusPageGroupInput represents the input of a CreateStatusPageGroup operation.
type CreateStatusPageGroupInput struct {
	_               struct{}
	StatusPageGroup *StatusPageGroup
	StatusPageID    *int64
}

// CreateStatusPageGroupOutput represents the output of a CreateStatusPageGroup operation.
type CreateStatusPageGroupOutput struct {
	_               struct{}
	StatusPageGroup *StatusPageGroup
}

// CreateStatusPageGroup creates a new status page group. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups/post
func (c *Client) CreateStatusPageGroup(input *CreateStatusPageGroupInput) (*CreateStatusPageGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageGroup == nil {
		return nil, errors.New("status page group input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups", apiRoutes.statusPages, *input.StatusPageID)
	resp, err := c.httpClient.R().SetBody(input.StatusPageGroup).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	StatusPageGroup := &StatusPageGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageGroup)
	if err != nil {
		return nil, err
	}

	return &CreateStatusPageGroupOutput{StatusPageGroup: StatusPageGroup}, nil
}

// GetStatusPageGroupInput represents the input of a GetStatusPageGroup operation.
type GetStatusPageGroupInput struct {
	_                 struct{}
	StatusPageGroupID *int64
	StatusPageID      *int64
}

// GetStatusPageGroupOutput represents the output of a GetStatusPageGroup operation.
type GetStatusPageGroupOutput struct {
	_               struct{}
	StatusPageGroup *StatusPageGroup
}

// GetStatusPageGroup gets the status page group with specified id. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups~1{group-id}/get
func (c *Client) GetStatusPageGroup(input *GetStatusPageGroupInput) (*GetStatusPageGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageGroupID == nil {
		return nil, errors.New("status page group id is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/%d", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageGroupID)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageGroup := &StatusPageGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageGroup)
	if err != nil {
		return nil, err
	}

	return &GetStatusPageGroupOutput{StatusPageGroup: StatusPageGroup}, nil
}

// GetStatusPageGroupsInput represents the input of a GetStatusPageGroups operation.
type GetStatusPageGroupsInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults   *int
	StatusPageID *int64
}

// GetStatusPageGroupsOutput represents the output of a GetStatusPageGroups operation.
type GetStatusPageGroupsOutput struct {
	_                struct{}
	StatusPageGroups []*StatusPageGroup
}

// GetStatusPageGroups lists existing status page groups. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups/get
func (c *Client) GetStatusPageGroups(input *GetStatusPageGroupsInput) (*GetStatusPageGroupsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	url := fmt.Sprintf("%s/%d/groups?%s", apiRoutes.statusPages, *input.StatusPageID, q.Encode())
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageGroups := make([]*StatusPageGroup, 0)
	err = json.Unmarshal(resp.Body(), &StatusPageGroups)
	if err != nil {
		return nil, err
	}

	return &GetStatusPageGroupsOutput{StatusPageGroups: StatusPageGroups}, nil
}

// SearchStatusPageGroupInput represents the input of a SearchStatusPageGroup operation.
type SearchStatusPageGroupInput struct {
	_                   struct{}
	StatusPageGroupName *string
	StatusPageID        *int64
}

// SearchStatusPageGroupOutput represents the output of a SearchStatusPageGroup operation.
type SearchStatusPageGroupOutput struct {
	_               struct{}
	StatusPageGroup *StatusPageGroup
}

// SearchStatusPageGroup gets the status page group with specified name.
func (c *Client) SearchStatusPageGroup(input *SearchStatusPageGroupInput) (*SearchStatusPageGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageGroupName == nil {
		return nil, errors.New("status page group name is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/name/%s", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageGroupName)
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageGroup := &StatusPageGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageGroup)
	if err != nil {
		return nil, err
	}

	return &SearchStatusPageGroupOutput{StatusPageGroup: StatusPageGroup}, nil
}

// UpdateStatusPageGroupInput represents the input of a UpdateStatusPageGroup operation.
type UpdateStatusPageGroupInput struct {
	_                 struct{}
	StatusPageGroup   *StatusPageGroup
	StatusPageGroupID *int64
	StatusPageID      *int64
}

// UpdateStatusPageGroupOutput represents the output of a UpdateStatusPageGroup operation.
type UpdateStatusPageGroupOutput struct {
	_               struct{}
	StatusPageGroup *StatusPageGroup
}

// UpdateStatusPageGroup updates an existing status page group. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups~1{group-id}/put
func (c *Client) UpdateStatusPageGroup(input *UpdateStatusPageGroupInput) (*UpdateStatusPageGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageGroup == nil {
		return nil, errors.New("status page group input is required")
	}
	if input.StatusPageGroupID == nil {
		return nil, errors.New("status page group id is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/%d", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageGroupID)
	resp, err := c.httpClient.R().SetBody(input.StatusPageGroup).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	StatusPageGroup := &StatusPageGroup{}
	err = json.Unmarshal(resp.Body(), StatusPageGroup)
	if err != nil {
		return nil, err
	}

	return &UpdateStatusPageGroupOutput{StatusPageGroup: StatusPageGroup}, nil
}

// DeleteStatusPageGroupInput represents the input of a DeleteStatusPageGroup operation.
type DeleteStatusPageGroupInput struct {
	_                 struct{}
	StatusPageGroupID *int64
	StatusPageID      *int64
}

// DeleteStatusPageGroupOutput represents the output of a DeleteStatusPageGroup operation.
type DeleteStatusPageGroupOutput struct {
	_ struct{}
}

// DeleteStatusPageGroup deletes the specified status page group. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1groups~1{group-id}/delete
func (c *Client) DeleteStatusPageGroup(input *DeleteStatusPageGroupInput) (*DeleteStatusPageGroupOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageGroupID == nil {
		return nil, errors.New("status page group id is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("status page id is required")
	}

	url := fmt.Sprintf("%s/%d/groups/%d", apiRoutes.statusPages, *input.StatusPageID, *input.StatusPageGroupID)
	resp, err := c.httpClient.R().Delete(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteStatusPageGroupOutput{}, nil
}
