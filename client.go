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
	return err != nil ||
		r.StatusCode() == http.StatusTooManyRequests ||
		r.StatusCode() >= http.StatusInternalServerError ||
		looksLikeUpstreamBlock(r)
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
		parts = append(parts, fmt.Sprintf("request-id=%s", id))
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
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Sprintf("network timeout - the ilert API did not respond within the client timeout (%dms): %v", apiTimeoutMs, err)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Sprintf("request deadline exceeded: %v", err)
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

	// ILERT_DEBUG=true|1 turns on full request/response tracing (method, URL,
	// headers and bodies) via resty's debug logger. This is the quickest way to
	// see what is really happening - e.g. that a 403 was returned by a WAF and
	// not by the ilert API - without code changes.
	if debug := getEnv("ILERT_DEBUG"); debug != nil && (*debug == "true" || *debug == "1") {
		c.httpClient.SetDebug(true)
	}

	// Always surface intermittent intermediary blocks and transport failures,
	// even without full debug enabled, so they are visible in normal logs.
	c.httpClient.OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
		if looksLikeUpstreamBlock(r) {
			log.Printf("[WARN] ilert-go: %s %s was answered by an intermediary, most likely a WAF/proxy block: %s",
				r.Request.Method, r.Request.URL, describeResponse(r))
		}
		return nil
	})
	c.httpClient.OnError(func(req *resty.Request, err error) {
		log.Printf("[WARN] ilert-go: %s %s failed: %s", req.Method, req.URL, classifyTransportError(err))
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

// WithRetry enables retry logic with exponential backoff for the following errors:
//
// - any network errors
//
// - 5xx errors: this indicates an error in iLert
//
// - 429 Too Many Requests: you have reached your rate limit
func WithRetry(retryCount int, retryWaitTime time.Duration, retryMaxWaitTime time.Duration) ClientOptions {
	return func(c *Client) {
		c.httpClient.
			SetRetryCount(retryCount).
			SetRetryWaitTime(retryWaitTime).
			SetRetryMaxWaitTime(retryMaxWaitTime).
			AddRetryCondition(retryCondition)
	}
}

// getGenericAPIError extract API response error
func getGenericAPIError(response *resty.Response, expectedStatusCode ...int) error {
	if !intSliceContains(expectedStatusCode, response.StatusCode()) {
		out := &GenericAPIError{}
		err := json.Unmarshal(response.Body(), out)
		if err != nil {
			// The body is not a JSON error envelope from the ilert API. This
			// almost always means the response was produced by an intermediary
			// (WAF, reverse proxy, load balancer) rather than the API itself -
			// e.g. an IP-based rate-block. Preserve every available clue
			// (status, server, request-id, body) so the caller can tell what
			// actually happened instead of seeing a bare "An error occurred".
			desc := describeResponse(response)
			if looksLikeUpstreamBlock(response) {
				return &RetryableAPIError{
					Status:  response.StatusCode(),
					Code:    "UPSTREAM_BLOCK",
					Message: fmt.Sprintf("the request did not reach the ilert API; it was answered by an intermediary, most likely a WAF/proxy rate-block on the client IP (re-running on a different egress IP usually succeeds). %s", desc),
				}
			}
			if retryCondition(response, nil) {
				return &RetryableAPIError{
					Status:  response.StatusCode(),
					Code:    "ERROR",
					Message: fmt.Sprintf("an error occurred and the response body was not valid JSON. %s", desc),
				}
			}
			return &GenericAPIError{
				Status:  response.StatusCode(),
				Code:    "ERROR",
				Message: fmt.Sprintf("an error occurred and the response body was not valid JSON. %s", desc),
			}
		}
		if out.Status == 0 {
			out.Status = response.StatusCode()
		}
		if out.Status == http.StatusNotFound {
			return &NotFoundAPIError{
				Status:  out.Status,
				Code:    out.Code,
				Message: out.Message,
			}
		}

		if out.Status == http.StatusBadRequest {
			return &BadRequestAPIError{
				Status:  out.Status,
				Code:    out.Code,
				Message: out.Message,
			}
		}
		// Only retry on genuinely transient conditions (429, 5xx, intermediary
		// block). A real client error such as a 401/403 from the API itself
		// carries a JSON envelope and must fail fast so the cause (e.g. bad
		// credentials) is reported immediately instead of after the full retry
		// timeout. Pass nil so the status - not the always-non-nil out value -
		// decides retryability.
		if retryCondition(response, nil) {
			return &RetryableAPIError{
				Status:  out.Status,
				Code:    out.Code,
				Message: out.Message,
			}
		}
		return out
	}

	return nil
}

// apiRoutes defines api routes
var apiRoutes = struct {
	alerts              string
	alertActions        string
	alertSources        string
	automationRules     string
	callFlows           string
	connections         string
	connectors          string
	deploymentPipelines string
	escalationPolicies  string
	events              string
	eventFlows          string
	heartbeats          string
	heartbeatMonitors   string
	incidents           string
	incidentTemplates   string
	metrics             string
	metricDataSources   string
	numbers             string
	schedules           string
	series              string
	services            string
	statusPages         string
	supportHours        string
	uptimeMonitors      string
	users               string
	teams               string
}{
	alerts:              "/api/alerts",
	alertActions:        "/api/alert-actions",
	alertSources:        "/api/alert-sources",
	automationRules:     "/api/automation-rules",
	callFlows:           "/api/call-flows",
	connections:         "/api/v1/connections",
	connectors:          "/api/connectors",
	deploymentPipelines: "/api/deployment-pipelines",
	escalationPolicies:  "/api/escalation-policies",
	events:              "/api/events",
	eventFlows:          "/api/event-flows",
	heartbeats:          "/api/heartbeats",
	heartbeatMonitors:   "/api/heartbeat-monitors",
	incidents:           "/api/incidents",
	incidentTemplates:   "/api/incident-templates",
	metrics:             "/api/metrics",
	metricDataSources:   "/api/metric-data-sources",
	numbers:             "/api/numbers",
	schedules:           "/api/schedules",
	series:              "/api/series",
	services:            "/api/services",
	statusPages:         "/api/status-pages",
	supportHours:        "/api/support-hours",
	uptimeMonitors:      "/api/uptime-monitors",
	users:               "/api/users",
	teams:               "/api/teams",
}

func getEnv(key string) *string {
	if v := os.Getenv(key); len(v) != 0 {
		return String(v)
	}

	return nil
}
