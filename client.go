package ilert

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	apiEndpoint  = "https://api.ilert.com"
	apiTimeoutMs = 30000
)

// Client wraps http client
type Client struct {
	apiEndpoint string
	httpClient  *resty.Client
}

// GenericAPIError describes generic API response error e.g. bad request
type GenericAPIError struct {
	error
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (aerr *GenericAPIError) Error() string {
	return fmt.Sprintf("Error occurred with status code: %d, error code: %s, message: %s", aerr.Status, aerr.Code, aerr.Message)
}

// RetryableAPIError describes retryable API response error e.g. too many requests
type RetryableAPIError struct {
	error
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (aerr *RetryableAPIError) Error() string {
	return fmt.Sprintf("Retryable error occurred with status code: %d, error code: %s, message: %s", aerr.Status, aerr.Code, aerr.Message)
}

// NotFoundAPIError describes not-found API response error e.g. resource deleted or never exists
type NotFoundAPIError struct {
	error
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (aerr *NotFoundAPIError) Error() string {
	return fmt.Sprintf("Not found: api respond with status code: %d, error code: %s, message: %s", aerr.Status, aerr.Code, aerr.Message)
}

// BadRequestAPIError describes not-found API response error e.g. resource deleted or never exists
type BadRequestAPIError struct {
	error
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (aerr *BadRequestAPIError) Error() string {
	return fmt.Sprintf("Bad request: api respond with status code: %d, error code: %s, message: %s", aerr.Status, aerr.Code, aerr.Message)
}

// GenericCountResponse describes generic resources count response
type GenericCountResponse struct {
	Count int `json:"count"`
}

func retryCondition(r *resty.Response, err error) bool {
	return err != nil || isTransientStatus(r.StatusCode()) || looksLikeUpstreamBlock(r)
}

// isTransientStatus reports whether an HTTP status is worth retrying: rate
// limits, server errors, and transient conflicts/locks that typically clear on
// a second attempt (concurrent edits or eventual consistency). Genuine client
// errors such as 401/403/404/400 are deliberately excluded so they fail fast.
func isTransientStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests, // 429
		http.StatusConflict, // 409 - concurrent edit / eventual consistency
		http.StatusLocked,   // 423
		http.StatusTooEarly: // 425
		return true
	}
	return code >= http.StatusInternalServerError
}

// maxErrorBodySnippet caps how much of a non-JSON error body we keep for
// diagnostics, so a full HTML block page does not flood the logs.
const maxErrorBodySnippet = 512

// responseIsJSON reports whether a response body is a JSON document, i.e. it
// most likely originated from the ilert API rather than an intermediary that
// answers with HTML or plain text (WAF, reverse proxy, load balancer).
func responseIsJSON(r *resty.Response) bool {
	if strings.Contains(strings.ToLower(r.Header().Get("Content-Type")), "json") {
		return true
	}
	body := bytes.TrimSpace(r.Body())
	return len(body) > 0 && (body[0] == '{' || body[0] == '[')
}

// looksLikeUpstreamBlock reports whether a response looks like it was produced
// by an intermediary transiently blocking the request (e.g. an IP-based WAF
// rate-block) rather than by the ilert API. A genuine API 403 carries a JSON
// error envelope; a WAF/proxy block returns HTML or plain text. Such blocks are
// transient and safe to retry, whereas a real authorization error is not.
func looksLikeUpstreamBlock(r *resty.Response) bool {
	switch r.StatusCode() {
	case http.StatusForbidden, http.StatusServiceUnavailable:
		return !responseIsJSON(r)
	default:
		return false
	}
}

// describeResponse renders the diagnostic context of a failed response so the
// caller can tell what actually went wrong (WAF vs API, which edge node,
// rate-limit hints) instead of just seeing a bare status code.
func describeResponse(r *resty.Response) string {
	parts := []string{fmt.Sprintf("status=%d", r.StatusCode())}
	if v := r.Header().Get("Content-Type"); v != "" {
		parts = append(parts, fmt.Sprintf("content-type=%q", v))
	}
	if v := r.Header().Get("Server"); v != "" {
		parts = append(parts, fmt.Sprintf("server=%q", v))
	}
	if id := responseRequestID(r); id != "" {
		// responseRequestID already returns a "<header>=<value>" pair.
		parts = append(parts, id)
	}
	if v := r.Header().Get("Retry-After"); v != "" {
		parts = append(parts, fmt.Sprintf("retry-after=%q", v))
	}
	if snippet := bodySnippet(r); snippet != "" {
		parts = append(parts, fmt.Sprintf("body=%q", snippet))
	}
	return strings.Join(parts, " ")
}

// responseRequestID extracts the most relevant correlation id from the common
// request-id headers set by ilert and the intermediaries in front of it, so the
// id can be handed to support to trace the exact request.
func responseRequestID(r *resty.Response) string {
	for _, h := range []string{
		"x-amzn-RequestId",
		"x-amz-cf-id",
		"cf-ray",
		"x-request-id",
		"x-amzn-trace-id",
		"apigw-requestid",
	} {
		if v := r.Header().Get(h); v != "" {
			return fmt.Sprintf("%s=%q", h, v)
		}
	}
	return ""
}

// bodySnippet returns a whitespace-collapsed, length-capped copy of the response
// body suitable for inclusion in an error message or log line.
func bodySnippet(r *resty.Response) string {
	body := strings.TrimSpace(string(r.Body()))
	if body == "" {
		return ""
	}
	body = strings.Join(strings.Fields(body), " ")
	if len(body) > maxErrorBodySnippet {
		body = body[:maxErrorBodySnippet] + "...(truncated)"
	}
	return body
}

// classifyTransportError turns a low-level transport error (no HTTP response was
// received) into a human-readable explanation: timeout, DNS, TLS or connection
// failure. Used for logging so the cause of a failed request is obvious.
func classifyTransportError(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Sprintf("request deadline exceeded - the ilert API did not respond within the client timeout (%dms): %v", apiTimeoutMs, err)
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Sprintf("network timeout - the ilert API did not respond within the client timeout (%dms): %v", apiTimeoutMs, err)
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return fmt.Sprintf("DNS resolution failed - check connectivity/proxy/DNS: %v", err)
	}
	var unknownAuthErr x509.UnknownAuthorityError
	var certInvalidErr x509.CertificateInvalidError
	if errors.As(err, &unknownAuthErr) || errors.As(err, &certInvalidErr) {
		return fmt.Sprintf("TLS certificate verification failed - a proxy or MITM may be intercepting traffic: %v", err)
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return fmt.Sprintf("network connection error (%s) - the ilert API could not be reached: %v", opErr.Op, err)
	}
	return fmt.Sprintf("request error: %v", err)
}

// debugEnabled parses an ILERT_DEBUG-style value, honouring the common boolean
// spellings (true/TRUE/1/...) via strconv.ParseBool. A nil or unparseable value
// is treated as disabled.
func debugEnabled(v *string) bool {
	if v == nil {
		return false
	}
	enabled, err := strconv.ParseBool(strings.TrimSpace(*v))
	return err == nil && enabled
}

// NewClient creates an API client using an API token
func NewClient(options ...ClientOptions) *Client {
	c := Client{
		apiEndpoint: apiEndpoint,
	}

	c.httpClient = resty.New()
	c.httpClient.SetBaseURL(apiEndpoint)
	c.httpClient.SetTimeout(apiTimeoutMs * time.Millisecond)
	c.httpClient.SetHeader("Accept", "application/json")
	c.httpClient.SetHeader("Content-Type", "application/json")
	c.httpClient.SetHeader("User-Agent", fmt.Sprintf("ilert-go/%s", Version))
	c.httpClient.SetHeader("Accept-Encoding", "gzip")
	c.httpClient.SetRetryCount(4).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(retryCondition)

	// Mask the Authorization header before resty's debug logger ever writes a
	// request, so enabling debug never prints the bearer token or basic-auth
	// credentials. Registered unconditionally; it only runs while debug logging
	// is active.
	c.httpClient.OnRequestLog(func(rl *resty.RequestLog) error {
		if rl.Header.Get("Authorization") != "" {
			rl.Header.Set("Authorization", "[REDACTED]")
		}
		return nil
	})

	// ILERT_DEBUG=true turns on full request/response tracing (method, URL,
	// headers and bodies) via resty's debug logger. This is the quickest way to
	// see what is really happening - e.g. that a 403 was returned by a WAF and
	// not by the ilert API - without code changes.
	if debugEnabled(getEnv("ILERT_DEBUG")) {
		c.httpClient.SetDebug(true)
	}

	// When debug is enabled, classify transport failures (timeout / DNS / TLS /
	// connection) into a readable line. Gated behind debug and emitted via the
	// callback so it does not produce unsilenceable output for normal use; the
	// returned error already carries the same detail once per call.
	c.httpClient.OnError(func(req *resty.Request, err error) {
		if c.httpClient.Debug {
			log.Printf("[ilert-go] %s %s failed: %s", req.Method, req.URL, classifyTransportError(err))
		}
	})

	endpoint := getEnv("ILERT_ENDPOINT")
	if endpoint != nil {
		c.httpClient.SetBaseURL(*endpoint)
	}

	apiToken := getEnv("ILERT_API_TOKEN")
	organizationID := getEnv("ILERT_ORGANIZATION")
	username := getEnv("ILERT_USERNAME")
	password := getEnv("ILERT_PASSWORD")

	if apiToken != nil {
		WithAPIToken(*apiToken)(&c)
	} else if organizationID != nil && username != nil && password != nil {
		WithBasicAuth(*organizationID, *username, *password)(&c)
	}

	for _, opt := range options {
		opt(&c)
	}

	return &c
}

// ClientOptions allows for options to be passed into the Client for customization
type ClientOptions func(*Client)

// WithBasicAuth adds an basic auth credentials to the client
func WithBasicAuth(organizationID string, username string, password string) ClientOptions {
	return func(c *Client) {
		c.httpClient.SetBasicAuth(fmt.Sprintf("%s@%s", username, organizationID), password)
	}
}

// WithAPIToken adds an api token to the client
func WithAPIToken(apiToken string) ClientOptions {
	return func(c *Client) {
		c.httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	}
}

// WithAPIEndpoint allows for a custom API endpoint to be passed into the client
func WithAPIEndpoint(endpoint string) ClientOptions {
	return func(c *Client) {
		c.apiEndpoint = endpoint
		c.httpClient.SetBaseURL(endpoint)
	}
}

// WithUserAgent replace user agent to the client
func WithUserAgent(agent string) ClientOptions {
	return func(c *Client) {
		c.httpClient.SetHeader("User-Agent", agent)
	}
}

// WithDebug enables (or disables) full request/response tracing - method, URL,
// headers and bodies - via resty's debug logger. This makes it possible to see
// what is really happening on the wire, e.g. that a 403 was returned by a WAF
// or proxy and not by the ilert API. It can also be toggled with ILERT_DEBUG.
//
// The Authorization header is masked in the debug output, but other potentially
// sensitive data (request and response bodies, cookies, custom headers) is
// printed verbatim. Only enable debug in trusted environments and avoid
// committing the resulting logs.
func WithDebug(debug bool) ClientOptions {
	return func(c *Client) {
		c.httpClient.SetDebug(debug)
	}
}

// WithProxy setting a Proxy URL and Port
func WithProxy(url string) ClientOptions {
	return func(c *Client) {
		c.httpClient.SetProxy(url)
	}
}

// WithRetry tunes the retry counts and backoff timing. The retry condition is
// already registered by NewClient and is not changed here. Requests are retried
// with exponential backoff on:
//
// - any network errors
//
// - 5xx errors: this indicates an error in iLert
//
// - 429 Too Many Requests: you have reached your rate limit
//
// - transient 409/423/425 responses (concurrent edit / eventual consistency)
//
// - non-JSON 403/503 responses from an intermediary (WAF / proxy / load
// balancer) transiently blocking the request
func WithRetry(retryCount int, retryWaitTime time.Duration, retryMaxWaitTime time.Duration) ClientOptions {
	return func(c *Client) {
		c.httpClient.
			SetRetryCount(retryCount).
			SetRetryWaitTime(retryWaitTime).
			SetRetryMaxWaitTime(retryMaxWaitTime)
	}
}

// getGenericAPIError extract API response error
func getGenericAPIError(response *resty.Response, expectedStatusCode ...int) error {
	if intSliceContains(expectedStatusCode, response.StatusCode()) {
		return nil
	}

	// Try to parse a JSON error envelope from the ilert API for its code and
	// message. A failure here means the body did not come from the API (most
	// likely an intermediary such as a WAF / proxy / load balancer).
	status := response.StatusCode()
	out := &GenericAPIError{}
	isJSON := json.Unmarshal(response.Body(), out) == nil
	if isJSON && out.Status != 0 {
		status = out.Status
	}

	code, message := out.Code, out.Message
	if !isJSON {
		// Preserve every available clue (status, server, request-id, body) so
		// the caller can tell what actually happened instead of seeing a bare
		// "An error occurred".
		code = "ERROR"
		message = fmt.Sprintf("the response body was not valid JSON, so it most likely did not come from the ilert API. %s", describeResponse(response))
	}

	// Classify by status so callers that branch on the error type (e.g.
	// detecting a deleted resource via *NotFoundAPIError) work whether or not
	// the body was JSON. Only genuinely transient conditions are retryable; a
	// real client error such as 401/403 fails fast so the cause is reported
	// immediately instead of after the full retry timeout.
	switch {
	case status == http.StatusNotFound:
		return &NotFoundAPIError{Status: status, Code: code, Message: message}
	case status == http.StatusBadRequest:
		return &BadRequestAPIError{Status: status, Code: code, Message: message}
	case looksLikeUpstreamBlock(response):
		return &RetryableAPIError{
			Status:  status,
			Code:    "UPSTREAM_BLOCK",
			Message: fmt.Sprintf("the request did not reach the ilert API; it was answered by an intermediary, most likely a WAF/proxy rate-block on the client IP (re-running on a different egress IP usually succeeds). %s", describeResponse(response)),
		}
	case isTransientStatus(status):
		return &RetryableAPIError{Status: status, Code: code, Message: message}
	default:
		return &GenericAPIError{Status: status, Code: code, Message: message}
	}
}

// apiRoutes defines api routes
var apiRoutes = struct {
	alerts                string
	alertActions          string
	alertSources          string
	automationRules       string
	callFlows             string
	connections           string
	connectors            string
	deploymentPipelines   string
	escalationPolicies    string
	events                string
	eventFlows            string
	eventFlowIntegrations string
	heartbeats            string
	heartbeatMonitors     string
	incidents             string
	incidentTemplates     string
	metrics               string
	metricDataSources     string
	numbers               string
	schedules             string
	series                string
	services              string
	statusPages           string
	supportHours          string
	telemetrySources      string
	uptimeMonitors        string
	users                 string
	teams                 string
}{
	alerts:                "/api/alerts",
	alertActions:          "/api/alert-actions",
	alertSources:          "/api/alert-sources",
	automationRules:       "/api/automation-rules",
	callFlows:             "/api/call-flows",
	connections:           "/api/v1/connections",
	connectors:            "/api/connectors",
	deploymentPipelines:   "/api/deployment-pipelines",
	escalationPolicies:    "/api/escalation-policies",
	events:                "/api/events",
	eventFlows:            "/api/event-flows",
	eventFlowIntegrations: "/api/event-flow-integrations",
	heartbeats:            "/api/heartbeats",
	heartbeatMonitors:     "/api/heartbeat-monitors",
	incidents:             "/api/incidents",
	incidentTemplates:     "/api/incident-templates",
	metrics:               "/api/metrics",
	metricDataSources:     "/api/metric-data-sources",
	numbers:               "/api/numbers",
	schedules:             "/api/schedules",
	series:                "/api/series",
	services:              "/api/services",
	statusPages:           "/api/status-pages",
	supportHours:          "/api/support-hours",
	telemetrySources:      "/api/telemetry-sources",
	uptimeMonitors:        "/api/uptime-monitors",
	users:                 "/api/users",
	teams:                 "/api/teams",
}

func getEnv(key string) *string {
	if v := os.Getenv(key); len(v) != 0 {
		return String(v)
	}

	return nil
}
