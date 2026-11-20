package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateActionSendsBetaHeader(t *testing.T) {
	var method, path, beta string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, beta = r.Method, r.URL.Path, r.Header.Get("ilert-beta")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":9,"name":"page ops","templateKey":"webhook_alert","entityKind":"ALERT",` +
			`"config":{"method":"POST"},"credentialId":3,"origin":"NATIVE","status":"NEVER_TRIGGERED"}`))
	}))
	defer srv.Close()

	credentialID := int64(3)
	result, err := newEarlyAccessTestClient(t, srv.URL).CreateAction(&CreateActionInput{
		Action: &Action{
			Name:          "page ops",
			TemplateKey:   "webhook_alert",
			EntityKind:    ActionEntityKind.Alert,
			TriggerEvents: []string{"alert.created"},
			Config:        json.RawMessage(`{"method":"POST"}`),
			CredentialID:  &credentialID,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/beta/actions" {
		t.Errorf("request = %s %s, want POST /api/beta/actions", method, path)
	}
	if beta == "" {
		t.Errorf("ilert-beta header is missing, the API rejects early access requests without it")
	}
	if body["templateKey"] != "webhook_alert" || body["entityKind"] != "ALERT" || body["credentialId"] != float64(3) {
		t.Errorf("body = %v, want the template, entity kind and credential", body)
	}
	if config, _ := body["config"].(map[string]any); config["method"] != "POST" {
		t.Errorf("config = %v, want it sent as a JSON object", body["config"])
	}
	// a create sends only what is set
	for _, field := range []string{"steps", "inbound", "fieldMapping", "templateVersion", "description", "enabled"} {
		if _, ok := body[field]; ok {
			t.Errorf("create sent %s although it is not set", field)
		}
	}
	if result.Action.ID != 9 || result.Action.Origin != ActionOrigin.Native || string(result.Action.Config) != `{"method":"POST"}` {
		t.Errorf("action = %+v, want action 9 with its config", result.Action)
	}
}

// The API merges an update onto the action and keeps what is left out, so
// UpdateAction sends the empty fields empty, to clear them, and leaves out the
// read-only ones and the transport.
func TestUpdateActionSendsTheWholeAction(t *testing.T) {
	var method, path string
	var body map[string]json.RawMessage
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":9,"name":"page ops","entityKind":"ALERT"}`))
	}))
	defer srv.Close()

	id := int64(9)
	_, err := newEarlyAccessTestClient(t, srv.URL).UpdateAction(&UpdateActionInput{
		ActionID: &id,
		Action: &Action{
			ID: 9, Name: "page ops", TemplateKey: "email_alert", EntityKind: ActionEntityKind.Alert,
			Config: json.RawMessage(`{"to":"ops@example.com","registered":"by ilert"}`),
			// as a read reports them
			Transport: "HTTP", Origin: ActionOrigin.Native, Status: ActionStatus.OK, IntegrationURL: "https://example.com/x",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPut || path != "/api/beta/actions/9" {
		t.Errorf("request = %s %s, want PUT /api/beta/actions/9", method, path)
	}
	want := map[string]string{
		"name":               `"page ops"`,
		"templateKey":        `"email_alert"`,
		"entityKind":         `"ALERT"`,
		"config":             `{"to":"ops@example.com","registered":"by ilert"}`,
		"description":        `""`,
		"templateVersion":    `null`,
		"triggerEvents":      `[]`,
		"triggerEventDelays": `[]`,
		"conditionICL":       `""`,
		"steps":              `[]`,
		"fieldMapping":       `[]`,
		"inbound":            `{}`,
		"credentialId":       `null`,
		"bidirectional":      `false`,
		"teams":              `[]`,
		"labels":             `{}`,
	}
	for field, value := range want {
		if got, ok := body[field]; !ok || string(got) != value {
			t.Errorf("%s = %s, want %s", field, got, value)
		}
	}
	for field := range body {
		if _, ok := want[field]; !ok {
			t.Errorf("update sent %s = %s, want it left out", field, body[field])
		}
	}
}

func TestUpdateActionKeepsWhatIsSet(t *testing.T) {
	var body map[string]json.RawMessage
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":9}`))
	}))
	defer srv.Close()

	id, credentialID, version, enabled := int64(9), int64(3), 2, false
	_, err := newEarlyAccessTestClient(t, srv.URL).UpdateAction(&UpdateActionInput{
		ActionID: &id,
		Action: &Action{
			Name: "page ops", Enabled: &enabled, TemplateVersion: &version, CredentialID: &credentialID,
			TriggerEvents: []string{"alert.created"}, Bidirectional: true, Labels: map[string]string{"env": "prod"},
			Teams: []TeamShort{{ID: 1}}, Steps: json.RawMessage(`[{"id":"post"}]`),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for field, value := range map[string]string{
		"enabled": `false`, "templateVersion": `2`, "credentialId": `3`, "triggerEvents": `["alert.created"]`,
		"bidirectional": `true`, "labels": `{"env":"prod"}`, "teams": `[{"id":1}]`, "steps": `[{"id":"post"}]`,
	} {
		if got := string(body[field]); got != value {
			t.Errorf("%s = %s, want %s", field, got, value)
		}
	}
}

func TestSearchActionEscapesName(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":9,"name":"page ops/eu"}`))
	}))
	defer srv.Close()

	result, err := newEarlyAccessTestClient(t, srv.URL).SearchAction(&SearchActionInput{ActionName: String("page ops/eu")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/beta/actions/name/page%20ops%2Feu" {
		t.Errorf("path = %s, want the name escaped", path)
	}
	if result.Action.ID != 9 {
		t.Errorf("action = %+v, want action 9", result.Action)
	}
}

func TestDeleteActionAndNotFound(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"status":404,"message":"Could not find action with id 9"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newEarlyAccessTestClient(t, srv.URL)
	if _, err := client.DeleteAction(&DeleteActionInput{ActionID: Int64(9)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodDelete || path != "/api/beta/actions/9" {
		t.Errorf("request = %s %s, want DELETE /api/beta/actions/9", method, path)
	}
	_, err := client.GetAction(&GetActionInput{ActionID: Int64(9)})
	if _, ok := err.(*NotFoundAPIError); !ok {
		t.Errorf("err = %v, want a *NotFoundAPIError", err)
	}
}
