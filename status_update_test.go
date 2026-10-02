package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Status updates moved off /api/incidents, which the API answers with operational incidents for
// tokens on contract version 2. The state filter is "states": the API ignores a "state" parameter.
func TestGetStatusUpdatesQueryParams(t *testing.T) {
	var path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":3,"summary":"checkout down","status":"MONITORING"}]`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetStatusUpdates(&GetStatusUpdatesInput{
		Include:  []*string{String(StatusUpdateInclude.Subscribed)},
		States:   []*string{String(StatusUpdateStatus.Investigating), String(StatusUpdateStatus.Monitoring)},
		Services: []*int64{Int64(7)},
		From:     String("2026-09-01T00:00:00Z"),
		Until:    String("2026-09-30T00:00:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/status-updates" {
		t.Errorf("path = %q, want /api/status-updates", path)
	}
	if got := query["states"]; len(got) != 2 || got[0] != "INVESTIGATING" || got[1] != "MONITORING" {
		t.Errorf("states = %v, want [INVESTIGATING MONITORING]", got)
	}
	if _, ok := query["state"]; ok {
		t.Errorf("state is present, the API only reads states")
	}
	for key, want := range map[string]string{
		"start-index": "0",
		"max-results": "10",
		"include":     "subscribed",
		"services":    "7",
		"from":        "2026-09-01T00:00:00Z",
		"until":       "2026-09-30T00:00:00Z",
	} {
		if got := query.Get(key); got != want {
			t.Errorf("query %s = %q, want %q", key, got, want)
		}
	}
	if len(result.StatusUpdates) != 1 || result.StatusUpdates[0].Status != StatusUpdateStatus.Monitoring {
		t.Errorf("status updates = %+v, want one MONITORING", result.StatusUpdates)
	}
}

// The API answers a created status update with 201 rather than the documented 200.
func TestCreateStatusUpdateAcceptsCreated(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":4,"summary":"checkout down","status":"INVESTIGATING"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).CreateStatusUpdate(&CreateStatusUpdateInput{
		StatusUpdate: &StatusUpdate{Summary: "checkout down", Status: StatusUpdateStatus.Investigating},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/status-updates" {
		t.Errorf("request = %s %s, want POST /api/status-updates", method, path)
	}
	if result.StatusUpdate.ID != 4 {
		t.Errorf("id = %d, want 4", result.StatusUpdate.ID)
	}
}

// The API answers the forecast with one entry per status page, and with an empty list when no
// status page would show the status update.
func TestGetStatusUpdateAffectedDecodesStatusPages(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"statusPagesInfo":[],"privateStatusPages":0,"publicStatusPages":0,"privateSubscribers":0,"publicSubscribers":0}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetStatusUpdateAffected(&GetStatusUpdateAffectedInput{StatusUpdate: &StatusUpdate{Summary: "checkout down"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/status-updates/publish-info" {
		t.Errorf("request = %s %s, want POST /api/status-updates/publish-info", method, path)
	}
	if result.Affected.StatusPagesInfo == nil || len(result.Affected.StatusPagesInfo) != 0 {
		t.Errorf("status pages = %#v, want an empty list", result.Affected.StatusPagesInfo)
	}
}

func TestGetStatusUpdateDecodesHistoryAndETag(t *testing.T) {
	var path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", `"abc-def"`)
		_, _ = w.Write([]byte(`{"id":4,"status":"RESOLVED","history":[{"id":"h1","content":"fixed","incidentStatus":"RESOLVED","creator":{"id":9,"firstName":"Ada"},"createdAt":"2026-09-30T10:00:00Z"}]}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetStatusUpdate(&GetStatusUpdateInput{
		StatusUpdateID: Int64(4),
		Include:        []*string{String(StatusUpdateInclude.History)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/status-updates/4" || query.Get("include") != "history" {
		t.Errorf("request = %s?%s, want /api/status-updates/4?include=history", path, query.Encode())
	}
	if result.ETag == nil || *result.ETag != `"abc-def"` {
		t.Errorf("etag = %v, want \"abc-def\"", result.ETag)
	}
	history := result.StatusUpdate.History
	if len(history) != 1 || history[0].Status != StatusUpdateStatus.Resolved || history[0].Creator == nil || history[0].Creator.ID != 9 {
		t.Errorf("history = %+v, want one RESOLVED entry by user 9", history)
	}
}

func TestUpdateStatusUpdateSendsIfMatch(t *testing.T) {
	var method, path, ifMatch string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, ifMatch = r.Method, r.URL.Path, r.Header.Get("If-Match")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":4,"status":"RESOLVED"}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).UpdateStatusUpdate(&UpdateStatusUpdateInput{
		StatusUpdateID: Int64(4),
		StatusUpdate:   &StatusUpdate{Status: StatusUpdateStatus.Resolved},
		ETag:           String(`"abc-def"`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPut || path != "/api/status-updates/4" {
		t.Errorf("request = %s %s, want PUT /api/status-updates/4", method, path)
	}
	if ifMatch != `"abc-def"` {
		t.Errorf("If-Match = %q, want \"abc-def\"", ifMatch)
	}
}

// The API answers adding subscribers with an empty 202, which must not be decoded.
func TestAddStatusUpdateSubscribersEmptyAccepted(t *testing.T) {
	var path string
	var body []Subscriber
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).AddStatusUpdateSubscribers(&AddStatusUpdateSubscribersInput{
		StatusUpdateID: Int64(4),
		Subscribers:    &[]Subscriber{{ID: 12, Type: SubscriberType.Team}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/status-updates/4/private-subscribers" {
		t.Errorf("path = %q, want /api/status-updates/4/private-subscribers", path)
	}
	if len(body) != 1 || body[0].ID != 12 || body[0].Type != "TEAM" {
		t.Errorf("body = %+v, want team 12", body)
	}
}

func TestGetStatusUpdateLogEntries(t *testing.T) {
	var path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"l1","type":"STATUS_UPDATE","subType":"CREATED","eventContext":{"summary":"checkout down"},"actorType":"USER"}]`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetStatusUpdateLogEntries(&GetStatusUpdateLogEntriesInput{
		StatusUpdateID: Int64(4),
		MaxResults:     Int(50),
		From:           String("2026-09-01T00:00:00Z"),
		Until:          String("2026-09-30T00:00:00Z"),
		SubTypes:       []*string{String("CREATED"), String("UPDATED")},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/status-updates/4/log-entries" {
		t.Errorf("path = %q, want /api/status-updates/4/log-entries", path)
	}
	if got := query["sub-types"]; len(got) != 2 || got[0] != "CREATED" || got[1] != "UPDATED" {
		t.Errorf("sub-types = %v, want [CREATED UPDATED]", got)
	}
	if query.Get("max-results") != "50" || query.Get("from") == "" || query.Get("until") == "" {
		t.Errorf("query = %s, want max-results, from and until", query.Encode())
	}
	if _, ok := query["start-index"]; ok {
		t.Errorf("start-index is present, want it absent when unset")
	}
	if len(result.LogEntries) != 1 || result.LogEntries[0].EventContext["summary"] != "checkout down" {
		t.Errorf("log entries = %+v, want one with its event context", result.LogEntries)
	}
}

// The API only accepts a time window with both ends, so a half-open one fails before the request.
func TestGetStatusUpdateLogEntriesRequiresBothEnds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("unexpected request")
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).GetStatusUpdateLogEntries(&GetStatusUpdateLogEntriesInput{
		StatusUpdateID: Int64(4),
		From:           String("2026-09-01T00:00:00Z"),
	})
	if err == nil {
		t.Fatal("expected an error for a window without an end")
	}
}

// The service include is still called "incidents" and returns status updates.
func TestServiceIncludeReturnsStatusUpdates(t *testing.T) {
	if ServiceInclude.StatusUpdates != "incidents" {
		t.Errorf("ServiceInclude.StatusUpdates = %q, want incidents", ServiceInclude.StatusUpdates)
	}
	service := &Service{}
	if err := json.Unmarshal([]byte(`{"id":1,"incidents":[{"id":4,"summary":"checkout down","status":"IDENTIFIED"}]}`), service); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(service.StatusUpdates) != 1 || service.StatusUpdates[0].ID != 4 {
		t.Errorf("status updates = %+v, want status update 4", service.StatusUpdates)
	}
}
