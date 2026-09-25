package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// CallFlowSession definition https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-sessions
type CallFlowSession struct {
	ID int64 `json:"id,omitempty"`

	// the id of the call at the telephony vendor
	CallSid string `json:"callSid,omitempty"`

	// the call flow that handled the call
	CallFlowID int64 `json:"callFlowId,omitempty"`

	// the call flow number the call came in on
	NumberID int64 `json:"numberId,omitempty"`

	// possible values: "DEFAULT", "ONGOING", "ENDED", "BLOCKED"
	Status string `json:"status,omitempty"`

	// the number that was called, in E.164 format
	ToNumber string `json:"toNumber,omitempty"`

	// the number the call was made from, in E.164 format
	FromNumber string `json:"fromNumber,omitempty"`

	FromCity    string `json:"fromCity,omitempty"`
	FromZip     string `json:"fromZip,omitempty"`
	FromState   string `json:"fromState,omitempty"`
	FromCountry string `json:"fromCountry,omitempty"`

	// the user the call was connected to, when the flow routed it to one
	ConnectedToUserID int64 `json:"connectedToUserId,omitempty"`

	// where a voicemail recorded during the call can be fetched
	RecordedMessageURL string `json:"recordedMessageUrl,omitempty"`

	// when the call itself ended, as opposed to the flow. Date time string in ISO format.
	CallEndedAt string `json:"callEndedAt,omitempty"`

	// when the flow finished running. Date time string in ISO format.
	FlowEndedAt string `json:"flowEndedAt,omitempty"`

	CreatedAt string `json:"createdAt,omitempty"` // Date time string in ISO format
	UpdatedAt string `json:"updatedAt,omitempty"` // Date time string in ISO format
}

// CallFlowSessionStatus defines the states a call flow session is returned in
var CallFlowSessionStatus = struct {
	Default string
	Ongoing string
	Ended   string
	Blocked string
}{
	Default: "DEFAULT",
	Ongoing: "ONGOING",
	Ended:   "ENDED",
	Blocked: "BLOCKED",
}

// CallFlowSessionStatusAll defines the call flow session status list
var CallFlowSessionStatusAll = []string{
	CallFlowSessionStatus.Default,
	CallFlowSessionStatus.Ongoing,
	CallFlowSessionStatus.Ended,
	CallFlowSessionStatus.Blocked,
}

// GetCallFlowSessionsInput represents the input of a GetCallFlowSessions operation.
type GetCallFlowSessionsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int

	// filters the result by session state. The API documents a maximum of one state and does
	// not reject more, it silently answers with an empty list, which is why this is a single
	// value rather than a slice.
	// possible values: "DEFAULT", "ONGOING", "ENDED", "BLOCKED"
	State *string

	// filters the result by the call flows that handled the call
	CallFlowIDs []*int64

	// filters the result by the call flow numbers the call came in on
	CallFlowNumberIDs []*int64

	// filters the result by the number the call was made from, in E.164 format
	FromNumber *string

	// only return calls that started at or after this point in time. Date time string in ISO format.
	From *string

	// only return calls that started before this point in time. Date time string in ISO format.
	Until *string
}

// GetCallFlowSessionsOutput represents the output of a GetCallFlowSessions operation.
type GetCallFlowSessionsOutput struct {
	_                struct{}
	CallFlowSessions []*CallFlowSession
}

// GetCallFlowSessions lists existing call flow sessions. https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-sessions
func (c *Client) GetCallFlowSessions(input *GetCallFlowSessionsInput) (*GetCallFlowSessionsOutput, error) {
	if input == nil {
		input = &GetCallFlowSessionsInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}
	if input.State != nil {
		q.Add("states", *input.State)
	}
	for _, callFlowID := range input.CallFlowIDs {
		q.Add("call-flows", strconv.FormatInt(*callFlowID, 10))
	}
	for _, numberID := range input.CallFlowNumberIDs {
		q.Add("call-flow-numbers", strconv.FormatInt(*numberID, 10))
	}
	if input.FromNumber != nil {
		q.Add("from-number", *input.FromNumber)
	}
	if input.From != nil {
		q.Add("from", *input.From)
	}
	if input.Until != nil {
		q.Add("until", *input.Until)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.callFlowSessions, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	sessions := make([]*CallFlowSession, 0)
	err = json.Unmarshal(resp.Body(), &sessions)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowSessionsOutput{CallFlowSessions: sessions}, nil
}

// GetCallFlowSessionInput represents the input of a GetCallFlowSession operation.
type GetCallFlowSessionInput struct {
	_                 struct{}
	CallFlowSessionID *int64
}

// GetCallFlowSessionOutput represents the output of a GetCallFlowSession operation.
type GetCallFlowSessionOutput struct {
	_               struct{}
	CallFlowSession *CallFlowSession
}

// GetCallFlowSession gets the call flow session with the specified id. https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-sessions
func (c *Client) GetCallFlowSession(input *GetCallFlowSessionInput) (*GetCallFlowSessionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowSessionID == nil {
		return nil, errors.New("call flow session id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.callFlowSessions, *input.CallFlowSessionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	session := &CallFlowSession{}
	err = json.Unmarshal(resp.Body(), session)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowSessionOutput{CallFlowSession: session}, nil
}

// GetCallFlowSessionDataInput represents the input of a GetCallFlowSessionData operation.
type GetCallFlowSessionDataInput struct {
	_                 struct{}
	CallFlowSessionID *int64
}

// GetCallFlowSessionDataOutput represents the output of a GetCallFlowSessionData operation.
type GetCallFlowSessionDataOutput struct {
	_ struct{}

	// the values collected during the session, keyed by name. The API documents it as a
	// free-form object rather than a fixed schema, since the keys come from the nodes of
	// the call flow that ran.
	Data map[string]any
}

// GetCallFlowSessionData gets the data that was collected during a call flow session, for example
// by its gather nodes. https://docs.ilert.com/developer-docs/rest-api/api-reference/call-flow-sessions
func (c *Client) GetCallFlowSessionData(input *GetCallFlowSessionDataInput) (*GetCallFlowSessionDataOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CallFlowSessionID == nil {
		return nil, errors.New("call flow session id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/data", apiRoutes.callFlowSessions, *input.CallFlowSessionID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	data := make(map[string]any)
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	return &GetCallFlowSessionDataOutput{Data: data}, nil
}
