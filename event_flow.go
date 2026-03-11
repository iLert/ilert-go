package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// EventFlow definition https://api.ilert.com/api-docs/#tag/event-flows
type EventFlow struct {
	ID       int64          `json:"id,omitempty"`
	Name     string         `json:"name"`
	Teams    []TeamShort    `json:"teams,omitempty"`
	RootNode *EventFlowNode `json:"root"`
}

// EventFlow definition https://api.ilert.com/api-docs/#tag/event-flows
type EventFlowOutput struct {
	ID       int64                `json:"id"`
	Name     string               `json:"name"`
	Teams    []TeamShort          `json:"teams,omitempty"`
	RootNode *EventFlowNodeOutput `json:"root"`
}

type EventFlowNode struct {
	ID       int64             `json:"id,omitempty"`
	Name     string            `json:"name,omitempty"`
	NodeType string            `json:"nodeType"` // one of EventFlowNodeType
	Metadata any               `json:"metadata,omitempty"`
	Branches []EventFlowBranch `json:"branches"`
}

type EventFlowNodeOutput struct {
	ID       int64                  `json:"id"`
	Name     string                 `json:"name,omitempty"`
	NodeType string                 `json:"nodeType"`
	Metadata *EventFlowNodeMetadata `json:"metadata,omitempty"`
	Branches []EventFlowBranch      `json:"branches"`
}

type EventFlowBranch struct {
	ID         int64          `json:"id,omitempty"`
	BranchType string         `json:"branchType"` // one of EventFlowBranchType
	Condition  string         `json:"condition,omitempty"`
	Target     *EventFlowNode `json:"target,omitempty"`
}

type EventFlowNodeMetadata struct {
	VarKey                  string                      `json:"varKey,omitempty"`                  // PLAIN
	VarValue                string                      `json:"varValue,omitempty"`                // PLAIN
	SupportHoursID          *int64                      `json:"supportHoursId,omitempty"`          // SUPPORT_HOURS
	AlertSourceID           *int64                      `json:"alertSourceId,omitempty"`           // ROUTE_EVENT
	OverwritePriority       string                      `json:"overwritePriority,omitempty"`       // ROUTE_EVENT, one of EventFlowNodeMetadataOverwritePriority
	EscalationPolicyID      *int64                      `json:"escalationPolicyId,omitempty"`      // ROUTE_EVENT
	Definitions             []EventFlowNodeDefinition   `json:"definitions,omitempty"`             // DEFINE_BRANCHES
	WaitForDuration         string                      `json:"waitForDuration,omitempty"`         // WAIT
	WaitStartSupportHoursID *int64                      `json:"waitStartSupportHoursId,omitempty"` // WAIT
	WaitEndSupportHoursID   *int64                      `json:"waitEndSupportHoursId,omitempty"`   // WAIT
	Condition               string                      `json:"condition,omitempty"`               // TRANSFORM
	Rules                   []EventFlowNodeRuleMetadata `json:"rules,omitempty"`                   // TRANSFORM
}

type EventFlowNodeDefinition struct {
	BranchName string `json:"branchName"`
	Conditions string `json:"conditions"`
}

type EventFlowNodeRuleMetadata struct {
	Name       string               `json:"name"`
	Target     string               `json:"target"`
	Operator   string               `json:"operator"` // one of EventFlowNodeRuleOperator
	Value      any                  `json:"value,omitempty"`
	Source     string               `json:"source,omitempty"`
	Mapping    map[string]*string   `json:"mapping,omitempty"`
	Default    any                  `json:"default,omitempty"`
	Properties map[string]*string   `json:"properties,omitempty"`
	Items      []map[string]*string `json:"items,omitempty"`
}

var EventFlowNodeType = struct {
	Root           string
	Plain          string
	SupportHours   string
	RouteEvent     string
	DefineBranches string
	Wait           string
	Transform      string
}{
	Root:           "ROOT",
	Plain:          "PLAIN",
	SupportHours:   "SUPPORT_HOURS",
	RouteEvent:     "ROUTE_EVENT",
	DefineBranches: "DEFINE_BRANCHES",
	Wait:           "WAIT",
	Transform:      "TRANSFORM",
}

var EventFlowNodeTypeAll = []string{
	EventFlowNodeType.Root,
	EventFlowNodeType.Plain,
	EventFlowNodeType.SupportHours,
	EventFlowNodeType.RouteEvent,
	EventFlowNodeType.DefineBranches,
	EventFlowNodeType.Wait,
	EventFlowNodeType.Transform,
}

var EventFlowBranchType = struct {
	Branch   string
	CatchAll string
	Accepted string // only used for Root node
}{
	Branch:   "BRANCH",
	CatchAll: "CATCH_ALL",
	Accepted: "ACCEPTED",
}

var EventFlowBranchTypeAll = []string{
	EventFlowBranchType.Branch,
	EventFlowBranchType.CatchAll,
	EventFlowBranchType.Accepted,
}

var EventFlowNodeMetadataOverwritePriority = struct {
	High string
	Low  string
}{
	High: "HIGH",
	Low:  "LOW",
}

var EventFlowNodeMetadataOverwritePriorityAll = []string{
	EventFlowNodeMetadataOverwritePriority.High,
	EventFlowNodeMetadataOverwritePriority.Low,
}

var EventFlowNodeRuleOperator = struct {
	Set         string
	Copy        string
	Map         string
	Template    string
	Merge       string
	AppendArray string
}{
	Set:         "SET",
	Copy:        "COPY",
	Map:         "MAP",
	Template:    "TEMPLATE",
	Merge:       "MERGE",
	AppendArray: "APPEND_ARRAY",
}

var EventFlowNodeRuleOperatorAll = []string{
	EventFlowNodeRuleOperator.Set,
	EventFlowNodeRuleOperator.Copy,
	EventFlowNodeRuleOperator.Map,
	EventFlowNodeRuleOperator.Template,
	EventFlowNodeRuleOperator.Merge,
	EventFlowNodeRuleOperator.AppendArray,
}

// CreateEventFlowInput represents the input of a CreateEventFlow operation.
type CreateEventFlowInput struct {
	_         struct{}
	EventFlow *EventFlow
}

// CreateEventFlowOutput represents the output of a CreateEventFlow operation.
type CreateEventFlowOutput struct {
	_         struct{}
	EventFlow *EventFlowOutput
}

// CreateEventFlow creates a new event flow resource. https://api.ilert.com/api-docs/#tag/event-flows/post/event-flows
func (c *Client) CreateEventFlow(input *CreateEventFlowInput) (*CreateEventFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlow == nil {
		return nil, errors.New("event flow input is required")
	}

	resp, err := c.httpClient.R().SetBody(input.EventFlow).Post(fmt.Sprintf("%s", apiRoutes.eventFlows))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	eventFlow := &EventFlowOutput{}
	err = json.Unmarshal(resp.Body(), eventFlow)
	if err != nil {
		return nil, err
	}

	return &CreateEventFlowOutput{EventFlow: eventFlow}, nil
}

// GetEventFlowInput represents the input of a GetEventFlow operation.
type GetEventFlowInput struct {
	_           struct{}
	EventFlowID *int64
}

// GetEventFlowOutput represents the output of a GetEventFlow operation.
type GetEventFlowOutput struct {
	_         struct{}
	EventFlow *EventFlowOutput
}

// GetEventFlow gets the event flow resource with specified id. https://api.ilert.com/api-docs/#tag/event-flows/get/event-flows/{id}
func (c *Client) GetEventFlow(input *GetEventFlowInput) (*GetEventFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowID == nil {
		return nil, errors.New("event flow id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.eventFlows, *input.EventFlowID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	eventFlow := &EventFlowOutput{}
	err = json.Unmarshal(resp.Body(), eventFlow)
	if err != nil {
		return nil, err
	}

	return &GetEventFlowOutput{EventFlow: eventFlow}, nil
}

// GetEventFlowsInput represents the input of a GetEventFlows operation.
type GetEventFlowsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetEventFlowsOutput represents the output of a GetEventFlows operation.
type GetEventFlowsOutput struct {
	_          struct{}
	EventFlows []*EventFlowOutput
}

// GetEventFlows lists existing event flow resources. https://api.ilert.com/api-docs/#tag/event-flows/get/event-flows
func (c *Client) GetEventFlows(input *GetEventFlowsInput) (*GetEventFlowsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.eventFlows, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	eventFlows := make([]*EventFlowOutput, 0)
	err = json.Unmarshal(resp.Body(), &eventFlows)
	if err != nil {
		return nil, err
	}

	return &GetEventFlowsOutput{EventFlows: eventFlows}, nil
}

// SearchEventFlowInput represents the input of a SearchEventFlow operation.
type SearchEventFlowInput struct {
	_             struct{}
	EventFlowName *string
}

// SearchEventFlowOutput represents the output of a SearchEventFlow operation.
type SearchEventFlowOutput struct {
	_         struct{}
	EventFlow *EventFlowOutput
}

// SearchEventFlow gets the event flow resource with specified name.
func (c *Client) SearchEventFlow(input *SearchEventFlowInput) (*SearchEventFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowName == nil {
		return nil, errors.New("event flow name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.eventFlows, *input.EventFlowName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	eventFlow := &EventFlowOutput{}
	err = json.Unmarshal(resp.Body(), eventFlow)
	if err != nil {
		return nil, err
	}

	return &SearchEventFlowOutput{EventFlow: eventFlow}, nil
}

// UpdateEventFlowInput represents the input of a UpdateEventFlow operation.
type UpdateEventFlowInput struct {
	_           struct{}
	EventFlowID *int64
	EventFlow   *EventFlow
}

// UpdateEventFlowOutput represents the output of a UpdateEventFlow operation.
type UpdateEventFlowOutput struct {
	_         struct{}
	EventFlow *EventFlowOutput
}

// UpdateEventFlow updates an existing event flow resource. https://api.ilert.com/api-docs/#tag/event-flows/put/event-flows/{id}
func (c *Client) UpdateEventFlow(input *UpdateEventFlowInput) (*UpdateEventFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlow == nil {
		return nil, errors.New("event flow input is required")
	}
	if input.EventFlowID == nil {
		return nil, errors.New("event flow id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.EventFlow).Put(fmt.Sprintf("%s/%d", apiRoutes.eventFlows, *input.EventFlowID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	eventFlow := &EventFlowOutput{}
	err = json.Unmarshal(resp.Body(), eventFlow)
	if err != nil {
		return nil, err
	}

	return &UpdateEventFlowOutput{EventFlow: eventFlow}, nil
}

// DeleteEventFlowInput represents the input of a DeleteEventFlow operation.
type DeleteEventFlowInput struct {
	_           struct{}
	EventFlowID *int64
}

// DeleteEventFlowOutput represents the output of a DeleteEventFlow operation.
type DeleteEventFlowOutput struct {
	_ struct{}
}

// DeleteEventFlow deletes the specified event flow resource. https://api.ilert.com/api-docs/#tag/event-flows/delete/event-flows/{id}
func (c *Client) DeleteEventFlow(input *DeleteEventFlowInput) (*DeleteEventFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.EventFlowID == nil {
		return nil, errors.New("event flow id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.eventFlows, *input.EventFlowID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteEventFlowOutput{}, nil
}
