package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Status-Page definition https://api.ilert.com/api-docs/#!/Status-Pages
type StatusPage struct {
	ID                        int64       `json:"id"`
	Name                      string      `json:"name"`
	Domain                    string      `json:"domain"`
	Subdomain                 string      `json:"subdomain"`
	CustomCss                 string      `json:"customCss"`
	FaviconUrl                string      `json:"faviconUrl"`
	LogoUrl                   string      `json:"logoUrl"`
	Visibility                string      `json:"visibility"`
	HiddenFromSearch          bool        `json:"hiddenFromSearch"`
	ShowSubscribeAction       bool        `json:"showSubscribeAction"`
	ShowIncidentHistoryOption bool        `json:"showIncidentHistoryOption"`
	PageTitle                 string      `json:"pageTitle"`
	PageDescription           string      `json:"pageDescription"`
	LogoRedirectUrl           string      `json:"logoRedirectUrl"`
	Activated                 bool        `json:"activated"`
	Status                    string      `json:"status"`
	Teams                     []TeamShort `json:"teams"`
	Timezone                  string      `json:"timezone,omitempty"`
	Services                  []Service   `json:"services,omitempty"`
	Subscribed                bool        `json:"subscribed,omitempty"`
}

// StatusPageVisibility defines status page visibility
var StatusPageVisibility = struct {
	Public  string
	Private string
}{
	Public:  "PUBLIC",
	Private: "PRIVATE",
}

// StatusPageVisibilityAll defines status page visibility list
var StatusPageVisibilityAll = []string{
	StatusPageVisibility.Public,
	StatusPageVisibility.Private,
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
		return nil, errors.New("statuspage input is required")
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

// GetStatusPagesInput represents the input of a GetStatusPages operation.
type GetStatusPagesInput struct {
	_ struct{}
	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Default: 50, Maximum: 100
	MaxResults *int

	// describes optional properties that should be included in the response
	Include []*string
}

// GetStatusPagesOutput represents the output of a GetStatusPages operation.
type GetStatusPagesOutput struct {
	_           struct{}
	StatusPages []*StatusPage
}

// GetStatusPages lists statusPage sources. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages/get
func (c *Client) GetStatusPages(input *GetStatusPagesInput) (*GetStatusPagesOutput, error) {
	if input == nil {
		input = &GetStatusPagesInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.statusPages, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusPages := make([]*StatusPage, 0)
	err = json.Unmarshal(resp.Body(), &statusPages)
	if err != nil {
		return nil, err
	}

	return &GetStatusPagesOutput{StatusPages: statusPages}, nil
}

// GetStatusPageInput represents the input of a GetStatusPage operation.
type GetStatusPageInput struct {
	_            struct{}
	StatusPageID *int64

	// describes optional properties that should be included in the response
	Include []*string
}

// GetStatusPageOutput represents the output of a GetStatusPage operation.
type GetStatusPageOutput struct {
	_          struct{}
	StatusPage *StatusPage
}

// GetStatusPage gets a statuspage by ID. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}/get
func (c *Client) GetStatusPage(input *GetStatusPageInput) (*GetStatusPageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statuspage id is required")
	}

	q := url.Values{}

	for _, include := range input.Include {
		q.Add("include", *include)
	}

	var url = fmt.Sprintf("%s/%d?%s", apiRoutes.statusPages, *input.StatusPageID, q.Encode())

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusPage := &StatusPage{}
	err = json.Unmarshal(resp.Body(), statusPage)
	if err != nil {
		return nil, err
	}

	return &GetStatusPageOutput{StatusPage: statusPage}, nil
}

// GetStatusPageSubscribersInput represents the input of a GetStatusPageSubscribers operation.
type GetStatusPageSubscribersInput struct {
	_            struct{}
	StatusPageID *int64
}

// GetStatusPageSubscribersOutput represents the output of a GetStatusPageSubscribers operation.
type GetStatusPageSubscribersOutput struct {
	_           struct{}
	Subscribers []*Subscriber
}

// GetStatusPageSubscribers gets subscribers of a statusPage by ID. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1private-subscribers/get
func (c *Client) GetStatusPageSubscribers(input *GetStatusPageSubscribersInput) (*GetStatusPageSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statuspage id is required")
	}

	var url = fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.statusPages, *input.StatusPageID)

	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	subscribers := make([]*Subscriber, 0)
	err = json.Unmarshal(resp.Body(), &subscribers)
	if err != nil {
		return nil, err
	}

	return &GetStatusPageSubscribersOutput{Subscribers: subscribers}, nil
}

// SearchStatusPageInput represents the input of a SearchStatusPage operation.
type SearchStatusPageInput struct {
	_              struct{}
	StatusPageName *string
}

// SearchStatusPageOutput represents the output of a SearchStatusPage operation.
type SearchStatusPageOutput struct {
	_          struct{}
	StatusPage *StatusPage
}

// SearchStatusPage gets the statusPage with specified name.
func (c *Client) SearchStatusPage(input *SearchStatusPageInput) (*SearchStatusPageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageName == nil {
		return nil, errors.New("status page name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.statusPages, *input.StatusPageName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusPage := &StatusPage{}
	err = json.Unmarshal(resp.Body(), statusPage)
	if err != nil {
		return nil, err
	}

	return &SearchStatusPageOutput{StatusPage: statusPage}, nil
}

// UpdateStatusPageInput represents the input of a UpdateStatusPage operation.
type UpdateStatusPageInput struct {
	_            struct{}
	StatusPageID *int64
	StatusPage   *StatusPage
}

// UpdateStatusPageOutput represents the output of a UpdateStatusPage operation.
type UpdateStatusPageOutput struct {
	_          struct{}
	StatusPage *StatusPage
}

// UpdateStatusPage updates the specific statuspage. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}/put
func (c *Client) UpdateStatusPage(input *UpdateStatusPageInput) (*UpdateStatusPageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statuspage id is required")
	}
	if input.StatusPage == nil {
		return nil, errors.New("statuspage input is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.statusPages, *input.StatusPageID)

	resp, err := c.httpClient.R().SetBody(input.StatusPage).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	statusPage := &StatusPage{}
	err = json.Unmarshal(resp.Body(), statusPage)
	if err != nil {
		return nil, err
	}

	return &UpdateStatusPageOutput{StatusPage: statusPage}, nil
}

// AddStatusPageSubscriberInput represents the input of a AddStatusPageSubscriber operation.
type AddStatusPageSubscriberInput struct {
	_            struct{}
	StatusPageID *int64
	Subscriber   *Subscriber
}

// AddStatusPageSubscriberOutput represents the output of a AddStatusPageSubscriber operation.
type AddStatusPageSubscriberOutput struct {
	_ struct{}
}

// AddStatusPageSubscriber adds a new subscriber to an statuspage. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1private-subscribers/post
func (c *Client) AddStatusPageSubscriber(input *AddStatusPageSubscribersInput) (*AddStatusPageSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statuspage id is required")
	}
	if input.Subscribers == nil {
		return nil, errors.New("subscriber input is required")
	}

	url := fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.statusPages, *input.StatusPageID)

	resp, err := c.httpClient.R().SetBody(input.Subscribers).Post(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return nil, apiErr
	}

	subscriber := make([]*Subscriber, 0)
	err = json.Unmarshal(resp.Body(), &subscriber)
	if err != nil {
		return nil, err
	}

	return &AddStatusPageSubscribersOutput{}, nil
}

// AddStatusPageSubscribersInput represents the input of a AddStatusPageSubscribers operation.
type AddStatusPageSubscribersInput struct {
	_            struct{}
	StatusPageID *int64
	Subscribers  *[]Subscriber
}

// AddStatusPageSubscribersOutput represents the output of a AddStatusPageSubscribers operation.
type AddStatusPageSubscribersOutput struct {
	_ struct{}
}

// AddStatusPageSubscribers adds a new subscriber to an statuspage. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}~1private-subscribers/post
func (c *Client) AddStatusPageSubscribers(input *AddStatusPageSubscribersInput) (*AddStatusPageSubscribersOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statuspage id is required")
	}
	if input.Subscribers == nil {
		return nil, errors.New("subscriber input is required")
	}

	url := fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.statusPages, *input.StatusPageID)

	resp, err := c.httpClient.R().SetBody(input.Subscribers).Put(url)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return nil, apiErr
	}

	subscribers := make([]*Subscriber, 0)
	err = json.Unmarshal(resp.Body(), &subscribers)
	if err != nil {
		return nil, err
	}

	return &AddStatusPageSubscribersOutput{}, nil
}

// DeleteStatusPageInput represents the input of a DeleteStatusPage operation.
type DeleteStatusPageInput struct {
	_            struct{}
	StatusPageID *int64
}

// DeleteStatusPageOutput represents the output of a DeleteStatusPage operation.
type DeleteStatusPageOutput struct {
	_ struct{}
}

// DeleteStatusPage deletes the specified statuspage. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}/delete
func (c *Client) DeleteStatusPage(input *DeleteStatusPageInput) (*DeleteStatusPageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statusPage id is required")
	}

	url := fmt.Sprintf("%s/%d", apiRoutes.statusPages, *input.StatusPageID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteStatusPageOutput{}, nil
}

// DeleteStatusPageSubscriberInput represents the input of a DeleteStatusPageSubscriber operation.
type DeleteStatusPageSubscriberInput struct {
	_            struct{}
	StatusPageID *int64
	Subscriber   *Subscriber
}

// DeleteStatusPageSubscriberOutput represents the output of a DeleteStatusPageSubscriber operation.
type DeleteStatusPageSubscriberOutput struct {
	_ struct{}
}

// DeleteStatusPageSubscriber deletes a subscriber of the specified statuspage. https://api.ilert.com/api-docs/#tag/Status-Pages/paths/~1status-pages~1{id}/delete
func (c *Client) DeleteStatusSubscriberPage(input *DeleteStatusPageSubscriberInput) (*DeleteStatusPageSubscriberOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.StatusPageID == nil {
		return nil, errors.New("statusPage id is required")
	}
	if input.Subscriber == nil {
		return nil, errors.New("subscriber is required")
	}

	url := fmt.Sprintf("%s/%d/private-subscribers", apiRoutes.statusPages, *input.StatusPageID)

	resp, err := c.httpClient.R().Delete(url)

	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 202); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteStatusPageSubscriberOutput{}, nil
}
