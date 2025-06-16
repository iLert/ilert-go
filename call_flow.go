package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// CallFlow definition https://api.ilert.com/api-docs/#tag/call-flows
type CallFlow struct {
	ID             int64           `json:"id"`
	Name           string          `json:"name"`
	Language       string          `json:"language"`
	AssignedNumber *CallFlowNumber `json:"assignedNumber,omitempty"`
	Teams          []TeamShort     `json:"teams,omitempty"`
	RootNode       *CallFlowNode   `json:"root"`
}

// CallFlow definition https://api.ilert.com/api-docs/#tag/call-flows
type CallFlowOutput struct {
	ID             int64               `json:"id"`
	Name           string              `json:"name"`
	Language       string              `json:"language"`
	AssignedNumber *CallFlowNumber     `json:"assignedNumber,omitempty"`
	Teams          []TeamShort         `json:"teams,omitempty"`
	RootNode       *CallFlowNodeOutput `json:"root"`
}

type CallFlowNode struct {
	ID       int64            `json:"id"`
	NodeType string           `json:"nodeType"`
	Metadata interface{}      `json:"metadata,omitempty"`
	Branches []CallFlowBranch `json:"branches"`
}

type CallFlowNodeOutput struct {
	ID       int64                 `json:"id"`
	NodeType string                `json:"nodeType"`
	Metadata *CallFlowNodeMetadata `json:"metadata,omitempty"`
	Branches []CallFlowBranch      `json:"branches"`
}

type CallFlowBranch struct {
	ID         int64         `json:"id"`
	BranchType string        `json:"branchType"`
	Condition  string        `json:"condition,omitempty"`
	Target     *CallFlowNode `json:"target"`
}

type CallFlowNumber struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`
}

type PhoneNumber struct {
	RegionCode string `json:"regionCode"`
	Number     string `json:"number"`
}

type CallFlowNodeMetadata struct {
	TextMessage    string                           `json:"textMessage,omitempty"`    // IVR_MENU or AUDIO_MESSAGE or VOICEMAIL or PIN_CODE
	CustomAudioUrl string                           `json:"customAudioUrl,omitempty"` // IVR_MENU or AUDIO_MESSAGE or VOICEMAIL or PIN_CODE
	AIVoiceModel   string                           `json:"aiVoiceModel,omitempty"`   // IVR_MENU or AUDIO_MESSAGE or VOICEMAIL or PIN_CODE, one of CallFlowNodeMetadataAIVoiceModel
	EnabledOptions []string                         `json:"enabledOptions,omitempty"` // IVR_MENU
	Language       string                           `json:"language,omitempty"`       // IVR_MENU or AUDIO_MESSAGE, one of CallFlowNodeMetadataLanguage
	VarKey         string                           `json:"varKey,omitempty"`         // PLAIN
	VarValue       string                           `json:"varValue,omitempty"`       // PLAIN
	Codes          []CallFlowNodeMetadataCode       `json:"codes,omitempty"`          // PIN_CODE
	SupportHoursId int64                            `json:"supportHoursId,omitempty"` // SUPPORT_HOURS
	HoldAudioUrl   string                           `json:"holdAudioUrl,omitempty"`   // ROUTE_CALL
	Targets        []CallFlowNodeMetadataCallTarget `json:"targets,omitempty"`        // ROUTE_CALL
	CallStyle      string                           `json:"callStyle,omitempty"`      // ROUTE_CALL, one of CallFlowNodeMetadataCallStyle
	AlertSourceId  int64                            `json:"alertSourceId,omitempty"`  // CREATE_ALERT
	Retries        int64                            `json:"retries,omitempty"`        // IVR_MENU or PIN_CODE or ROUTE_CALL
	CallTimeoutSec int64                            `json:"callTimeoutSec,omitempty"` // ROUTE_CALL
	Blacklist      []string                         `json:"blacklist,omitempty"`      // BLOCK_NUMBERS
	Intents        []CallFlowNodeMetadataIntent     `json:"intents,omitempty"`        // AGENTIC
	Gathers        []CallFlowNodeMetadataGather     `json:"gathers,omitempty"`        // AGENTIC
	Enrichment     *CallFlowNodeMetadataEnrichment  `json:"enrichment,omitempty"`     // AGENTIC
}

// CallFlowNodeMetadataAIVoiceModel defines the voice model used
var CallFlowNodeMetadataAIVoiceModel = struct {
	Emma     string
	Liam     string
	Oliver   string
	Andreas  string
	Sophie   string
	Isabelle string
	Gordon   string
	Bruce    string
	Alfred   string
	Ellen    string
	Barbara  string
}{
	Emma:     "emma",
	Liam:     "liam",
	Oliver:   "oliver",
	Andreas:  "andreas",
	Sophie:   "sophie",
	Isabelle: "isabelle",
	Gordon:   "gordon",
	Bruce:    "bruce",
	Alfred:   "alfred",
	Ellen:    "ellen",
	Barbara:  "barbara",
}

var CallFlowNodeMetadataAIVoiceModelAll = []string{
	CallFlowNodeMetadataAIVoiceModel.Emma,
	CallFlowNodeMetadataAIVoiceModel.Liam,
	CallFlowNodeMetadataAIVoiceModel.Oliver,
	CallFlowNodeMetadataAIVoiceModel.Andreas,
	CallFlowNodeMetadataAIVoiceModel.Sophie,
	CallFlowNodeMetadataAIVoiceModel.Isabelle,
	CallFlowNodeMetadataAIVoiceModel.Gordon,
	CallFlowNodeMetadataAIVoiceModel.Bruce,
	CallFlowNodeMetadataAIVoiceModel.Alfred,
	CallFlowNodeMetadataAIVoiceModel.Ellen,
	CallFlowNodeMetadataAIVoiceModel.Barbara,
}

// CallFlowNodeMetadataLanguage defines the language
var CallFlowNodeMetadataLanguage = struct {
	English string
	German  string
	French  string
	Spanish string
	Dutch   string
	Russian string
	Italian string
}{
	English: "en",
	German:  "de",
	French:  "fr",
	Spanish: "es",
	Dutch:   "nl",
	Russian: "ru",
	Italian: "it",
}

var CallFlowNodeMetadataLanguageAll = []string{
	CallFlowNodeMetadataLanguage.English,
	CallFlowNodeMetadataLanguage.German,
	CallFlowNodeMetadataLanguage.French,
	CallFlowNodeMetadataLanguage.Spanish,
	CallFlowNodeMetadataLanguage.Dutch,
	CallFlowNodeMetadataLanguage.Russian,
	CallFlowNodeMetadataLanguage.Italian,
}

type CallFlowNodeMetadataCode struct {
	Code  int64  `json:"code,omitempty"`
	Label string `json:"label"`
}

type CallFlowNodeMetadataCallTarget struct {
	Target string `json:"target"`
	Type   string `json:"type"` // one of CallFlowNodeMetadataCallTargetType
}

var CallFlowNodeMetadataCallTargetType = struct {
	User           string
	OnCallSchedule string
	Number         string
}{
	User:           "USER",
	OnCallSchedule: "ON_CALL_SCHEDULE",
	Number:         "NUMBER",
}

var CallFlowNodeMetadataCallTargetTypeAll = []string{
	CallFlowNodeMetadataCallTargetType.User,
	CallFlowNodeMetadataCallTargetType.OnCallSchedule,
	CallFlowNodeMetadataCallTargetType.Number,
}

var CallFlowNodeMetadataCallStyle = struct {
	Ordered  string
	Random   string
	Parallel string
}{
	Ordered:  "ORDERED",
	Random:   "RANDOM",
	Parallel: "PARALLEL",
}

var CallFlowNodeMetadataCallStyleAll = []string{
	CallFlowNodeMetadataCallStyle.Ordered,
	CallFlowNodeMetadataCallStyle.Random,
	CallFlowNodeMetadataCallStyle.Parallel,
}

type CallFlowNodeMetadataIntent struct {
	Type        string   `json:"type,omitempty"` // one of CallFlowNodeMetadataIntentType
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

var CallFlowNodeMetadataIntentType = struct {
	Incident         string
	SystemOutage     string
	SecurityBreach   string
	TechnicalSupport string
	Inquiry          string
}{
	Incident:         "INCIDENT",
	SystemOutage:     "SYSTEM_OUTAGE",
	SecurityBreach:   "SECURITY_BREACH",
	TechnicalSupport: "TECHNICAL_SUPPORT",
	Inquiry:          "INQUIRY",
}

var CallFlowNodeMetadataIntentTypeAll = []string{
	CallFlowNodeMetadataIntentType.Incident,
	CallFlowNodeMetadataIntentType.SystemOutage,
	CallFlowNodeMetadataIntentType.SecurityBreach,
	CallFlowNodeMetadataIntentType.TechnicalSupport,
	CallFlowNodeMetadataIntentType.Inquiry,
}

type CallFlowNodeMetadataGather struct {
	Type     string `json:"type,omitempty"` // one of CallFlowNodeMetadataGatherType
	Label    string `json:"label,omitempty"`
	VarType  string `json:"varType,omitempty"` // one of CallFlowNodeMetadataGatherVarType
	Required bool   `json:"required,omitempty"`
	Question string `json:"question,omitempty"`
}

var CallFlowNodeMetadataGatherType = struct {
	CallerName       string
	ContactNumber    string
	Email            string
	Incident         string
	AffectedServices string
}{
	CallerName:       "CALLER_NAME",
	ContactNumber:    "CONTACT_NUMBER",
	Email:            "EMAIL",
	Incident:         "INCIDENT",
	AffectedServices: "AFFECTED_SERVICES",
}

var CallFlowNodeMetadataGatherTypeAll = []string{
	CallFlowNodeMetadataGatherType.CallerName,
	CallFlowNodeMetadataGatherType.ContactNumber,
	CallFlowNodeMetadataGatherType.Email,
	CallFlowNodeMetadataGatherType.Incident,
	CallFlowNodeMetadataGatherType.AffectedServices,
}

var CallFlowNodeMetadataGatherVarType = struct {
	Number  string
	Date    string
	Boolean string
	String  string
}{
	Number:  "NUMBER",
	Date:    "DATE",
	Boolean: "BOOLEAN",
	String:  "STRING",
}

var CallFlowNodeMetadataGatherVarTypeAll = []string{
	CallFlowNodeMetadataGatherVarType.Number,
	CallFlowNodeMetadataGatherVarType.Date,
	CallFlowNodeMetadataGatherVarType.Boolean,
	CallFlowNodeMetadataGatherVarType.String,
}

type CallFlowNodeMetadataEnrichment struct {
	Enabled          bool                                   `json:"enabled"`
	InformationTypes []string                               `json:"informationTypes"` // one or more of CallFlowNodeMetadataEnrichmentInformationType
	Sources          []CallFlowNodeMetadataEnrichmentSource `json:"sources"`
}

var CallFlowNodeMetadataEnrichmentInformationType = struct {
	Incident      string
	Maintenance   string
	ServiceStatus string
}{
	Incident:      "INCIDENT",
	Maintenance:   "MAINTENANCE",
	ServiceStatus: "SERVICE_STATUS",
}

type CallFlowNodeMetadataEnrichmentSource struct {
	ID   string `json:"id"`
	Type string `json:"type"` // one of CallFlowNodeMetadataEnrichmentSourceType
}

var CallFlowNodeMetadataEnrichmentSourceType = struct {
	StatusPage string
	Service    string
}{
	StatusPage: "STATUS_PAGE",
	Service:    "SERVICE",
}

// CreateCallFlowInput represents the input of a CreateCallFlow operation.
type CreateCallFlowInput struct {
	_        struct{}
	CallFlow *CallFlow

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// CreateCallFlowOutput represents the output of a CreateCallFlow operation.
type CreateCallFlowOutput struct {
	_        struct{}
	CallFlow *CallFlowOutput
}

// CreateCallFlow creates a new call flow resource. https://api.ilert.com/api-docs/#tag/call-flows/post/call-flows
func (c *Client) CreateCallFlow(input *CreateCallFlowInput) (*CreateCallFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlow == nil {
		return nil, errors.New("call flow input is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.CallFlow).Post(fmt.Sprintf("%s?%s", apiRoutes.callFlows, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	callFlow := &CallFlowOutput{}
	err = json.Unmarshal(resp.Body(), callFlow)
	if err != nil {
		return nil, err
	}

	return &CreateCallFlowOutput{CallFlow: callFlow}, nil
}

// GetCallFlowInput represents the input of a GetCallFlow operation.
type GetCallFlowInput struct {
	_          struct{}
	CallFlowID *int64

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// GetCallFlowOutput represents the output of a GetCallFlow operation.
type GetCallFlowOutput struct {
	_        struct{}
	CallFlow *CallFlowOutput
}

// GetCallFlow gets the call flows resource with specified id. https://api.ilert.com/api-docs/#tag/call-flows/get/call-flows/{id}
func (c *Client) GetCallFlow(input *GetCallFlowInput) (*GetCallFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowID == nil {
		return nil, errors.New("call flow id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.callFlows, *input.CallFlowID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlow := &CallFlowOutput{}
	err = json.Unmarshal(resp.Body(), callFlow)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowOutput{CallFlow: callFlow}, nil
}

// GetCallFlowsInput represents the input of a GetCallFlows operation.
type GetCallFlowsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetCallFlowsOutput represents the output of a GetCallFlows operation.
type GetCallFlowsOutput struct {
	_         struct{}
	CallFlows []*CallFlowOutput
}

// GetCallFlows lists existing call flow resources. https://api.ilert.com/api-docs/#tag/call-flows/get/call-flows
func (c *Client) GetCallFlows(input *GetCallFlowsInput) (*GetCallFlowsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.callFlows, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlows := make([]*CallFlowOutput, 0)
	err = json.Unmarshal(resp.Body(), &callFlows)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowsOutput{CallFlows: callFlows}, nil
}

// SearchCallFlowInput represents the input of a SearchCallFlow operation.
type SearchCallFlowInput struct {
	_            struct{}
	CallFlowName *string
}

// SearchCallFlowOutput represents the output of a SearchCallFlow operation.
type SearchCallFlowOutput struct {
	_        struct{}
	CallFlow *CallFlowOutput
}

// SearchCallFlow gets the call flow resource with specified name.
func (c *Client) SearchCallFlow(input *SearchCallFlowInput) (*SearchCallFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowName == nil {
		return nil, errors.New("call flow name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.callFlows, *input.CallFlowName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlow := &CallFlowOutput{}
	err = json.Unmarshal(resp.Body(), callFlow)
	if err != nil {
		return nil, err
	}

	return &SearchCallFlowOutput{CallFlow: callFlow}, nil
}

// UpdateCallFlowInput represents the input of a UpdateCallFlow operation.
type UpdateCallFlowInput struct {
	_          struct{}
	CallFlowID *int64
	CallFlow   *CallFlow

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// UpdateCallFlowOutput represents the output of a UpdateCallFlow operation.
type UpdateCallFlowOutput struct {
	_        struct{}
	CallFlow *CallFlowOutput
}

// UpdateCallFlow updates an existing call flow resource. https://api.ilert.com/api-docs/#tag/call-flows/put/call-flows/{id}
func (c *Client) UpdateCallFlow(input *UpdateCallFlowInput) (*UpdateCallFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlow == nil {
		return nil, errors.New("call flow input is required")
	}
	if input.CallFlowID == nil {
		return nil, errors.New("call flow id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.CallFlow).Put(fmt.Sprintf("%s/%d?%s", apiRoutes.callFlows, *input.CallFlowID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	callFlow := &CallFlowOutput{}
	err = json.Unmarshal(resp.Body(), callFlow)
	if err != nil {
		return nil, err
	}

	return &UpdateCallFlowOutput{CallFlow: callFlow}, nil
}

// DeleteCallFlowInput represents the input of a DeleteCallFlow operation.
type DeleteCallFlowInput struct {
	_          struct{}
	CallFlowID *int64
}

// DeleteCallFlowOutput represents the output of a DeleteCallFlow operation.
type DeleteCallFlowOutput struct {
	_ struct{}
}

// DeleteCallFlow deletes the specified call flow resource. https://api.ilert.com/api-docs/#tag/call-flows/delete/call-flows/{id}
func (c *Client) DeleteCallFlow(input *DeleteCallFlowInput) (*DeleteCallFlowOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowID == nil {
		return nil, errors.New("call flow id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.callFlows, *input.CallFlowID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteCallFlowOutput{}, nil
}
