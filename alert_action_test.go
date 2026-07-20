package ilert

import (
	"encoding/json"
	"testing"
)

// TestAlertActionOutputParamsUnmarshalAutotaskIDs verifies that the Autotask
// companyId/queueId fields decode into the int64 struct fields whether the API
// returns them as JSON strings (which it does) or numbers, and that sibling
// fields keep decoding through the custom UnmarshalJSON. See
// iLert/engineering-tasks#2265.
func TestAlertActionOutputParamsUnmarshalAutotaskIDs(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		company int64
		queue   int64
	}{
		{"string values", `{"companyId":"12345","queueId":"8","noteType":"1"}`, 12345, 8},
		{"number values", `{"companyId":12345,"queueId":8,"noteType":"1"}`, 12345, 8},
		{"missing", `{"noteType":"1"}`, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p AlertActionOutputParams
			if err := json.Unmarshal([]byte(tc.payload), &p); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.CompanyID != tc.company {
				t.Errorf("CompanyID = %d, want %d", p.CompanyID, tc.company)
			}
			if p.QueueID != tc.queue {
				t.Errorf("QueueID = %d, want %d", p.QueueID, tc.queue)
			}
			// sibling fields must still decode through the custom UnmarshalJSON
			if p.NoteType != "1" {
				t.Errorf("NoteType = %q, want %q", p.NoteType, "1")
			}
		})
	}
}
