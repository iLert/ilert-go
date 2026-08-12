package ilert

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestTeamsMarshalEmptySlice verifies that an empty Teams slice survives into the
// payload as "teams":[]. The API only clears a resource's teams on an explicit
// empty array, so a field dropped by omitempty silently leaves them assigned.
// See iLert/engineering-tasks#2360.
func TestTeamsMarshalEmptySlice(t *testing.T) {
	empty := []TeamShort{}
	cases := []struct {
		name  string
		value any
	}{
		{"escalation policy", &EscalationPolicy{Name: "test", Teams: empty}},
		{"heartbeat monitor", &HeartbeatMonitor{Name: "test", Teams: empty}},
		{"incident template", &IncidentTemplate{Name: "test", Teams: empty}},
		{"schedule", &Schedule{Name: "test", Teams: empty}},
		{"support hour", &SupportHour{Name: "test", Teams: empty}},
		// already tagged json:"teams", covered here to keep the guarantee pinned
		{"alert source", &AlertSource{Name: "test", Teams: empty}},
		{"service", &Service{Name: "test", Teams: empty}},
		{"status page", &StatusPage{Name: "test", Teams: empty}},
		// pointer fields: a pointer to an empty slice is the "clear" signal
		{"alert action", &AlertAction{Name: "test", Teams: &empty}},
		{"call flow", &CallFlow{Name: "test", Teams: &empty}},
		{"event flow", &EventFlow{Name: "test", Teams: &empty}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(string(payload), `"teams":[]`) {
				t.Errorf("payload = %s, want it to contain \"teams\":[]", payload)
			}
		})
	}
}

// TestTeamsUnmarshalIntoOutputStructs verifies that the read path is untouched by
// the pointer-typed input fields: every API response decodes into the *Output
// structs, whose Teams stay plain slices.
func TestTeamsUnmarshalIntoOutputStructs(t *testing.T) {
	payload := []byte(`{"id":1,"name":"test","teams":[{"id":9501,"name":"TeamA"},{"id":9502,"name":"TeamB"}]}`)

	var callFlow CallFlowOutput
	if err := json.Unmarshal(payload, &callFlow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(callFlow.Teams) != 2 || callFlow.Teams[0].ID != 9501 {
		t.Errorf("CallFlowOutput.Teams = %v, want the two teams from the payload", callFlow.Teams)
	}

	var eventFlow EventFlowOutput
	if err := json.Unmarshal(payload, &eventFlow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(eventFlow.Teams) != 2 || eventFlow.Teams[1].ID != 9502 {
		t.Errorf("EventFlowOutput.Teams = %v, want the two teams from the payload", eventFlow.Teams)
	}
}

// TestTeamsUnmarshalIntoInputStructs verifies the pointer fields still decode, so
// a caller round-tripping a response through the input struct keeps working.
func TestTeamsUnmarshalIntoInputStructs(t *testing.T) {
	payload := []byte(`{"id":1,"name":"test","teams":[{"id":9501,"name":"TeamA"}]}`)

	var callFlow CallFlow
	if err := json.Unmarshal(payload, &callFlow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callFlow.Teams == nil || len(*callFlow.Teams) != 1 {
		t.Errorf("CallFlow.Teams = %v, want one team", callFlow.Teams)
	}

	var eventFlow EventFlow
	if err := json.Unmarshal(payload, &eventFlow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eventFlow.Teams == nil || len(*eventFlow.Teams) != 1 {
		t.Errorf("EventFlow.Teams = %v, want one team", eventFlow.Teams)
	}
}

// TestTeamsMarshalNil verifies that a caller who never touches Teams does not
// clear them. Most endpoints treat both an omitted field and a null as a no-op,
// but the call flow and event flow endpoints clear on null, so those two fields
// are pointers and have to disappear from the payload entirely.
func TestTeamsMarshalNil(t *testing.T) {
	cases := []struct {
		name     string
		value    any
		omitted  bool // field must not appear at all, null would clear
		expected string
	}{
		{"escalation policy", &EscalationPolicy{Name: "test"}, false, `"teams":null`},
		{"heartbeat monitor", &HeartbeatMonitor{Name: "test"}, false, `"teams":null`},
		{"incident template", &IncidentTemplate{Name: "test"}, false, `"teams":null`},
		{"schedule", &Schedule{Name: "test"}, false, `"teams":null`},
		{"support hour", &SupportHour{Name: "test"}, false, `"teams":null`},
		{"alert source", &AlertSource{Name: "test"}, false, `"teams":null`},
		{"service", &Service{Name: "test"}, false, `"teams":null`},
		{"status page", &StatusPage{Name: "test"}, false, `"teams":null`},
		{"alert action", &AlertAction{Name: "test"}, true, ""},
		{"call flow", &CallFlow{Name: "test"}, true, ""},
		{"event flow", &EventFlow{Name: "test"}, true, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.omitted {
				if strings.Contains(string(payload), `"teams"`) {
					t.Errorf("payload = %s, want no teams field", payload)
				}
				return
			}
			if !strings.Contains(string(payload), tc.expected) {
				t.Errorf("payload = %s, want it to contain %s", payload, tc.expected)
			}
		})
	}
}
