package ilert

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// The API rejects more than one state, so State is a single value and lands on the
// "states" parameter once. The id filters are repeatable.
func TestGetCallFlowSessionsQueryParams(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).GetCallFlowSessions(&GetCallFlowSessionsInput{
		StartIndex:        Int(0),
		MaxResults:        Int(25),
		State:             String(CallFlowSessionStatus.Ended),
		CallFlowIDs:       []*int64{Int64(11), Int64(12)},
		CallFlowNumberIDs: []*int64{Int64(21)},
		FromNumber:        String("+447700900123"),
		From:              String("2026-09-01T00:00:00Z"),
		Until:             String("2026-09-17T00:00:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := query["states"]; len(got) != 1 || got[0] != "ENDED" {
		t.Errorf("states = %v, want exactly one ENDED", got)
	}
	if got := query["call-flows"]; len(got) != 2 || got[0] != "11" || got[1] != "12" {
		t.Errorf("call-flows = %v, want [11 12]", got)
	}
	if got := query["call-flow-numbers"]; len(got) != 1 || got[0] != "21" {
		t.Errorf("call-flow-numbers = %v, want [21]", got)
	}
	for key, want := range map[string]string{
		"start-index": "0",
		"max-results": "25",
		"from-number": "+447700900123",
		"from":        "2026-09-01T00:00:00Z",
		"until":       "2026-09-17T00:00:00Z",
	} {
		if got := query.Get(key); got != want {
			t.Errorf("query %s = %q, want %q", key, got, want)
		}
	}
}

func TestGetCallFlowSessionsSendsNoUnsetParams(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	if _, err := newTestClient(t, srv.URL).GetCallFlowSessions(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"states", "call-flows", "call-flow-numbers", "from-number", "from", "until", "start-index", "max-results"} {
		if _, ok := query[key]; ok {
			t.Errorf("query %s is present, want it absent", key)
		}
	}
}

func TestGetCallFlowSessionDecodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":5,"callSid":"CA1","status":"ENDED","callFlowId":11,"numberId":21,"fromNumber":"+447700900123","fromCountry":"GB","connectedToUserId":42,"recordedMessageUrl":"https://example.com/voicemail.mp3"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetCallFlowSession(&GetCallFlowSessionInput{CallFlowSessionID: Int64(5)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := result.CallFlowSession
	if s.Status != CallFlowSessionStatus.Ended || s.CallFlowID != 11 || s.NumberID != 21 || s.ConnectedToUserID != 42 {
		t.Errorf("session = %+v, want the ended session on flow 11", s)
	}
	if s.RecordedMessageURL == "" {
		t.Error("recordedMessageUrl is empty, want the voicemail link")
	}
}

// The collected data has no fixed schema, its keys come from the nodes of the flow that ran,
// so it decodes into a map rather than a struct.
func TestGetCallFlowSessionDataDecodesFreeForm(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"pin":"1234","reason":"outage","attempts":2}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetCallFlowSessionData(&GetCallFlowSessionDataInput{CallFlowSessionID: Int64(5)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/call-flow-sessions/5/data" {
		t.Errorf("path = %q, want /api/call-flow-sessions/5/data", path)
	}
	if result.Data["pin"] != "1234" || result.Data["reason"] != "outage" {
		t.Errorf("data = %v, want the collected pin and reason", result.Data)
	}
}

func TestCallFlowSessionRequiredInputs(t *testing.T) {
	c := NewClient()
	if _, err := c.GetCallFlowSession(&GetCallFlowSessionInput{}); err == nil {
		t.Error("GetCallFlowSession without an id = nil error, want an error")
	}
	if _, err := c.GetCallFlowSessionData(&GetCallFlowSessionDataInput{}); err == nil {
		t.Error("GetCallFlowSessionData without an id = nil error, want an error")
	}
	if _, err := c.GetCallFlowSession(nil); err == nil {
		t.Error("GetCallFlowSession(nil) = nil error, want an error")
	}
}
