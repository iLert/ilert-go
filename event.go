package ilert

import (
	"encoding/json"
	"errors"
)

// EventComment represents a comment in an event
type EventComment struct {
	// Required. The creator of the comment
	Creator string `json:"creator"`
	// Required. The content of the comment
	Content string `json:"content"`
}

// EventLog represents a log entry in an event
type EventLog struct {
	// Required. The timestamp of the log entry
	Timestamp string `json:"timestamp"`
	// Required. The log level (e.g., INFO, WARN, ERROR)
	Level string `json:"level"`
	// Required. The log body/message
	Body string `json:"body"`
	// Optional. Labels for the log entry
	Labels map[string]string `json:"labels,omitempty"`
}

// Event represents the alert event https://api.ilert.com/api-docs/#tag/Events
type Event struct {
	// Required. The API key of the alert source.
	APIKey string `json:"apiKey"`

	// Required. Must be either ALERT, ACCEPT, or RESOLVE
	EventType string `json:"eventType"`

	// Required. The event summary. Will be used as the alert summary if a new alert is created.
	Summary string `json:"summary"`

	// Optional. The event details. Will be used as the alert details if a new alert is created.
	Details string `json:"details,omitempty"`

	// Optional. For ALERT events, the alert key can be used to deduplicate or group events. If an open alert with the key already exists, the event will be appended to the alert's event log. Otherwise a new alert will be created. For ACCEPT and RESOLVE events, the alert key is used to reference the open alert which is to be accepted or resolved by this event.
	AlertKey string `json:"alertKey,omitempty"`

	// Optional. For ALERT events, the alert key can be used to deduplicate or group events. If an open alert with the key already exists, the event will be appended to the alert's event log. Otherwise a new alert will be created. For ACCEPT and RESOLVE events, the alert key is used to reference the open alert which is to be accepted or resolved by this event.
	Priority string `json:"priority,omitempty"`

	// Optional. A list of images containing src, href and alt.
	Images []AlertImage `json:"images,omitempty"`

	// Optional. A list of links, containing href and text.
	Links []AlertLink `json:"links,omitempty"`

	// Optional. A list of labels.
	Labels map[string]string `json:"labels,omitempty"`

	// Optional. Additional custom details for the event.
	CustomDetails map[string]interface{} `json:"customDetails,omitempty"`

	// Optional. A list of comments for the event.
	Comments []EventComment `json:"comments,omitempty"`

	// Optional. A list of log entries for the event.
	Logs []EventLog `json:"logs,omitempty"`
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
	AlertKey     string `json:"alertKey"`
	AlertURL     string `json:"alertUrl"`
	ResponseCode string `json:"responseCode"`
}

// CreateEventInput represents the input of a CreateEvent operation.
type CreateEventInput struct {
	_ struct{}
	// alert event
	Event *Event
	// (optional) request url
	URL *string
}

// CreateEventOutput represents the output of a CreateEvent operation.
type CreateEventOutput struct {
	_             struct{}
	EventResponse *EventResponse
}

// CreateEvent creates an alert event. https://api.ilert.com/api-docs/#tag/Events/paths/~1events/post
func (c *Client) CreateEvent(input *CreateEventInput) (*CreateEventOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Event == nil {
		return nil, errors.New("input event is required")
	}
	url := apiRoutes.events
	if input.URL != nil && *input.URL != "" {
		url = *input.URL
	}
	resp, err := c.httpClient.R().SetBody(input.Event).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200, 202); apiErr != nil {
		return nil, apiErr
	}
	eventResponse := &EventResponse{}
	err = json.Unmarshal(resp.Body(), eventResponse)
	if err != nil {
		return nil, err
	}

	return &CreateEventOutput{EventResponse: eventResponse}, nil
}
