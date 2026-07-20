package ilert

import (
	"encoding/json"
	"testing"
)

// TestConnectionOutputParamsUnmarshalAutotaskIDs verifies that the legacy
// connection params decode Autotask companyId/queueId into their int64 fields
// whether the API returns them as JSON strings (which it does) or numbers, and
// that sibling fields keep decoding through the custom UnmarshalJSON.
func TestConnectionOutputParamsUnmarshalAutotaskIDs(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		company int64
		queue   int64
	}{
		{"string values", `{"companyId":"12345","queueId":"8","ticketType":"t"}`, 12345, 8},
		{"number values", `{"companyId":12345,"queueId":8,"ticketType":"t"}`, 12345, 8},
		{"missing", `{"ticketType":"t"}`, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p ConnectionOutputParams
			if err := json.Unmarshal([]byte(tc.payload), &p); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.CompanyID != tc.company {
				t.Errorf("CompanyID = %d, want %d", p.CompanyID, tc.company)
			}
			if p.QueueID != tc.queue {
				t.Errorf("QueueID = %d, want %d", p.QueueID, tc.queue)
			}
			if p.TicketType != "t" {
				t.Errorf("TicketType = %q, want %q", p.TicketType, "t")
			}
		})
	}
}
