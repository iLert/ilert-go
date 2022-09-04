package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Incident Template definition https://api.ilert.com/api-docs/#tag/Incident-Templates
type IncidentTemplate struct {
	ID               int64       `json:"id"`
	Name             string      `json:"name"`
	Summary          string      `json:"summary"`
	Status           string      `json:"status"`
	Message          string      `json:"message"`
	SendNotification bool        `json:"send_notification"`
	Teams            []TeamShort `json:"teams,omitempty"`
}

// CreateIncidentTemplateInput represents the input of a CreateIncidentTemplate operation.
type CreateIncidentTemplateInput struct {
	_                struct{}
	IncidentTemplate *IncidentTemplate
}

// CreateIncidentTemplateOutput represents the output of a CreateIncidentTemplate operation.
type CreateIncidentTemplateOutput struct {
	_                struct{}
	IncidentTemplate *IncidentTemplate
}

// CreateIncidentTemplate creates a new incidentTemplate. https://api.ilert.com/api-docs/#tag/Incident-Templates/paths/~1incident-templates/post
func (c *Client) CreateIncidentTemplate(input *CreateIncidentTemplateInput) (*CreateIncidentTemplateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentTemplate == nil {
		return nil, errors.New("incidentTemplate input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.IncidentTemplate).Post(apiRoutes.incidentTemplates)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	incidentTemplate := &IncidentTemplate{}
	err = json.Unmarshal(resp.Body(), incidentTemplate)
	if err != nil {
		return nil, err
	}

	return &CreateIncidentTemplateOutput{IncidentTemplate: incidentTemplate}, nil
}

// GetIncidentTemplatesInput represents the input of a GetIncidentTemplates operation.
type GetIncidentTemplatesInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults *int
}

// GetIncidentTemplatesOutput represents the output of a GetIncidentTemplates operation.
type GetIncidentTemplatesOutput struct {
	_                 struct{}
	IncidentTemplates []*IncidentTemplate
}

// GetIncidentTemplates lists incidentTemplate sources. https://api.ilert.com/api-docs/#tag/Incident-Templates/paths/~1incident-templates/get
func (c *Client) GetIncidentTemplates(input *GetIncidentTemplatesInput) (*GetIncidentTemplatesOutput, error) {
	if input == nil {
		input = &GetIncidentTemplatesInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.incidentTemplates, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incidentTemplates := make([]*IncidentTemplate, 0)
	err = json.Unmarshal(resp.Body(), &incidentTemplates)
	if err != nil {
		return nil, err
	}

	return &GetIncidentTemplatesOutput{IncidentTemplates: incidentTemplates}, nil
}

// GetIncidentTemplateInput represents the input of a GetIncidentTemplate operation.
type GetIncidentTemplateInput struct {
	_                  struct{}
	IncidentTemplateID *int64
}

// GetIncidentTemplateOutput represents the output of a GetIncidentTemplate operation.
type GetIncidentTemplateOutput struct {
	_                struct{}
	IncidentTemplate *IncidentTemplate
}

// GetIncidentTemplate gets a incidentTemplate by ID. https://api.ilert.com/api-docs/#tag/Incident-Templates/paths/~1incident-templates~1{id}/get
func (c *Client) GetIncidentTemplate(input *GetIncidentTemplateInput) (*GetIncidentTemplateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentTemplateID == nil {
		return nil, errors.New("incidentTemplate id is required")
	}

	q := url.Values{}

	var url = fmt.Sprintf("%s/%d?%s", apiRoutes.incidentTemplates, *input.IncidentTemplateID, q.Encode())

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incidentTemplate := &IncidentTemplate{}
	err = json.Unmarshal(resp.Body(), incidentTemplate)
	if err != nil {
		return nil, err
	}

	return &GetIncidentTemplateOutput{IncidentTemplate: incidentTemplate}, nil
}

// SearchIncidentTemplateInput represents the input of a SearchIncidentTemplate operation.
type SearchIncidentTemplateInput struct {
	_                    struct{}
	IncidentTemplateName *string
}

// SearchIncidentTemplateOutput represents the output of a SearchIncidentTemplate operation.
type SearchIncidentTemplateOutput struct {
	_                struct{}
	IncidentTemplate *IncidentTemplate
}

// SearchIncidentTemplate gets the incidentTemplate with specified name.
func (c *Client) SearchIncidentTemplate(input *SearchIncidentTemplateInput) (*SearchIncidentTemplateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentTemplateName == nil {
		return nil, errors.New("incident template name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.incidentTemplates, *input.IncidentTemplateName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incidentTemplate := &IncidentTemplate{}
	err = json.Unmarshal(resp.Body(), incidentTemplate)
	if err != nil {
		return nil, err
	}

	return &SearchIncidentTemplateOutput{IncidentTemplate: incidentTemplate}, nil
}

// UpdateIncidentTemplateInput represents the input of a UpdateIncidentTemplate operation.
type UpdateIncidentTemplateInput struct {
	_                  struct{}
	IncidentTemplateID *int64
	IncidentTemplate   *IncidentTemplate
}

// UpdateIncidentTemplateOutput represents the output of a UpdateIncidentTemplate operation.
type UpdateIncidentTemplateOutput struct {
	_                struct{}
	IncidentTemplate *IncidentTemplate
}

// UpdateIncidentTemplate updates the specific incidentTemplate. https://api.ilert.com/api-docs/#tag/Incident-Templates/paths/~1incident-templates~1{id}/put
func (c *Client) UpdateIncidentTemplate(input *UpdateIncidentTemplateInput) (*UpdateIncidentTemplateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentTemplateID == nil {
		return nil, errors.New("incidentTemplate id is required")
	}
	if input.IncidentTemplate == nil {
		return nil, errors.New("incidentTemplate input is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.incidentTemplates, *input.IncidentTemplateID)

	resp, err := c.httpClient.R().SetBody(input.IncidentTemplate).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	incidentTemplate := &IncidentTemplate{}
	err = json.Unmarshal(resp.Body(), incidentTemplate)
	if err != nil {
		return nil, err
	}

	return &UpdateIncidentTemplateOutput{IncidentTemplate: incidentTemplate}, nil
}

// DeleteIncidentTemplateInput represents the input of a DeleteIncidentTemplate operation.
type DeleteIncidentTemplateInput struct {
	_                  struct{}
	IncidentTemplateID *int64
}

// DeleteIncidentTemplateOutput represents the output of a DeleteIncidentTemplate operation.
type DeleteIncidentTemplateOutput struct {
	_ struct{}
}

// DeleteIncidentTemplate deletes the specified incidentTemplate. https://api.ilert.com/api-docs/#tag/Incident-Templates/paths/~1incident-templates~1{id}/delete
func (c *Client) DeleteIncidentTemplate(input *DeleteIncidentTemplateInput) (*DeleteIncidentTemplateOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.IncidentTemplateID == nil {
		return nil, errors.New("incidentTemplate id is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.incidentTemplates, *input.IncidentTemplateID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteIncidentTemplateOutput{}, nil
}
