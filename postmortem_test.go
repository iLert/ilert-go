package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequestIncidentPostmortem(t *testing.T) {
	var method, path string
	var query url.Values
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, query = r.Method, r.URL.Path, r.URL.Query()
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":5,"status":"REQUESTED","visibility":"PRIVATE"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).RequestIncidentPostmortem(&RequestIncidentPostmortemInput{
		IncidentID: Int64(8),
		Request: &PostmortemRequest{
			Alerts:     []PostmortemRequestReference{{ID: 31}},
			Channels:   []PostmortemRequestChannel{{ID: "C1", Type: IncidentChannelType.Slack, FromEpochSec: 1790000000}},
			RootCause:  "a bad deploy",
			Visibility: PostmortemVisibility.Private,
			Language:   "en",
		},
		Overwrite: Bool(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/incidents/8/postmortems" {
		t.Errorf("request = %s %s, want POST /api/incidents/8/postmortems", method, path)
	}
	if query.Get("overwrite") != "true" {
		t.Errorf("overwrite = %q, want true", query.Get("overwrite"))
	}
	if body["rootCause"] != "a bad deploy" || body["visibility"] != "PRIVATE" {
		t.Errorf("body = %v, want root cause and visibility", body)
	}
	if _, ok := body["deployments"]; ok {
		t.Errorf("deployments is present, want it absent when unset")
	}
	if result.Postmortem.ID != 5 || result.Postmortem.Status != PostmortemStatus.Requested {
		t.Errorf("postmortem = %+v, want postmortem 5 REQUESTED", result.Postmortem)
	}
}

func TestRequestIncidentPostmortemSendsNoUnsetOverwrite(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":5,"status":"REQUESTED"}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).RequestIncidentPostmortem(&RequestIncidentPostmortemInput{
		IncidentID: Int64(8),
		Request:    &PostmortemRequest{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := query["overwrite"]; ok {
		t.Errorf("overwrite is present, want it absent when unset")
	}
}

func TestLinkIncidentPostmortem(t *testing.T) {
	var method, path string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":6,"status":"LINKED","linkUrl":"https://wiki.example.com/pm","visibility":"PUBLIC"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).LinkIncidentPostmortem(&LinkIncidentPostmortemInput{
		IncidentID: Int64(8),
		Postmortem: &Postmortem{LinkURL: "https://wiki.example.com/pm", Visibility: PostmortemVisibility.Public},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/incidents/8/postmortem-links" {
		t.Errorf("request = %s %s, want POST /api/incidents/8/postmortem-links", method, path)
	}
	if body["linkUrl"] != "https://wiki.example.com/pm" {
		t.Errorf("body = %v, want the link url", body)
	}
	if result.Postmortem.Status != PostmortemStatus.Linked {
		t.Errorf("status = %q, want LINKED", result.Postmortem.Status)
	}
}

func TestLinkIncidentPostmortemRequiresLinkURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("unexpected request")
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).LinkIncidentPostmortem(&LinkIncidentPostmortemInput{
		IncidentID: Int64(8),
		Postmortem: &Postmortem{Visibility: PostmortemVisibility.Public},
	})
	if err == nil {
		t.Fatal("expected an error for a postmortem without a link url")
	}
}

func TestUpdateAndDeleteIncidentPostmortem(t *testing.T) {
	var requests []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":6,"status":"LINKED","linkUrl":"https://wiki.example.com/pm-v2","visibility":"PRIVATE"}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	updated, err := client.UpdateIncidentPostmortem(&UpdateIncidentPostmortemInput{
		IncidentID:   Int64(8),
		PostmortemID: Int64(6),
		Postmortem:   &Postmortem{Status: PostmortemStatus.Linked, LinkURL: "https://wiki.example.com/pm-v2", Visibility: PostmortemVisibility.Private},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Postmortem.LinkURL != "https://wiki.example.com/pm-v2" {
		t.Errorf("link url = %q, want the new link", updated.Postmortem.LinkURL)
	}
	if _, err := client.DeleteIncidentPostmortem(&DeleteIncidentPostmortemInput{IncidentID: Int64(8), PostmortemID: Int64(6)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"PUT /api/incidents/8/postmortems/6", "DELETE /api/incidents/8/postmortems/6"}
	if len(requests) != 2 || requests[0] != want[0] || requests[1] != want[1] {
		t.Errorf("requests = %v, want %v", requests, want)
	}
}

func TestGetIncidentPostmortems(t *testing.T) {
	var requests []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/incidents/8/postmortems" {
			_, _ = w.Write([]byte(`[{"id":6,"status":"LINKED","linkUrl":"https://wiki.example.com/pm","visibility":"PRIVATE","creator":{"id":9},"createdAt":"2026-10-01T10:00:00Z"}]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":6,"status":"CREATED","markdownUrl":"https://example.com/pm.md","visibility":"PUBLIC"}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	listed, err := client.GetIncidentPostmortems(&GetIncidentPostmortemsInput{IncidentID: Int64(8)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(listed.Postmortems) != 1 || listed.Postmortems[0].Creator == nil || listed.Postmortems[0].Creator.ID != 9 || listed.Postmortems[0].CreatedAt == "" {
		t.Errorf("postmortems = %+v, want one created by user 9", listed.Postmortems)
	}
	single, err := client.GetIncidentPostmortem(&GetIncidentPostmortemInput{IncidentID: Int64(8), PostmortemID: Int64(6)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if single.Postmortem.Status != PostmortemStatus.Created || single.Postmortem.MarkdownURL == "" {
		t.Errorf("postmortem = %+v, want a CREATED one with its document", single.Postmortem)
	}
	want := []string{"GET /api/incidents/8/postmortems", "GET /api/incidents/8/postmortems/6"}
	if len(requests) != 2 || requests[0] != want[0] || requests[1] != want[1] {
		t.Errorf("requests = %v, want %v", requests, want)
	}
}
