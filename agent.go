package ilert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// Agent definition, an early access resource behind /api/beta/agents. It may
// change or be removed while it is in early access, and its operations need a
// client opted in with WithEarlyAccess(EarlyAccess.Agents).
//
// The API replaces an agent with what an update sends, and clears most of the
// fields it leaves out, among them the avatar and the MCP servers and tools.
// Update an agent by sending it back as GetAgent returned it, with the changes
// applied.
type Agent struct {
	ID        int64       `json:"id,omitempty"`
	Name      string      `json:"name"`
	AvatarURL string      `json:"avatarUrl,omitempty"`
	Status    string      `json:"status,omitempty"`
	Teams     []TeamShort `json:"teams"`

	// McpServers are the MCP servers the agent can call tools on. A read of a
	// single agent always returns them, the list only on request.
	McpServers []AgentMcpServer `json:"mcpServers"`

	// McpTools are the tools of the servers the agent may call. A read of a
	// single agent always returns them, the list only on request.
	McpTools []AgentMcpTool `json:"mcpTools"`

	// AdditionalInstructions, RcaPolicy and AutonomyLevel are only taken from
	// accounts in the beta program. For every other account the API ignores them
	// and keeps what is stored.
	AdditionalInstructions string          `json:"additionalInstructions,omitempty"`
	RcaPolicy              *AgentRcaPolicy `json:"rcaPolicy,omitempty"`
	AutonomyLevel          string          `json:"autonomyLevel,omitempty"`

	// CustomConfig, RepositoryConfig, EnableCodingSessions and
	// EnablePullRequestReviews are not available to customer accounts, the API
	// rejects any other value than an empty one with a 400. CustomConfig from
	// an account outside the beta program is ignored instead, like the fields
	// above. They are here so an update can send back what a read returned.
	CustomConfig             json.RawMessage `json:"customConfig,omitempty"`
	RepositoryConfig         json.RawMessage `json:"repositoryConfig,omitempty"`
	EnableCodingSessions     bool            `json:"enableCodingSessions,omitempty"`
	EnablePullRequestReviews bool            `json:"enablePullRequestReviews,omitempty"`

	CreatedAt string `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt string `json:"updatedAt,omitempty"` // date time string in ISO 8601

	// read holds the MCP servers and tools as the API returned them, see MarshalJSON
	read agentMcpRead
}

// agentMcpRead is the MCP servers and tools of an agent as the API returned them:
// the raw JSON, and the same decoded a second time, which the fields of the agent
// are compared with to tell whether the caller changed them
type agentMcpRead struct {
	servers        json.RawMessage
	tools          json.RawMessage
	decodedServers []AgentMcpServer
	decodedTools   []AgentMcpTool
}

// UnmarshalJSON decodes an agent and keeps its MCP servers and tools as they were
// returned, for MarshalJSON.
func (a *Agent) UnmarshalJSON(data []byte) error {
	type plain Agent
	var agent plain
	if err := json.Unmarshal(data, &agent); err != nil {
		return err
	}
	var raw struct {
		McpServers json.RawMessage `json:"mcpServers"`
		McpTools   json.RawMessage `json:"mcpTools"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*a = Agent(agent)
	a.read = agentMcpRead{servers: raw.McpServers, tools: raw.McpTools}
	if raw.McpServers != nil {
		if err := json.Unmarshal(raw.McpServers, &a.read.decodedServers); err != nil {
			return err
		}
	}
	if raw.McpTools != nil {
		if err := json.Unmarshal(raw.McpTools, &a.read.decodedTools); err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON sends the MCP servers and tools back exactly as the API returned them
// while the caller left them as they were read. The API stores them as free-form
// objects, and the types above only know what the web app writes today: decoded
// and encoded again, a field they do not know would be dropped, and a server saved
// without its active flag would be sent as inactive.
func (a Agent) MarshalJSON() ([]byte, error) {
	type plain Agent
	data, err := json.Marshal(plain(a))
	if err != nil {
		return nil, err
	}

	keepServers := a.read.servers != nil && reflect.DeepEqual(a.McpServers, a.read.decodedServers)
	keepTools := a.read.tools != nil && reflect.DeepEqual(a.McpTools, a.read.decodedTools)
	if !keepServers && !keepTools {
		return data, nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	if keepServers {
		fields["mcpServers"] = a.read.servers
	}
	if keepTools {
		fields["mcpTools"] = a.read.tools
	}
	return json.Marshal(fields)
}

// AgentMcpServer is an MCP server attached to an agent
type AgentMcpServer struct {
	// Name of a built-in server, such as ilert-workforce-github-mcp rather than
	// its short name github, or of a custom server the agent defines
	Name   string `json:"name"`
	Active bool   `json:"active"`

	// DisplayName and Brand are only used by custom servers, the brand picks
	// the icon the web app shows for the server
	DisplayName string `json:"displayName,omitempty"`
	Brand       string `json:"brand,omitempty"`

	// Credentials point at the credential the server authenticates with. The
	// credential has to be in the MCP category and readable by the caller. A
	// custom server takes its URL from the credential, whose sub category is
	// custom or starts with custom-.
	Credentials *AgentMcpServerCredentials `json:"credentials,omitempty"`
}

// AgentMcpServerCredentials references the credential of an MCP server
type AgentMcpServerCredentials struct {
	ID int64 `json:"id"`
}

// AgentMcpTool is a tool of an MCP server attached to an agent
type AgentMcpTool struct {
	Server          string                       `json:"server"`
	Name            string                       `json:"name"`
	Active          bool                         `json:"active"`
	InputOverwrites []AgentMcpToolInputOverwrite `json:"inputOverwrites"`
}

// AgentMcpToolInputOverwrite fixes the value of an input of a tool, instead of
// leaving it to the agent
type AgentMcpToolInputOverwrite struct {
	Name  string                  `json:"name"`
	Type  string                  `json:"type"`
	Value AgentMcpToolInputValues `json:"value"`
}

// AgentMcpToolInputValues are the values an input of a tool is overwritten with.
// The agent uses a single value as it is and chooses from several. An input
// without values is left to the agent, as if it had no overwrite.
//
// Values are always written as strings, the agent converts them to the type of
// the input. A read also accepts a single value, and numbers and booleans, which
// it turns into strings the way the web app does. An object or a list among the
// values is turned into its JSON.
type AgentMcpToolInputValues []string

// UnmarshalJSON accepts null, a single value or a list of values
func (v *AgentMcpToolInputValues) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var raw any
	if err := decoder.Decode(&raw); err != nil {
		return err
	}

	switch value := raw.(type) {
	case nil:
		*v = nil
	case string:
		if value == "" {
			*v = nil
		} else {
			*v = AgentMcpToolInputValues{value}
		}
	case []any:
		values := make(AgentMcpToolInputValues, 0, len(value))
		for _, item := range value {
			s, err := agentMcpToolInputValueString(item)
			if err != nil {
				return err
			}
			values = append(values, s)
		}
		*v = values
	default:
		s, err := agentMcpToolInputValueString(value)
		if err != nil {
			return err
		}
		*v = AgentMcpToolInputValues{s}
	}

	return nil
}

func agentMcpToolInputValueString(value any) (string, error) {
	switch item := value.(type) {
	case nil:
		return "", nil
	case string:
		return item, nil
	case map[string]any, []any:
		encoded, err := json.Marshal(item)
		return string(encoded), err
	default:
		return fmt.Sprint(item), nil
	}
}

// AgentRcaPolicy defines when an agent starts a root cause investigation
type AgentRcaPolicy struct {
	Type           string `json:"type,omitempty"`
	Condition      string `json:"condition,omitempty"`
	ScoreThreshold *int   `json:"scoreThreshold,omitempty"`
}

// AgentStatus defines agent status
var AgentStatus = struct {
	Enabled  string
	Disabled string
}{
	Enabled:  "ENABLED",
	Disabled: "DISABLED",
}

// AgentStatusAll defines agent status list
var AgentStatusAll = []string{
	AgentStatus.Enabled,
	AgentStatus.Disabled,
}

// AgentMcpToolInputType defines the type of a tool input
var AgentMcpToolInputType = struct {
	String  string
	Number  string
	Boolean string
	Any     string
	Null    string
}{
	String:  "string",
	Number:  "number",
	Boolean: "boolean",
	Any:     "any",
	Null:    "null",
}

// AgentMcpToolInputTypeAll defines tool input type list
var AgentMcpToolInputTypeAll = []string{
	AgentMcpToolInputType.String,
	AgentMcpToolInputType.Number,
	AgentMcpToolInputType.Boolean,
	AgentMcpToolInputType.Any,
	AgentMcpToolInputType.Null,
}

// CreateAgentInput represents the input of a CreateAgent operation.
type CreateAgentInput struct {
	_     struct{}
	Agent *Agent
}

// CreateAgentOutput represents the output of a CreateAgent operation.
type CreateAgentOutput struct {
	_     struct{}
	Agent *Agent
}

// CreateAgent creates a new agent resource.
func (c *Client) CreateAgent(input *CreateAgentInput) (*CreateAgentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Agent == nil {
		return nil, errors.New("agent input is required")
	}

	req, err := c.betaRequest(EarlyAccess.Agents)
	if err != nil {
		return nil, err
	}
	resp, err := req.SetBody(input.Agent).Post(apiRoutes.agents)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	agent := &Agent{}
	err = json.Unmarshal(resp.Body(), agent)
	if err != nil {
		return nil, err
	}

	return &CreateAgentOutput{Agent: agent}, nil
}

// GetAgentInput represents the input of a GetAgent operation.
type GetAgentInput struct {
	_       struct{}
	AgentID *int64
}

// GetAgentOutput represents the output of a GetAgent operation.
type GetAgentOutput struct {
	_     struct{}
	Agent *Agent
}

// GetAgent gets the agent resource with specified id.
func (c *Client) GetAgent(input *GetAgentInput) (*GetAgentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AgentID == nil {
		return nil, errors.New("agent id is required")
	}

	req, err := c.betaRequest(EarlyAccess.Agents)
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(fmt.Sprintf("%s/%d", apiRoutes.agents, *input.AgentID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	agent := &Agent{}
	err = json.Unmarshal(resp.Body(), agent)
	if err != nil {
		return nil, err
	}

	return &GetAgentOutput{Agent: agent}, nil
}

// GetAgentsInput represents the input of a GetAgents operation.
type GetAgentsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetAgentsOutput represents the output of a GetAgents operation.
type GetAgentsOutput struct {
	_      struct{}
	Agents []*Agent
}

// GetAgents lists existing agent resources. The agents in the list carry no MCP
// servers and tools, GetAgent returns them.
func (c *Client) GetAgents(input *GetAgentsInput) (*GetAgentsOutput, error) {
	if input == nil {
		input = &GetAgentsInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	req, err := c.betaRequest(EarlyAccess.Agents)
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(fmt.Sprintf("%s?%s", apiRoutes.agents, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	agents := make([]*Agent, 0)
	err = json.Unmarshal(resp.Body(), &agents)
	if err != nil {
		return nil, err
	}

	return &GetAgentsOutput{Agents: agents}, nil
}

// SearchAgentInput represents the input of a SearchAgent operation.
type SearchAgentInput struct {
	_         struct{}
	AgentName *string
}

// SearchAgentOutput represents the output of a SearchAgent operation.
type SearchAgentOutput struct {
	_     struct{}
	Agent *Agent
}

// SearchAgent gets the agent resource with specified name.
func (c *Client) SearchAgent(input *SearchAgentInput) (*SearchAgentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AgentName == nil {
		return nil, errors.New("agent name is required")
	}

	req, err := c.betaRequest(EarlyAccess.Agents)
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(fmt.Sprintf("%s/name/%s", apiRoutes.agents, url.PathEscape(*input.AgentName)))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	agent := &Agent{}
	err = json.Unmarshal(resp.Body(), agent)
	if err != nil {
		return nil, err
	}

	return &SearchAgentOutput{Agent: agent}, nil
}

// UpdateAgentInput represents the input of a UpdateAgent operation.
type UpdateAgentInput struct {
	_       struct{}
	AgentID *int64
	Agent   *Agent
}

// UpdateAgentOutput represents the output of a UpdateAgent operation.
type UpdateAgentOutput struct {
	_     struct{}
	Agent *Agent
}

// UpdateAgent updates an existing agent resource. The API replaces the agent with
// what is sent, so send it as GetAgent returned it with the changes applied.
func (c *Client) UpdateAgent(input *UpdateAgentInput) (*UpdateAgentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Agent == nil {
		return nil, errors.New("agent input is required")
	}
	if input.AgentID == nil {
		return nil, errors.New("agent id is required")
	}

	req, err := c.betaRequest(EarlyAccess.Agents)
	if err != nil {
		return nil, err
	}
	resp, err := req.SetBody(input.Agent).Put(fmt.Sprintf("%s/%d", apiRoutes.agents, *input.AgentID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	agent := &Agent{}
	err = json.Unmarshal(resp.Body(), agent)
	if err != nil {
		return nil, err
	}

	return &UpdateAgentOutput{Agent: agent}, nil
}

// DeleteAgentInput represents the input of a DeleteAgent operation.
type DeleteAgentInput struct {
	_       struct{}
	AgentID *int64
}

// DeleteAgentOutput represents the output of a DeleteAgent operation.
type DeleteAgentOutput struct {
	_ struct{}
}

// DeleteAgent deletes the specified agent resource.
func (c *Client) DeleteAgent(input *DeleteAgentInput) (*DeleteAgentOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.AgentID == nil {
		return nil, errors.New("agent id is required")
	}

	req, err := c.betaRequest(EarlyAccess.Agents)
	if err != nil {
		return nil, err
	}
	resp, err := req.Delete(fmt.Sprintf("%s/%d", apiRoutes.agents, *input.AgentID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteAgentOutput{}, nil
}
