package ilert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
		r.StatusCode() >= http.StatusInternalServerError
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
			return &GenericAPIError{
				Status:  response.StatusCode(),
				Code:    "ERROR",
				Message: "An error occurred",
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
		if retryCondition(response, out) {
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
	connections         string
	connectors          string
	deploymentPipelines string
	escalationPolicies  string
	events              string
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
	connections:         "/api/v1/connections",
	connectors:          "/api/connectors",
	deploymentPipelines: "/api/deployment-pipelines",
	escalationPolicies:  "/api/escalation-policies",
	events:              "/api/events",
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
