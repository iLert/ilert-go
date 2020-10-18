package ilert

import (
	"encoding/json"
	"errors"
)

// Event represents the incident event https://api.ilert.com/api-docs/#tag/Events
type Event struct {
	// Required. The API key of the alert source.
	APIKey string `json:"apiKey"`

	// Required. Must be either ALERT, ACCEPT, or RESOLVE
	EventType string `json:"eventType"`

	// Required. The event summary. Will be used as the incident summary if a new incident is created.
	Summary string `json:"summary"`

	// Optional. The event details. Will be used as the incident details if a new incident is created.
	Details string `json:"details"`

	// Optional. For ALERT events, the incident key can be used to deduplicate or group events. If an open incident with the key already exists, the event will be appended to the incident's event log. Otherwise a new incident will be created. For ACCEPT and RESOLVE events, the incident key is used to reference the open incident which is to be accepted or resolved by this event.
	IncidentKey string `json:"incidentKey"`

	// Optional. For ALERT events, the incident key can be used to deduplicate or group events. If an open incident with the key already exists, the event will be appended to the incident's event log. Otherwise a new incident will be created. For ACCEPT and RESOLVE events, the incident key is used to reference the open incident which is to be accepted or resolved by this event.
	Priority string `json:"priority"`

	// Optional. A list of images containing src, href and alt.
	Images []EventImage `json:"images"`

	// Optional. A list of links, containing href and text.
	Links []EventLink `json:"links"`

	// Optional. Additional custom details for the event.
	CustomDetails map[string]interface{} `json:"customDetails"`
}

// EventImage represents event image
type EventImage struct {
	Src  string `json:"src"`
	Href string `json:"href"`
	Alt  string `json:"alt"`
}

// EventLink represents event link
type EventLink struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

// EventTypes defines event types
var EventTypes = struct {
	Alert   string
	Accept  string
	Resolve string
}{
	Alert:   "ALERT",
	Accept:  "ACCEPT",
	Resolve: "RESOLVE",
}

// EventResponse describes event API response body
type EventResponse struct {
	IncidentKey  string `json:"incidentKey"`
	IncidentURL  string `json:"incidentUrl"`
	ResponseCode string `json:"responseCode"`
}

// CreateEventInput represents the input of a CreateEvent operation.
type CreateEventInput struct {
	_     struct{}
	Event *Event
}

// CreateEventOutput represents the output of a CreateEvent operation.
type CreateEventOutput struct {
	_             struct{}
	EventResponse *EventResponse
}

// CreateEvent creates an incident event. https://api.ilert.com/api-docs/#tag/Events/paths/~1events/post
func (c *Client) CreateEvent(input *CreateEventInput) (*CreateEventOutput, error) {
	if input == nil {
		return nil, errors.New("Input is required")
	}
	if input.Event == nil {
		return nil, errors.New("Input event is required")
	}
	output := &CreateEventOutput{}
	resp, err := c.httpClient.R().SetBody(input.Event).Post("/api/v1/events")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}
	eventResponse := &EventResponse{}
	err = json.Unmarshal(resp.Body(), eventResponse)
	if err != nil {
		return nil, err
	}
	output.EventResponse = eventResponse

	return output, nil
}
