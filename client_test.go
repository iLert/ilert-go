package ilert

import (
	"bytes"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

// errorResponse spins up a server that replies once with the given status,
// content-type and body and returns the error produced by GetEscalationPolicy.
func errorResponse(t *testing.T, status int, contentType, body string, headers map[string]string) error {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	id := int64(1)
	_, err := newTestClient(t, srv.URL).GetEscalationPolicy(&GetEscalationPolicyInput{EscalationPolicyID: &id})
	return err
}

// TestErrorClassificationByStatus covers the status-based classification,
// including the non-JSON 404/400 path and the transient-status retry branches.
func TestErrorClassificationByStatus(t *testing.T) {
	jsonCT := "application/json"
	htmlCT := "text/html"
	tests := []struct {
		name        string
		status      int
		contentType string
		body        string
		wantType    interface{}
	}{
		{"json 500 -> retryable", 500, jsonCT, `{"status":500,"code":"X","message":"boom"}`, &RetryableAPIError{}},
		{"json 429 -> retryable", 429, jsonCT, `{"status":429,"code":"RATE","message":"slow"}`, &RetryableAPIError{}},
		{"json 409 -> retryable", 409, jsonCT, `{"status":409,"code":"CONFLICT","message":"again"}`, &RetryableAPIError{}},
		{"json 403 -> fail fast", 403, jsonCT, `{"status":403,"code":"KEY_ERROR","message":"bad token"}`, &GenericAPIError{}},
		{"json 404 -> not found", 404, jsonCT, `{"status":404,"code":"NF","message":"gone"}`, &NotFoundAPIError{}},
		{"json 400 -> bad request", 400, jsonCT, `{"status":400,"code":"BAD","message":"nope"}`, &BadRequestAPIError{}},
		{"non-json 404 -> not found", 404, htmlCT, "<html>not found</html>", &NotFoundAPIError{}},
		{"non-json 400 -> bad request", 400, htmlCT, "<html>bad</html>", &BadRequestAPIError{}},
		{"non-json 403 -> upstream block", 403, htmlCT, "blocked", &RetryableAPIError{}},
		{"non-json 503 -> upstream block", 503, htmlCT, "unavailable", &RetryableAPIError{}},
		{"non-json 500 -> retryable", 500, htmlCT, "oops", &RetryableAPIError{}},
		{"non-json 401 -> fail fast", 401, htmlCT, "denied", &GenericAPIError{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := errorResponse(t, tc.status, tc.contentType, tc.body, nil)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if fmt.Sprintf("%T", err) != fmt.Sprintf("%T", tc.wantType) {
				t.Fatalf("expected %T, got %T: %v", tc.wantType, err, err)
			}
		})
	}
}

// TestJSONBodySniffedDespiteContentType verifies a JSON body is recognised even
// when the Content-Type does not say so, so a real JSON 403 still fails fast.
func TestJSONBodySniffedDespiteContentType(t *testing.T) {
	err := errorResponse(t, 403, "text/plain", `{"status":403,"code":"KEY_ERROR","message":"bad token"}`, nil)
	if _, ok := err.(*RetryableAPIError); ok {
		t.Fatalf("JSON body must be sniffed and fail fast, got *RetryableAPIError: %v", err)
	}
	if !strings.Contains(err.Error(), "bad token") {
		t.Errorf("expected API message preserved, got: %s", err.Error())
	}
}

// TestRequestIDPrecedenceAndTruncation checks the request-id header precedence
// and the body-snippet length cap appear in the surfaced message.
func TestRequestIDPrecedenceAndTruncation(t *testing.T) {
	longBody := strings.Repeat("x", maxErrorBodySnippet+50)
	err := errorResponse(t, 500, "text/html", longBody, map[string]string{
		"cf-ray":           "ray-should-lose",
		"x-amzn-RequestId": "amzn-should-win",
	})
	msg := err.Error()
	if !strings.Contains(msg, `x-amzn-RequestId="amzn-should-win"`) {
		t.Errorf("expected highest-precedence request id, got: %s", msg)
	}
	if strings.Contains(msg, "ray-should-lose") {
		t.Errorf("lower-precedence request id should not appear, got: %s", msg)
	}
	if !strings.Contains(msg, "...(truncated)") {
		t.Errorf("expected truncated body marker, got: %s", msg)
	}
}

func TestResponseRequestIDPrecedence(t *testing.T) {
	// no request-id headers -> empty handled via getGenericAPIError already;
	// here assert helper precedence directly through describeResponse output.
	err := errorResponse(t, 500, "text/html", "x", map[string]string{"x-request-id": "rid"})
	if !strings.Contains(err.Error(), `x-request-id="rid"`) {
		t.Errorf("expected x-request-id surfaced, got: %s", err.Error())
	}
}

func TestDebugEnabled(t *testing.T) {
	s := func(v string) *string { return &v }
	cases := map[string]struct {
		in   *string
		want bool
	}{
		"nil":         {nil, false},
		"true":        {s("true"), true},
		"TRUE":        {s("TRUE"), true},
		"True":        {s("True"), true},
		"one":         {s("1"), true},
		"padded true": {s("  true  "), true},
		"false":       {s("false"), false},
		"zero":        {s("0"), false},
		"empty":       {s(""), false},
		"garbage":     {s("yes-please"), false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := debugEnabled(tc.in); got != tc.want {
				t.Errorf("debugEnabled(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// timeoutError is a net.Error reporting a timeout, for classification tests.
type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestClassifyTransportError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"nil", nil, ""},
		{"timeout", timeoutError{}, "network timeout"},
		{"deadline", context.DeadlineExceeded, "request deadline exceeded"},
		{"dns", &net.DNSError{Err: "no such host", IsNotFound: true}, "DNS resolution failed"},
		{"tls", x509.UnknownAuthorityError{}, "TLS certificate verification failed"},
		{"conn", &net.OpError{Op: "dial", Err: errors.New("connection refused")}, "network connection error (dial)"},
		{"generic", errors.New("boom"), "request error: boom"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyTransportError(tc.err)
			if tc.want == "" {
				if got != "" {
					t.Fatalf("expected empty, got %q", got)
				}
				return
			}
			if !strings.Contains(got, tc.want) {
				t.Errorf("classifyTransportError(%v) = %q, want substring %q", tc.err, got, tc.want)
			}
		})
	}
}

// captureLogger is a resty.Logger that records everything it is asked to log.
type captureLogger struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (l *captureLogger) write(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf.WriteString(fmt.Sprintf(format, v...))
}
func (l *captureLogger) Errorf(format string, v ...interface{}) { l.write(format, v...) }
func (l *captureLogger) Warnf(format string, v ...interface{})  { l.write(format, v...) }
func (l *captureLogger) Debugf(format string, v ...interface{}) { l.write(format, v...) }
func (l *captureLogger) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

// TestDebugMasksAuthorization verifies that enabling debug never prints the
// bearer token: the Authorization header must be redacted in the trace.
func TestDebugMasksAuthorization(t *testing.T) {
	const token = "super-secret-token-value"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"name":"x"}`))
	}))
	defer srv.Close()

	log := &captureLogger{}
	c := NewClient(WithAPIEndpoint(srv.URL), WithAPIToken(token), WithDebug(true))
	c.httpClient.SetLogger(log)

	id := int64(1)
	if _, err := c.GetEscalationPolicy(&GetEscalationPolicyInput{EscalationPolicyID: &id}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := log.String()
	if strings.Contains(out, token) {
		t.Fatalf("debug output leaked the bearer token")
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Errorf("expected Authorization to be redacted in debug output")
	}
}
