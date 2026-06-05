package ilert

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient builds a client pointed at a test server with retry tuned down
// so retry-driven tests stay fast.
func newTestClient(t *testing.T, url string) *Client {
	t.Helper()
	c := NewClient(WithAPIEndpoint(url))
	c.httpClient.SetRetryWaitTime(1 * time.Millisecond).SetRetryMaxWaitTime(2 * time.Millisecond)
	return c
}

// TestWAFBlockIsRetryable verifies that a non-JSON 403 (as produced by a WAF or
// proxy) is classified as a retryable upstream block - the bug behind the
// intermittent Terraform 403s - rather than the opaque "An error occurred".
func TestWAFBlockIsRetryable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Server", "awselb/2.0")
		w.Header().Set("x-amzn-RequestId", "abc-123")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("<html><body>403 Forbidden - Request blocked</body></html>"))
	}))
	defer srv.Close()

	id := int64(1)
	_, err := newTestClient(t, srv.URL).GetEscalationPolicy(&GetEscalationPolicyInput{EscalationPolicyID: &id})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*RetryableAPIError); !ok {
		t.Fatalf("expected *RetryableAPIError so the provider retries the block, got %T: %v", err, err)
	}
	for _, want := range []string{"WAF", "awselb/2.0", "abc-123", "403"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error message missing %q: %s", want, err.Error())
		}
	}
}

// TestRealAuthErrorIsNotRetried verifies that a genuine JSON 403 from the API
// (e.g. an invalid token) is NOT treated as a retryable block.
func TestRealAuthErrorIsNotRetried(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"status":403,"code":"NO_PERMISSION","message":"insufficient permissions"}`))
	}))
	defer srv.Close()

	id := int64(1)
	_, err := newTestClient(t, srv.URL).GetEscalationPolicy(&GetEscalationPolicyInput{EscalationPolicyID: &id})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*RetryableAPIError); ok {
		t.Fatalf("a real JSON 403 must not be retryable, got *RetryableAPIError: %v", err)
	}
	if !strings.Contains(err.Error(), "insufficient permissions") {
		t.Errorf("expected the API message to be preserved, got: %s", err.Error())
	}
}

// TestWAFBlockSelfHeals verifies the client transparently retries a transient
// block and succeeds when a later attempt gets through (different egress IP).
func TestWAFBlockSelfHeals(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("blocked"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":42,"name":"default"}`))
	}))
	defer srv.Close()

	id := int64(42)
	out, err := newTestClient(t, srv.URL).GetEscalationPolicy(&GetEscalationPolicyInput{EscalationPolicyID: &id})
	if err != nil {
		t.Fatalf("expected self-heal on retry, got error: %v", err)
	}
	if out == nil || out.EscalationPolicy == nil || out.EscalationPolicy.Name != "default" {
		t.Fatalf("unexpected output: %+v", out)
	}
	if got := atomic.LoadInt32(&calls); got < 2 {
		t.Fatalf("expected at least one retry, server saw %d call(s)", got)
	}
}

// TestNotFoundStillClassified guards that the JSON error classification for
// known statuses is unchanged.
func TestNotFoundStillClassified(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":404,"code":"NOT_FOUND","message":"missing"}`))
	}))
	defer srv.Close()

	id := int64(1)
	_, err := newTestClient(t, srv.URL).GetEscalationPolicy(&GetEscalationPolicyInput{EscalationPolicyID: &id})
	if _, ok := err.(*NotFoundAPIError); !ok {
		t.Fatalf("expected *NotFoundAPIError, got %T: %v", err, err)
	}
}
