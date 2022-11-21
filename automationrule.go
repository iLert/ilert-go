package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
type AutomationRule struct {
	ID               string            `json:"id"`
	AlertType        string            `json:"alertType"`
	ResolveIncident  bool              `json:"resolveIncident"`
	ResolveService   bool              `json:"resolveService"`
	ServiceStatus    string            `json:"serviceStatus"`
	Template         *IncidentTemplate `json:"template"`
	Service          *Service          `json:"service"`
	AlertSource      *AlertSource      `json:"alertSource"`
	SendNotification bool              `json:"sendNotification"`
}

// AlertType defines the alert type in an automation rule
var AlertType = struct {
	Created  string
	Accepted string
}{
	Created:  "CREATED",
	Accepted: "ACCEPTED",
}

// AlertType defines the alert type list
var AlertTypeAll = []string{
	AlertType.Created,
	AlertType.Accepted,
}

// CreateAutomationRuleInput represents the input of a CreateAutomationRule operation.
type CreateAutomationRuleInput struct {
	_              struct{}
	AutomationRule *AutomationRule
}

// CreateAutomationRuleOutput represents the output of a CreateAutomationRule operation.
type CreateAutomationRuleOutput struct {
	_              struct{}
	AutomationRule *AutomationRule
}

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
func (c *Client) CreateAutomationRule(input *CreateAutomationRuleInput) (*CreateAutomationRuleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AutomationRule == nil {
		return nil, errors.New("automationRule input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.AutomationRule).Post(apiRoutes.automationRules)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	automationRule := &AutomationRule{}
	err = json.Unmarshal(resp.Body(), automationRule)
	if err != nil {
		return nil, err
	}

	return &CreateAutomationRuleOutput{AutomationRule: automationRule}, nil
}

// GetAutomationRulesInput represents the input of a GetAutomationRules operation.
type GetAutomationRulesInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults *int

	// The service id for which automation rules are filtered for, this param is required
	Service *int
}

// GetAutomationRulesOutput represents the output of a GetAutomationRules operation.
type GetAutomationRulesOutput struct {
	_               struct{}
	AutomationRules []*AutomationRule
}

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
func (c *Client) GetAutomationRules(input *GetAutomationRulesInput) (*GetAutomationRulesOutput, error) {
	if input == nil {
		input = &GetAutomationRulesInput{}
	}
	if input.Service == nil {
		return nil, errors.New("service id is required")
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}
	if input.Service != nil {
		q.Add("service", strconv.Itoa(*input.Service))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.automationRules, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	automationRules := make([]*AutomationRule, 0)
	err = json.Unmarshal(resp.Body(), &automationRules)
	if err != nil {
		return nil, err
	}

	return &GetAutomationRulesOutput{AutomationRules: automationRules}, nil
}

// GetAutomationRuleInput represents the input of a GetAutomationRule operation.
type GetAutomationRuleInput struct {
	_                struct{}
	AutomationRuleID *string
}

// GetAutomationRuleOutput represents the output of a GetAutomationRule operation.
type GetAutomationRuleOutput struct {
	_              struct{}
	AutomationRule *AutomationRule
}

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
func (c *Client) GetAutomationRule(input *GetAutomationRuleInput) (*GetAutomationRuleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AutomationRuleID == nil {
		return nil, errors.New("automationRule id is required")
	}

	q := url.Values{}

	var url = fmt.Sprintf("%s/%s?%s", apiRoutes.automationRules, *input.AutomationRuleID, q.Encode())

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	automationRule := &AutomationRule{}
	err = json.Unmarshal(resp.Body(), automationRule)
	if err != nil {
		return nil, err
	}

	return &GetAutomationRuleOutput{AutomationRule: automationRule}, nil
}

// UpdateAutomationRuleInput represents the input of a UpdateAutomationRule operation.
type UpdateAutomationRuleInput struct {
	_                struct{}
	AutomationRuleID *string
	AutomationRule   *AutomationRule
}

// UpdateAutomationRuleOutput represents the output of a UpdateAutomationRule operation.
type UpdateAutomationRuleOutput struct {
	_              struct{}
	AutomationRule *AutomationRule
}

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
func (c *Client) UpdateAutomationRule(input *UpdateAutomationRuleInput) (*UpdateAutomationRuleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AutomationRuleID == nil {
		return nil, errors.New("automationRule id is required")
	}
	if input.AutomationRule == nil {
		return nil, errors.New("automationRule input is required")
	}

	url := fmt.Sprintf("%s/%s", apiRoutes.automationRules, *input.AutomationRuleID)

	resp, err := c.httpClient.R().SetBody(input.AutomationRule).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	automationRule := &AutomationRule{}
	err = json.Unmarshal(resp.Body(), automationRule)
	if err != nil {
		return nil, err
	}

	return &UpdateAutomationRuleOutput{AutomationRule: automationRule}, nil
}

// DeleteAutomationRuleInput represents the input of a DeleteAutomationRule operation.
type DeleteAutomationRuleInput struct {
	_                struct{}
	AutomationRuleID *string
}

// DeleteAutomationRuleOutput represents the output of a DeleteAutomationRule operation.
type DeleteAutomationRuleOutput struct {
	_ struct{}
}

// Legacy API - please use alert-actions of type 'automation_rule' - for more information see https://api.ilert.com/api-docs/#tag/Connectors/paths/~1connectors/post
func (c *Client) DeleteAutomationRule(input *DeleteAutomationRuleInput) (*DeleteAutomationRuleOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AutomationRuleID == nil {
		return nil, errors.New("automationRule id is required")
	}

	url := fmt.Sprintf("%s/%s", apiRoutes.automationRules, *input.AutomationRuleID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteAutomationRuleOutput{}, nil
}
