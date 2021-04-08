package ilert

import (
	"errors"
	"fmt"
)

// HeartbeatMethods defines uptime monitor regions
var HeartbeatMethods = struct {
	HEAD string
	GET  string
	POST string
}{
	HEAD: "HEAD",
	GET:  "GET",
	POST: "POST",
}

// PingHeartbeatInput represents the input of a PingHeartbeat operation.
type PingHeartbeatInput struct {
	_      struct{}
	Method *string
	APIKey *string
}

// PingHeartbeatOutput represents the output of a PingHeartbeat operation.
type PingHeartbeatOutput struct {
	_ struct{}
}

// PingHeartbeat gets list available iLert phone numbers. https://api.ilert.com/api-docs/#tag/Heartbeats/paths/~1heartbeats~1{key}/get
func (c *Client) PingHeartbeat(input *PingHeartbeatInput) (*PingHeartbeatOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.APIKey == nil {
		return nil, errors.New("APIKey is required")
	}
	if input.Method == nil {
		input.Method = String(HeartbeatMethods.HEAD)
	}

	resp, err := c.httpClient.R().Execute(*input.Method, fmt.Sprintf("%s/%s", apiRoutes.heartbeats, *input.APIKey))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return nil, apiErr
	}

	return &PingHeartbeatOutput{}, nil
}
