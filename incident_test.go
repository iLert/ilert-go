package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetIncidentsQueryParams(t *testing.T) {
	var path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":8,"number":42,"title":"checkout down","status":"DECLARED","severity":2}]`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetIncidents(&GetIncidentsInput{
		MaxResults: Int(25),
		Include:    []*string{String(IncidentInclude.AffectedServices), String(IncidentInclude.Responders)},
		States:     []*string{String(IncidentStatus.Declared), String(IncidentStatus.Investigating)},
		Services:   []*int64{Int64(7), Int64(9)},
		Severity:   Int(2),
		From:       String("2026-09-01T00:00:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/incidents" {
		t.Errorf("path = %q, want /api/incidents", path)
	}
	if got := query["states"]; len(got) != 2 || got[0] != "DECLARED" || got[1] != "INVESTIGATING" {
		t.Errorf("states = %v, want [DECLARED INVESTIGATING]", got)
	}
	if got := query["services"]; len(got) != 2 || got[0] != "7" || got[1] != "9" {
		t.Errorf("services = %v, want [7 9]", got)
	}
	if got := query["include"]; len(got) != 2 {
		t.Errorf("include = %v, want two includes", got)
	}
	if query.Get("severity") != "2" || query.Get("max-results") != "25" || query.Get("from") != "2026-09-01T00:00:00Z" {
		t.Errorf("query = %s, want severity, max-results and from", query.Encode())
	}
	for _, key := range []string{"start-index", "until"} {
		if _, ok := query[key]; ok {
			t.Errorf("query %s is present, want it absent when unset", key)
		}
	}
	if len(result.Incidents) != 1 || result.Incidents[0].Number != 42 || result.Incidents[0].Severity != 2 {
		t.Errorf("incidents = %+v, want incident #42 with severity 2", result.Incidents)
	}
}

// The API answers a created incident with 201 rather than the documented 200. A nil list of
// affected services goes out as null, which the API reads as "leave them alone".
func TestCreateIncident(t *testing.T) {
	var method, path string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":8,"number":42,"title":"checkout down","status":"DECLARED","declaredAt":"2026-09-30T10:00:00Z"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).CreateIncident(&CreateIncidentInput{
		Incident: &Incident{Title: "checkout down", Severity: 2},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/incidents" {
		t.Errorf("request = %s %s, want POST /api/incidents", method, path)
	}
	if body["title"] != "checkout down" || body["severity"] != float64(2) {
		t.Errorf("body = %v, want title and severity", body)
	}
	if v, ok := body["affectedServices"]; !ok || v != nil {
		t.Errorf("affectedServices = %v (present %v), want an explicit null", v, ok)
	}
	for _, key := range []string{"id", "number", "summary", "status", "declaredAt", "responders"} {
		if _, ok := body[key]; ok {
			t.Errorf("body %s is present, want it absent when unset", key)
		}
	}
	if result.Incident.ID != 8 || result.Incident.Number != 42 {
		t.Errorf("incident = %+v, want id 8 and number 42", result.Incident)
	}
}

// A non-nil empty list of affected services is how an update removes all of them.
func TestUpdateIncidentClearsAffectedServicesAndSendsIfMatch(t *testing.T) {
	var method, path, ifMatch, raw string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, ifMatch = r.Method, r.URL.Path, r.Header.Get("If-Match")
		b, _ := io.ReadAll(r.Body)
		raw = string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":8,"status":"RESOLVED","affectedServices":[]}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).UpdateIncident(&UpdateIncidentInput{
		IncidentID: Int64(8),
		Incident:   &Incident{Status: IncidentStatus.Resolved, AffectedServices: []AffectedServices{}},
		ETag:       String(`"abc-def"`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPut || path != "/api/incidents/8" {
		t.Errorf("request = %s %s, want PUT /api/incidents/8", method, path)
	}
	if ifMatch != `"abc-def"` {
		t.Errorf("If-Match = %q, want \"abc-def\"", ifMatch)
	}
	if !strings.Contains(raw, `"affectedServices":[]`) {
		t.Errorf("body = %s, want an empty affectedServices list", raw)
	}
}

func TestGetIncidentDecodes(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", `"abc-def"`)
		_, _ = w.Write([]byte(`{
			"id":8,"number":42,"title":"checkout down","status":"INVESTIGATING","severity":1,
			"declaredBy":{"id":9,"firstName":"Ada"},
			"channel":{"id":"C1","name":"inc-42","type":"SLACK","isPrivate":true},
			"conferenceBridge":{"type":"ZOOM","href":"https://zoom.example.com/j/1"},
			"affectedServices":[{"impact":"MAJOR_OUTAGE","service":{"id":7,"name":"checkout"}}],
			"responders":[{"user":{"id":9},"status":"JOINED","role":"INCIDENT_COMMANDER","pagedByUserId":3}]
		}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetIncident(&GetIncidentInput{IncidentID: Int64(8)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/incidents/8" {
		t.Errorf("path = %q, want /api/incidents/8", path)
	}
	if result.ETag == nil || *result.ETag != `"abc-def"` {
		t.Errorf("etag = %v, want \"abc-def\"", result.ETag)
	}
	incident := result.Incident
	if incident.DeclaredBy == nil || incident.DeclaredBy.ID != 9 {
		t.Errorf("declaredBy = %+v, want user 9", incident.DeclaredBy)
	}
	if incident.Channel == nil || incident.Channel.Type != IncidentChannelType.Slack || !incident.Channel.IsPrivate {
		t.Errorf("channel = %+v, want a private Slack channel", incident.Channel)
	}
	if incident.ConferenceBridge == nil || incident.ConferenceBridge.Type != IncidentConferenceBridgeType.Zoom {
		t.Errorf("conference bridge = %+v, want Zoom", incident.ConferenceBridge)
	}
	if len(incident.AffectedServices) != 1 || incident.AffectedServices[0].Service.ID != 7 {
		t.Errorf("affected services = %+v, want service 7", incident.AffectedServices)
	}
	if len(incident.Responders) != 1 || incident.Responders[0].Role != IncidentResponderRole.IncidentCommander || incident.Responders[0].PagedByUserID != 3 {
		t.Errorf("responders = %+v, want the incident commander paged by user 3", incident.Responders)
	}
}

func TestCreateIncidentStatusUpdate(t *testing.T) {
	var method, path string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":4,"summary":"checkout down","status":"IDENTIFIED"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).CreateIncidentStatusUpdate(&CreateIncidentStatusUpdateInput{
		IncidentID:   Int64(8),
		StatusUpdate: &StatusUpdate{Summary: "checkout down", Status: StatusUpdateStatus.Identified, SendNotification: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/incidents/8/status-updates" {
		t.Errorf("request = %s %s, want POST /api/incidents/8/status-updates", method, path)
	}
	if body["status"] != "IDENTIFIED" || body["sendNotification"] != true {
		t.Errorf("body = %v, want status and sendNotification", body)
	}
	if result.StatusUpdate.ID != 4 {
		t.Errorf("status update id = %d, want 4", result.StatusUpdate.ID)
	}
}

func TestAddIncidentSubscribersEmptyAccepted(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).AddIncidentSubscribers(&AddIncidentSubscribersInput{
		IncidentID:  Int64(8),
		Subscribers: &[]Subscriber{{ID: 9, Type: SubscriberType.User}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/incidents/8/private-subscribers" {
		t.Errorf("path = %q, want /api/incidents/8/private-subscribers", path)
	}
}

func TestGetIncidentAffected(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"statusPagesInfo":[{"id":21,"label":"public page"},{"id":22,"label":"internal page"}],"privateStatusPages":1,"publicStatusPages":1,"privateSubscribers":3,"publicSubscribers":4}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetIncidentAffected(&GetIncidentAffectedInput{
		Incident: &Incident{AffectedServices: []AffectedServices{{Impact: "DEGRADED", Service: Service{ID: 7}}}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/incidents/publish-info" {
		t.Errorf("request = %s %s, want POST /api/incidents/publish-info", method, path)
	}
	if result.Affected.PublicSubscribers != 4 {
		t.Errorf("affected = %+v, want 4 public subscribers", result.Affected)
	}
	// the API answers with one entry per status page, not the single object the spec documents
	if pages := result.Affected.StatusPagesInfo; len(pages) != 2 || pages[1].ID != 22 || pages[1].Label != "internal page" {
		t.Errorf("status pages = %+v, want both pages", pages)
	}
}

func TestGetIncidentLogEntries(t *testing.T) {
	var path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"l1","timestamp":"2026-09-30T10:00:00Z","type":"OPERATIONAL_INCIDENT","subType":"COMMENT","eventContext":{"message":"rolling back"},"diff":{"status":"RESOLVED"},"actorId":"u9","actorType":"USER"}]`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetIncidentLogEntries(&GetIncidentLogEntriesInput{
		IncidentID: Int64(8),
		StartIndex: Int(100),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/incidents/8/log-entries" {
		t.Errorf("path = %q, want /api/incidents/8/log-entries", path)
	}
	if query.Get("start-index") != "100" {
		t.Errorf("start-index = %q, want 100", query.Get("start-index"))
	}
	for _, key := range []string{"max-results", "from", "until", "sub-types"} {
		if _, ok := query[key]; ok {
			t.Errorf("query %s is present, want it absent when unset", key)
		}
	}
	if len(result.LogEntries) != 1 || result.LogEntries[0].Diff["status"] != "RESOLVED" || result.LogEntries[0].ActorID != "u9" {
		t.Errorf("log entries = %+v, want one with its diff and actor", result.LogEntries)
	}
}

func TestGetIncidentLogEntriesRequiresBothEnds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("unexpected request")
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).GetIncidentLogEntries(&GetIncidentLogEntriesInput{
		IncidentID: Int64(8),
		Until:      String("2026-09-30T00:00:00Z"),
	})
	if err == nil {
		t.Fatal("expected an error for a window without a start")
	}
}

func TestGetIncidentSubscribers(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":9,"name":"Ada","type":"USER"}]`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetIncidentSubscribers(&GetIncidentSubscribersInput{IncidentID: Int64(8)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/incidents/8/private-subscribers" {
		t.Errorf("path = %q, want /api/incidents/8/private-subscribers", path)
	}
	if len(result.Subscribers) != 1 || result.Subscribers[0].ID != 9 || result.Subscribers[0].Type != SubscriberType.User {
		t.Errorf("subscribers = %+v, want user 9", result.Subscribers)
	}
}
