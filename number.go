package ilert

import "encoding/json"

// Number definition https://api.ilert.com/api-docs/#tag/Numbers
type Number struct {
	CountryCode        string   `json:"countryCode"`
	PhoneNumber        string   `json:"phoneNumber"`
	SupportsInboundSMS bool     `json:"supportsInboundSms"`
	Types              []string `json:"types"`
}

// NumberTypes defines number types
var NumberTypes = struct {
	SMS   string
	Voice string
}{
	SMS:   "SMS",
	Voice: "VOICE",
}

// GetNumbersInput represents the input of a GetNumbers operation.
type GetNumbersInput struct {
	_ struct{}
}

// GetNumbersOutput represents the output of a GetNumbers operation.
type GetNumbersOutput struct {
	_       struct{}
	Numbers []*Number
}

// GetNumbers gets list available iLert phone numbers. https://api.ilert.com/api-docs/#tag/Numbers/paths/~1numbers/get
func (c *Client) GetNumbers(input *GetNumbersInput) (*GetNumbersOutput, error) {
	resp, err := c.httpClient.R().Get("/api/v1/numbers")
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	numbers := make([]*Number, 0)
	err = json.Unmarshal(resp.Body(), &numbers)
	if err != nil {
		return nil, err
	}

	output := &GetNumbersOutput{Numbers: numbers}

	return output, nil
}
