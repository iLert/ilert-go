package ilert

import (
	"encoding/json"
	"errors"
)

// Status-Page definition https://api.ilert.com/api-docs/#!/Status-Pages
type StatusPage struct {
	ID                        int64               `json:"id"`
	Name                      string              `json:"name"`
	Domain                    string              `json:"domain"`
	Subdomain                 string              `json:"subdomain"`
	CustomCss                 string              `json:"customCss"`
	FaviconUrl                string              `json:"faviconUrl"`
	LogoUrl                   string              `json:"logoUrl"`
	Visibility                string              `json:"visibility"`
	HiddenFromSearch          string              `json:"hiddenFromSearch"`
	ShowSubscribeAction       string              `json:"showSubscribeAction"`
	ShowIncidentHistoryOption string              `json:"showIncidentHistoryOption"`
	PageTitle                 string              `json:"pageTitle"`
	PageDescription           string              `json:"pageDescription"`
	LogoRedirectUrl           string              `json:"logoRedirectUrl"`
	Activated                 string              `json:"activated"`
	Status                    string              `json:"status"`
	Teams                     []TeamShort         `json:"teams"`
	Services                  []ServiceUptimeOnly `json:"services"`
}

// StatusPage defines status-page visibility
var StatusPageVisibility = struct {
	Public  string
	Private string
}{
	Public:  "PUBLIC",
	Private: "PRIVATE",
}

// CreateStatusPageInput represents the input of a CreateStatusPage operation.
type CreateStatusPageInput struct {
	_          struct{}
	StatusPage *StatusPage
}

// CreateStatusPageOutput represents the output of a CreateStatusPage operation.
type CreateStatusPageOutput struct {
	_          struct{}
	StatusPage *StatusPage
}

// CreateStatusPage creates a new status-page. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages/post
func (c *Client) CreateStatusPage(input *CreateStatusPageInput) (*CreateStatusPageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPage == nil {
		return nil, errors.New("User input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.StatusPage).Post(apiRoutes.statusPages)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	statusPage := &StatusPage{}
	err = json.Unmarshal(resp.Body(), statusPage)
	if err != nil {
		return nil, err
	}

	return &CreateStatusPageOutput{StatusPage: statusPage}, nil
}
