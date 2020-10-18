package ilert

import (
	"encoding/json"
	"fmt"
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

// GenericErrorResponse describes generic API error response
type GenericErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// NewClient creates an API client using an API token
func NewClient(options ...ClientOptions) *Client {
	c := Client{
		apiEndpoint: apiEndpoint,
	}

	c.httpClient = resty.New()
	c.httpClient.SetHostURL(apiEndpoint)
	c.httpClient.SetTimeout(apiTimeoutMs * time.Millisecond)
	c.httpClient.SetHeader("Accept", "application/json")
	c.httpClient.SetHeader("Content-Type", "application/json")
	c.httpClient.SetHeader("User-Agent", fmt.Sprintf("ilert-go/%s", Version))
	c.httpClient.SetHeader("Accept-Encoding", "gzip")

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
		c.httpClient.SetHostURL(endpoint)
	}
}

func catchGenericAPIError(response *resty.Response, expectedStatusCode int) error {
	if response.StatusCode() != expectedStatusCode {
		restErr := fmt.Errorf("Wrong status code %d", response.StatusCode())
		respBody := &GenericErrorResponse{}
		err := json.Unmarshal(response.Body(), respBody)
		if err == nil && respBody.Message != "" {
			restErr = fmt.Errorf("%s: %s", respBody.Code, respBody.Message)
		}
		return restErr
	}

	return nil
}
