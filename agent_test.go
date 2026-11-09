package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateAgentSendsBetaHeader(t *testing.T) {
	var method, path, beta string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, beta = r.Method, r.URL.Path, r.Header.Get("ilert-beta")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":7,"name":"triage","status":"ENABLED","mcpServers":[{"name":"github","active":true,"credentials":{"id":3}}],"mcpTools":[]}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).CreateAgent(&CreateAgentInput{
		Agent: &Agent{
			Name:       "triage",
			Status:     AgentStatus.Enabled,
			McpServers: []AgentMcpServer{{Name: "github", Active: true, Credentials: &AgentMcpServerCredentials{ID: 3}}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/beta/agents" {
		t.Errorf("request = %s %s, want POST /api/beta/agents", method, path)
	}
	if beta == "" {
		t.Errorf("ilert-beta header is missing, the API rejects early access requests without it")
	}
	servers, _ := body["mcpServers"].([]any)
	if len(servers) != 1 || servers[0].(map[string]any)["credentials"].(map[string]any)["id"] != float64(3) {
		t.Errorf("mcpServers = %v, want the github server with credential 3", body["mcpServers"])
	}
	if result.Agent.ID != 7 || result.Agent.McpServers[0].Credentials.ID != 3 {
		t.Errorf("agent = %+v, want agent 7 with the github server", result.Agent)
	}
}

// The API replaces an agent with what an update sends, so the fields a caller
// read and left alone have to go back as they were.
func TestUpdateAgentSendsBackWhatWasRead(t *testing.T) {
	const stored = `{"id":7,"name":"triage","avatarUrl":"https://example.com/a.png","status":"DISABLED",` +
		`"teams":[{"id":1,"name":"ops"}],"mcpServers":[],"mcpTools":[],` +
		`"additionalInstructions":"be brief","rcaPolicy":{"type":"CUSTOM","condition":"alert.priority == 'HIGH'","scoreThreshold":70},` +
		`"autonomyLevel":"FULLY","customConfig":{"model":"x"}}`
	var method, path string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(stored))
			return
		}
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = w.Write([]byte(stored))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	read, err := client.GetAgent(&GetAgentInput{AgentID: Int64(7)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	agent := read.Agent
	agent.Name = "triage-2"
	if _, err := client.UpdateAgent(&UpdateAgentInput{AgentID: Int64(7), Agent: agent}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if method != http.MethodPut || path != "/api/beta/agents/7" {
		t.Errorf("request = %s %s, want PUT /api/beta/agents/7", method, path)
	}
	want := map[string]any{
		"name":                   "triage-2",
		"avatarUrl":              "https://example.com/a.png",
		"status":                 "DISABLED",
		"additionalInstructions": "be brief",
		"autonomyLevel":          "FULLY",
		"customConfig":           map[string]any{"model": "x"},
		"rcaPolicy":              map[string]any{"type": "CUSTOM", "condition": "alert.priority == 'HIGH'", "scoreThreshold": float64(70)},
	}
	for key, value := range want {
		if !reflect.DeepEqual(body[key], value) {
			t.Errorf("%s = %#v, want %#v", key, body[key], value)
		}
	}
}

// The API stores the MCP servers and tools as free-form objects. Left as they were
// read, they go back byte for byte: a server saved without its active flag stays
// active, a credential reference with a null id stays one the API skips, and fields
// the types do not know survive.
func TestUpdateAgentKeepsUnchangedMcpListsAsRead(t *testing.T) {
	const servers = `[{"name":"ilert-workforce-github-mcp","credentials":{"id":null},"futureField":{"a":1}}]`
	const tools = `[{"server":"ilert-workforce-github-mcp","name":"get_pull_request","inputOverwrites":[{"name":"limit","type":"number","value":[12345678901234567890]}]}]`
	const stored = `{"id":7,"name":"triage","customConfig":{"big":12345678901234567890},"mcpServers":` + servers + `,"mcpTools":` + tools + `}`
	var sent []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut {
			sent, _ = io.ReadAll(r.Body)
		}
		_, _ = w.Write([]byte(stored))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	read, err := client.GetAgent(&GetAgentInput{AgentID: Int64(7)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	read.Agent.Name = "triage-2"
	if _, err := client.UpdateAgent(&UpdateAgentInput{AgentID: Int64(7), Agent: read.Agent}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(sent, &body); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body["mcpServers"]) != servers {
		t.Errorf("mcpServers = %s, want them as read: %s", body["mcpServers"], servers)
	}
	if string(body["mcpTools"]) != tools {
		t.Errorf("mcpTools = %s, want them as read: %s", body["mcpTools"], tools)
	}
	if string(body["customConfig"]) != `{"big":12345678901234567890}` {
		t.Errorf("customConfig = %s, want the number kept exactly", body["customConfig"])
	}
	if string(body["name"]) != `"triage-2"` {
		t.Errorf("name = %s, want the change sent", body["name"])
	}
}

// Changed lists are sent as the caller set them, a list left alone next to one that
// changed still goes back as read.
func TestUpdateAgentSendsChangedMcpLists(t *testing.T) {
	const tools = `[{"server":"s","name":"t","futureField":true}]`
	const stored = `{"id":7,"name":"triage","mcpServers":[{"name":"old"}],"mcpTools":` + tools + `}`
	var sent []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut {
			sent, _ = io.ReadAll(r.Body)
		}
		_, _ = w.Write([]byte(stored))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	read, err := client.GetAgent(&GetAgentInput{AgentID: Int64(7)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	read.Agent.McpServers = []AgentMcpServer{{Name: "new", Active: true}}
	if _, err := client.UpdateAgent(&UpdateAgentInput{AgentID: Int64(7), Agent: read.Agent}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(sent, &body); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body["mcpServers"]) != `[{"name":"new","active":true}]` {
		t.Errorf("mcpServers = %s, want the changed servers", body["mcpServers"])
	}
	if string(body["mcpTools"]) != tools {
		t.Errorf("mcpTools = %s, want them as read: %s", body["mcpTools"], tools)
	}
}

func TestAgentMcpToolInputValuesUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		json string
		want AgentMcpToolInputValues
	}{
		{"null", `null`, nil},
		{"empty list", `[]`, AgentMcpToolInputValues{}},
		{"single string", `"main"`, AgentMcpToolInputValues{"main"}},
		{"empty string", `""`, nil},
		{"mixed list", `["a",5,1.5,true]`, AgentMcpToolInputValues{"a", "5", "1.5", "true"}},
		{"large number", `[12345678901234567890]`, AgentMcpToolInputValues{"12345678901234567890"}},
		{"nested values", `[{"a":1},[2,3]]`, AgentMcpToolInputValues{`{"a":1}`, `[2,3]`}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var overwrite AgentMcpToolInputOverwrite
			if err := json.Unmarshal([]byte(`{"name":"repo","type":"string","value":`+tc.json+`}`), &overwrite); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(overwrite.Value, tc.want) {
				t.Errorf("value = %#v, want %#v", overwrite.Value, tc.want)
			}
		})
	}
}

func TestAgentMcpToolInputValuesMarshalNilAsNull(t *testing.T) {
	encoded, err := json.Marshal(AgentMcpToolInputOverwrite{Name: "repo", Type: "string"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != `{"name":"repo","type":"string","value":null}` {
		t.Errorf("encoded = %s, want the value sent as null", encoded)
	}
}

func TestSearchAgentEscapesName(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":7,"name":"on call triage"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).SearchAgent(&SearchAgentInput{AgentName: String("on call triage")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/beta/agents/name/on%20call%20triage" {
		t.Errorf("path = %s, want the name escaped", path)
	}
	if result.Agent.ID != 7 {
		t.Errorf("agent = %+v, want agent 7", result.Agent)
	}
}
