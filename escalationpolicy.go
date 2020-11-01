package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EscalationPolicy definition https://api.ilert.com/api-docs/#!/Escalation_Policies
type EscalationPolicy struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	EscalationRules []EscalationRule `json:"escalationRules"`
	Repeating       bool             `json:"repeating"`
	Frequency       int              `json:"frequency"`
}

// EscalationRule definition
type EscalationRule struct {
	User              *User     `json:"user,omitempty"`
	Schedule          *Schedule `json:"schedule,omitempty"`
	EscalationTimeout int       `json:"escalationTimeout"`
}

// CreateEscalationPolicyInput represents the input of a CreateEscalationPolicy operation.
type CreateEscalationPolicyInput struct {
	_                struct{}
	EscalationPolicy *EscalationPolicy
}

// CreateEscalationPolicyOutput represents the output of a CreateEscalationPolicy operation.
type CreateEscalationPolicyOutput struct {
	_                struct{}
	EscalationPolicy *EscalationPolicy
}

// CreateEscalationPolicy creates a new escalation policy. https://api.ilert.com/api-docs/#tag/Escalation-Policies/paths/~1escalation-policies/post
func (c *Client) CreateEscalationPolicy(input *CreateEscalationPolicyInput) (*CreateEscalationPolicyOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.EscalationPolicy == nil {
		return nil, errors.New("Escalation policy input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.EscalationPolicy).Post(fmt.Sprintf("%s", apiRoutes.escalationPolicies))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 201); err != nil {
		return nil, err
	}

	escalationPolicy := &EscalationPolicy{}
	err = json.Unmarshal(resp.Body(), escalationPolicy)
	if err != nil {
		return nil, err
	}

	output := &CreateEscalationPolicyOutput{EscalationPolicy: escalationPolicy}

	return output, nil
}

// GetEscalationPolicyInput represents the input of a GetEscalationPolicy operation.
type GetEscalationPolicyInput struct {
	_                  struct{}
	EscalationPolicyID *int64
}

// GetEscalationPolicyOutput represents the output of a GetEscalationPolicy operation.
type GetEscalationPolicyOutput struct {
	_                struct{}
	EscalationPolicy *EscalationPolicy
}

// GetEscalationPolicy gets the escalation policy with specified id. https://api.ilert.com/api-docs/#tag/Escalation-Policies/paths/~1escalation-policies~1{id}/get
func (c *Client) GetEscalationPolicy(input *GetEscalationPolicyInput) (*GetEscalationPolicyOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.EscalationPolicyID == nil {
		return nil, errors.New("EscalationPolicy id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.escalationPolicies, *input.EscalationPolicyID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	escalationPolicy := &EscalationPolicy{}
	err = json.Unmarshal(resp.Body(), escalationPolicy)
	if err != nil {
		return nil, err
	}

	output := &GetEscalationPolicyOutput{
		EscalationPolicy: escalationPolicy,
	}

	return output, nil
}

// GetEscalationPoliciesInput represents the input of a GetEscalationPolicies operation.
type GetEscalationPoliciesInput struct {
	_ struct{}
}

// GetEscalationPoliciesOutput represents the output of a GetEscalationPolicies operation.
type GetEscalationPoliciesOutput struct {
	_                  struct{}
	EscalationPolicies []*EscalationPolicy
}

// GetEscalationPolicies lists escalation policies. https://api.ilert.com/api-docs/#tag/Escalation-Policies/paths/~1escalation-policies/get
func (c *Client) GetEscalationPolicies(input *GetEscalationPoliciesInput) (*GetEscalationPoliciesOutput, error) {
	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s", apiRoutes.escalationPolicies))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	escalationPolicies := make([]*EscalationPolicy, 0)
	err = json.Unmarshal(resp.Body(), &escalationPolicies)
	if err != nil {
		return nil, err
	}

	output := &GetEscalationPoliciesOutput{EscalationPolicies: escalationPolicies}

	return output, nil
}

// UpdateEscalationPolicyInput represents the input of a UpdateEscalationPolicy operation.
type UpdateEscalationPolicyInput struct {
	_                  struct{}
	EscalationPolicyID *int64
	EscalationPolicy   *EscalationPolicy
}

// UpdateEscalationPolicyOutput represents the output of a UpdateEscalationPolicy operation.
type UpdateEscalationPolicyOutput struct {
	_                struct{}
	EscalationPolicy *EscalationPolicy
}

// UpdateEscalationPolicy updates an existing escalation policy. https://api.ilert.com/api-docs/#tag/Escalation-Policies/paths/~1escalation-policies~1{id}/put
func (c *Client) UpdateEscalationPolicy(input *UpdateEscalationPolicyInput) (*UpdateEscalationPolicyOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.EscalationPolicy == nil {
		return nil, errors.New("EscalationPolicy input is required")
	}
	if input.EscalationPolicyID == nil {
		return nil, errors.New("Escalation policy id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.EscalationPolicy).Put(fmt.Sprintf("%s/%d", apiRoutes.escalationPolicies, *input.EscalationPolicyID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	escalationPolicy := &EscalationPolicy{}
	err = json.Unmarshal(resp.Body(), escalationPolicy)
	if err != nil {
		return nil, err
	}

	output := &UpdateEscalationPolicyOutput{EscalationPolicy: escalationPolicy}

	return output, nil
}

// DeleteEscalationPolicyInput represents the input of a DeleteEscalationPolicy operation.
type DeleteEscalationPolicyInput struct {
	_                  struct{}
	EscalationPolicyID *int64
}

// DeleteEscalationPolicyOutput represents the output of a DeleteEscalationPolicy operation.
type DeleteEscalationPolicyOutput struct {
	_ struct{}
}

// DeleteEscalationPolicy deletes the specified escalation policy. https://api.ilert.com/api-docs/#tag/Escalation-Policies/paths/~1escalation-policies~1{id}/delete
func (c *Client) DeleteEscalationPolicy(input *DeleteEscalationPolicyInput) (*DeleteEscalationPolicyOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.EscalationPolicyID == nil {
		return nil, errors.New("EscalationPolicy id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.escalationPolicies, *input.EscalationPolicyID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 204); err != nil {
		return nil, err
	}

	output := &DeleteEscalationPolicyOutput{}
	return output, nil
}
