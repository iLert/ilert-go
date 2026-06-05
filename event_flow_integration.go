package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// EventFlowIntegration definition
type EventFlowIntegration struct {
	ID              int64  `json:"id,omitempty"`
	IntegrationType string `json:"integrationType"` // one of AlertSourceIntegrationTypes
	IntegrationKey  string `json:"integrationKey,omitempty"`
	EventFlowID     int64  `json:"eventFlowId"`

	// Read-only fields populated by the server
	Catalogue      *IntegrationCatalogEntry `json:"catalogue,omitempty"`
	IntegrationURL string                   `json:"integrationUrl,omitempty"`
}

// IntegrationCatalogEntry describes the integration catalogue entry returned for an EventFlowIntegration.
type IntegrationCatalogEntry struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Slug             string   `json:"slug,omitempty"`
	Type             string   `json:"type,omitempty"` // one of IntegrationType
	DocumentationURL string   `json:"documentationUrl,omitempty"`
	IconURL          string   `json:"iconUrl,omitempty"`
	LargeIconURL     string   `json:"largeIconUrl,omitempty"`
	LightIconURL     string   `json:"lightIconUrl,omitempty"`
	DarkIconURL      string   `json:"darkIconUrl,omitempty"`
	IntegrationPath  string   `json:"integrationPath,omitempty"`
	Description      string   `json:"description,omitempty"`
	Categories       []string `json:"categories,omitempty"`
	Featured         *bool    `json:"featured,omitempty"`
}

// IntegrationType defines integration catalogue entry types
var IntegrationType = struct {
	Inbound  string
	Outbound string
	Native   string
}{
	Inbound:  "INBOUND",
	Outbound: "OUTBOUND",
	Native:   "NATIVE",
}

// IntegrationTypeAll defines all integration catalogue entry types
var IntegrationTypeAll = []string{
	IntegrationType.Inbound,
	IntegrationType.Outbound,
	IntegrationType.Native,
}

// CreateEventFlowIntegrationInput represents the input of a CreateEventFlowIntegration operation.
type CreateEventFlowIntegrationInput struct {
	_                    struct{}
	EventFlowIntegration *EventFlowIntegration
}

// CreateEventFlowIntegrationOutput represents the output of a CreateEventFlowIntegration operation.
type CreateEventFlowIntegrationOutput struct {
	_                    struct{}
	EventFlowIntegration *EventFlowIntegration
}

// CreateEventFlowIntegration creates a new event flow integration.
func (c *Client) CreateEventFlowIntegration(input *CreateEventFlowIntegrationInput) (*CreateEventFlowIntegrationOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowIntegration == nil {
		return nil, errors.New("event flow integration input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.EventFlowIntegration).Post(apiRoutes.eventFlowIntegrations)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	integration := &EventFlowIntegration{}
	if err = json.Unmarshal(resp.Body(), integration); err != nil {
		return nil, err
	}

	return &CreateEventFlowIntegrationOutput{EventFlowIntegration: integration}, nil
}

// GetEventFlowIntegrationInput represents the input of a GetEventFlowIntegration operation.
type GetEventFlowIntegrationInput struct {
	_                      struct{}
	EventFlowIntegrationID *int64
}

// GetEventFlowIntegrationOutput represents the output of a GetEventFlowIntegration operation.
type GetEventFlowIntegrationOutput struct {
	_                    struct{}
	EventFlowIntegration *EventFlowIntegration
}

// GetEventFlowIntegration gets the event flow integration with the specified id.
func (c *Client) GetEventFlowIntegration(input *GetEventFlowIntegrationInput) (*GetEventFlowIntegrationOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowIntegrationID == nil {
		return nil, errors.New("event flow integration id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.eventFlowIntegrations, *input.EventFlowIntegrationID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	integration := &EventFlowIntegration{}
	if err = json.Unmarshal(resp.Body(), integration); err != nil {
		return nil, err
	}

	return &GetEventFlowIntegrationOutput{EventFlowIntegration: integration}, nil
}

// GetEventFlowIntegrationsInput represents the input of a GetEventFlowIntegrations operation.
type GetEventFlowIntegrationsInput struct {
	_ struct{}

	// Optional: filter the result to integrations belonging to the given event flow id.
	EventFlowID *int64

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetEventFlowIntegrationsOutput represents the output of a GetEventFlowIntegrations operation.
type GetEventFlowIntegrationsOutput struct {
	_                     struct{}
	EventFlowIntegrations []*EventFlowIntegration
}

// GetEventFlowIntegrations lists existing event flow integrations.
func (c *Client) GetEventFlowIntegrations(input *GetEventFlowIntegrationsInput) (*GetEventFlowIntegrationsOutput, error) {
	q := url.Values{}
	if input != nil && input.EventFlowID != nil {
		q.Add("event-flow", strconv.FormatInt(*input.EventFlowID, 10))
	}
	if input != nil && input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input != nil && input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.eventFlowIntegrations, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	integrations := make([]*EventFlowIntegration, 0)
	if err = json.Unmarshal(resp.Body(), &integrations); err != nil {
		return nil, err
	}

	return &GetEventFlowIntegrationsOutput{EventFlowIntegrations: integrations}, nil
}

// UpdateEventFlowIntegrationInput represents the input of an UpdateEventFlowIntegration operation.
type UpdateEventFlowIntegrationInput struct {
	_                      struct{}
	EventFlowIntegrationID *int64
	EventFlowIntegration   *EventFlowIntegration
}

// UpdateEventFlowIntegrationOutput represents the output of an UpdateEventFlowIntegration operation.
type UpdateEventFlowIntegrationOutput struct {
	_                    struct{}
	EventFlowIntegration *EventFlowIntegration
}

// UpdateEventFlowIntegration updates the event flow integration with the specified id.
func (c *Client) UpdateEventFlowIntegration(input *UpdateEventFlowIntegrationInput) (*UpdateEventFlowIntegrationOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowIntegrationID == nil {
		return nil, errors.New("event flow integration id is required")
	}
	if input.EventFlowIntegration == nil {
		return nil, errors.New("event flow integration input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.EventFlowIntegration).Put(fmt.Sprintf("%s/%d", apiRoutes.eventFlowIntegrations, *input.EventFlowIntegrationID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	integration := &EventFlowIntegration{}
	if err = json.Unmarshal(resp.Body(), integration); err != nil {
		return nil, err
	}

	return &UpdateEventFlowIntegrationOutput{EventFlowIntegration: integration}, nil
}

// DeleteEventFlowIntegrationInput represents the input of a DeleteEventFlowIntegration operation.
type DeleteEventFlowIntegrationInput struct {
	_                      struct{}
	EventFlowIntegrationID *int64
}

// DeleteEventFlowIntegrationOutput represents the output of a DeleteEventFlowIntegration operation.
type DeleteEventFlowIntegrationOutput struct {
	_ struct{}
}

// DeleteEventFlowIntegration deletes the event flow integration with the specified id.
func (c *Client) DeleteEventFlowIntegration(input *DeleteEventFlowIntegrationInput) (*DeleteEventFlowIntegrationOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowIntegrationID == nil {
		return nil, errors.New("event flow integration id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.eventFlowIntegrations, *input.EventFlowIntegrationID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteEventFlowIntegrationOutput{}, nil
}
