package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCredentialSendsBetaHeader(t *testing.T) {
	var method, path, beta string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, beta = r.Method, r.URL.Path, r.Header.Get("ilert-beta")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":3,"name":"github","typeHint":"BEARER","hostUrl":"https://api.githubcopilot.com","category":"MCP","teams":[]}`))
	}))
	defer srv.Close()

	result, err := newEarlyAccessTestClient(t, srv.URL).CreateCredential(&CreateCredentialInput{
		Credential: &Credential{
			Name:        "github",
			TypeHint:    CredentialType.Bearer,
			SecretValue: "ghp_secret",
			HostURL:     "https://api.githubcopilot.com",
			Category:    "MCP",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPost || path != "/api/beta/credentials" {
		t.Errorf("request = %s %s, want POST /api/beta/credentials", method, path)
	}
	if beta == "" {
		t.Errorf("ilert-beta header is missing, the API rejects early access requests without it")
	}
	if body["secretValue"] != "ghp_secret" || body["typeHint"] != "BEARER" {
		t.Errorf("body = %v, want the secret and type sent", body)
	}
	if _, ok := body["username"]; ok {
		t.Errorf("username is present, want it absent when unset: the API rejects one on a BEARER credential")
	}
	if result.Credential.ID != 3 || result.Credential.SecretValue != "" {
		t.Errorf("credential = %+v, want credential 3 without a secret", result.Credential)
	}
}

func TestGetCredentialNotFound(t *testing.T) {
	var path, beta string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, beta = r.URL.Path, r.Header.Get("ilert-beta")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":404,"code":"NOT_FOUND","message":"Could not find credential with id 3"}`))
	}))
	defer srv.Close()

	_, err := newEarlyAccessTestClient(t, srv.URL).GetCredential(&GetCredentialInput{CredentialID: Int64(3)})
	if _, ok := err.(*NotFoundAPIError); !ok {
		t.Fatalf("expected *NotFoundAPIError, got %T: %v", err, err)
	}
	if path != "/api/beta/credentials/3" || beta == "" {
		t.Errorf("request = %s with ilert-beta %q, want /api/beta/credentials/3 with the header", path, beta)
	}
}

func TestDeleteCredential(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	if _, err := newEarlyAccessTestClient(t, srv.URL).DeleteCredential(&DeleteCredentialInput{CredentialID: Int64(3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodDelete || path != "/api/beta/credentials/3" {
		t.Errorf("request = %s %s, want DELETE /api/beta/credentials/3", method, path)
	}
}
