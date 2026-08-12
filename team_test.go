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
