package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// DeploymentPipeline definition https://api.ilert.com/api-docs/#tag/deployment-pipelines
type DeploymentPipeline struct {
	ID              int64       `json:"id"`
	Name            string      `json:"name"`
	IntegrationType string      `json:"integrationType"`
	IntegrationKey  string      `json:"integrationKey,omitempty"`
	Teams           []TeamShort `json:"teams,omitempty"`
	CreatedAt       string      `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt       string      `json:"updatedAt,omitempty"` // date time string in ISO 8601
	IntegrationUrl  string      `json:"integrationUrl,omitempty"`
	Params          interface{} `json:"params"`
}

// DeploymentPipelineOutput definition https://api.ilert.com/api-docs/#tag/deployment-pipelines
type DeploymentPipelineOutput struct {
	ID              int64                           `json:"id"`
	Name            string                          `json:"name"`
	IntegrationType string                          `json:"integrationType"`
	IntegrationKey  string                          `json:"integrationKey,omitempty"`
	Teams           []TeamShort                     `json:"teams,omitempty"`
	CreatedAt       string                          `json:"createdAt,omitempty"` // date time string in ISO 8601
	UpdatedAt       string                          `json:"updatedAt,omitempty"` // date time string in ISO 8601
	IntegrationUrl  string                          `json:"integrationUrl,omitempty"`
	Params          *DeploymentPipelineOutputParams `json:"params"`
}

// DeploymentPipelineParams defines settings for a deployment pipeline
type DeploymentPipelineOutputParams struct {
	BranchFilters []string `json:"branchFilters,omitempty"` // used for GitHub
	EventFilters  []string `json:"eventFilters,omitempty"`  // used for GitHub
}

// DeploymentPipelineGitHubParams definition
type DeploymentPipelineGitHubParams struct {
	BranchFilters []string `json:"branchFilters,omitempty"`
	EventFilters  []string `json:"eventFilters,omitempty"`
}

// IntegrationType defines integration type
var DeploymentPipelineIntegrationType = struct {
	Api    string
	GitHub string
}{
	Api:    "API",
	GitHub: "GITHUB",
}

// IntegrationTypeAll defines integration type list
var DeploymentPipelineIntegrationTypeAll = []string{
	DeploymentPipelineIntegrationType.Api,
	DeploymentPipelineIntegrationType.GitHub,
}

// GitHubEventFilterType defines github event filter type
var GitHubEventFilterType = struct {
	PullRequest string
	Push        string
	Release     string
}{
	PullRequest: "pull_request",
	Push:        "push",
	Release:     "release",
}

// GitHubEventFilterTypeAll defines github event filter type list
var GitHubEventFilterTypeAll = []string{
	GitHubEventFilterType.PullRequest,
	GitHubEventFilterType.Push,
	GitHubEventFilterType.Release,
}

// CreateDeploymentPipelineInput represents the input of a CreateDeploymentPipeline operation.
type CreateDeploymentPipelineInput struct {
	_                  struct{}
	DeploymentPipeline *DeploymentPipeline

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// CreateDeploymentPipelineOutput represents the output of a CreateDeploymentPipeline operation.
type CreateDeploymentPipelineOutput struct {
	_                  struct{}
	DeploymentPipeline *DeploymentPipeline
}

// CreateDeploymentPipeline creates a new deployment pipeline resource. https://api.ilert.com/api-docs/#tag/deployment-pipelines/post/deployment-pipelines
func (c *Client) CreateDeploymentPipeline(input *CreateDeploymentPipelineInput) (*CreateDeploymentPipelineOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.DeploymentPipeline == nil {
		return nil, errors.New("deployment pipeline input is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.DeploymentPipeline).Post(fmt.Sprintf("%s?%s", apiRoutes.deploymentPipelines, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	deploymentPipeline := &DeploymentPipeline{}
	err = json.Unmarshal(resp.Body(), deploymentPipeline)
	if err != nil {
		return nil, err
	}

	return &CreateDeploymentPipelineOutput{DeploymentPipeline: deploymentPipeline}, nil
}

// GetDeploymentPipelineInput represents the input of a GetDeploymentPipeline operation.
type GetDeploymentPipelineInput struct {
	_                    struct{}
	DeploymentPipelineID *int64

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// GetDeploymentPipelineOutput represents the output of a GetDeploymentPipeline operation.
type GetDeploymentPipelineOutput struct {
	_                  struct{}
	DeploymentPipeline *DeploymentPipeline
}

// GetDeploymentPipeline gets the deployment pipelines resource with specified id. https://api.ilert.com/api-docs/#tag/deployment-pipelines/get/deployment-pipelines/{id}
func (c *Client) GetDeploymentPipeline(input *GetDeploymentPipelineInput) (*GetDeploymentPipelineOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.DeploymentPipelineID == nil {
		return nil, errors.New("deployment pipeline id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d?%s", apiRoutes.deploymentPipelines, *input.DeploymentPipelineID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	deploymentPipeline := &DeploymentPipeline{}
	err = json.Unmarshal(resp.Body(), deploymentPipeline)
	if err != nil {
		return nil, err
	}

	return &GetDeploymentPipelineOutput{DeploymentPipeline: deploymentPipeline}, nil
}

// GetDeploymentPipelinesInput represents the input of a GetDeploymentPipelines operation.
type GetDeploymentPipelinesInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetDeploymentPipelinesOutput represents the output of a GetDeploymentPipelines operation.
type GetDeploymentPipelinesOutput struct {
	_                   struct{}
	DeploymentPipelines []*DeploymentPipeline
}

// GetDeploymentPipelines lists existing deployment pipeline resources. https://api.ilert.com/api-docs/#tag/deployment-pipelines/get/deployment-pipelines
func (c *Client) GetDeploymentPipelines(input *GetDeploymentPipelinesInput) (*GetDeploymentPipelinesOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.deploymentPipelines, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	deploymentPipelines := make([]*DeploymentPipeline, 0)
	err = json.Unmarshal(resp.Body(), &deploymentPipelines)
	if err != nil {
		return nil, err
	}

	return &GetDeploymentPipelinesOutput{DeploymentPipelines: deploymentPipelines}, nil
}

// SearchDeploymentPipelineInput represents the input of a SearchDeploymentPipeline operation.
type SearchDeploymentPipelineInput struct {
	_                      struct{}
	DeploymentPipelineName *string
}

// SearchDeploymentPipelineOutput represents the output of a SearchDeploymentPipeline operation.
type SearchDeploymentPipelineOutput struct {
	_                  struct{}
	DeploymentPipeline *DeploymentPipeline
}

// SearchDeploymentPipeline gets the deployment pipeline resource with specified name.
func (c *Client) SearchDeploymentPipeline(input *SearchDeploymentPipelineInput) (*SearchDeploymentPipelineOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.DeploymentPipelineName == nil {
		return nil, errors.New("deployment pipeline name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.deploymentPipelines, *input.DeploymentPipelineName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	deploymentPipeline := &DeploymentPipeline{}
	err = json.Unmarshal(resp.Body(), deploymentPipeline)
	if err != nil {
		return nil, err
	}

	return &SearchDeploymentPipelineOutput{DeploymentPipeline: deploymentPipeline}, nil
}

// UpdateDeploymentPipelineInput represents the input of a UpdateDeploymentPipeline operation.
type UpdateDeploymentPipelineInput struct {
	_                    struct{}
	DeploymentPipelineID *int64
	DeploymentPipeline   *DeploymentPipeline

	// describes optional properties that should be included in the response
	// possible values: "integrationUrl"
	Include []*string
}

// UpdateDeploymentPipelineOutput represents the output of a UpdateDeploymentPipeline operation.
type UpdateDeploymentPipelineOutput struct {
	_                  struct{}
	DeploymentPipeline *DeploymentPipeline
}

// UpdateDeploymentPipeline updates an existing deployment pipeline resource. https://api.ilert.com/api-docs/#tag/deployment-pipelines/put/deployment-pipelines/{id}
func (c *Client) UpdateDeploymentPipeline(input *UpdateDeploymentPipelineInput) (*UpdateDeploymentPipelineOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.DeploymentPipeline == nil {
		return nil, errors.New("deployment pipeline input is required")
	}
	if input.DeploymentPipelineID == nil {
		return nil, errors.New("deployment pipeline id is required")
	}

	q := url.Values{}
	for _, include := range input.Include {
		q.Add("include", *include)
	}

	resp, err := c.httpClient.R().SetBody(input.DeploymentPipeline).Put(fmt.Sprintf("%s/%d?%s", apiRoutes.deploymentPipelines, *input.DeploymentPipelineID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	deploymentPipeline := &DeploymentPipeline{}
	err = json.Unmarshal(resp.Body(), deploymentPipeline)
	if err != nil {
		return nil, err
	}

	return &UpdateDeploymentPipelineOutput{DeploymentPipeline: deploymentPipeline}, nil
}

// DeleteDeploymentPipelineInput represents the input of a DeleteDeploymentPipeline operation.
type DeleteDeploymentPipelineInput struct {
	_                    struct{}
	DeploymentPipelineID *int64
}

// DeleteDeploymentPipelineOutput represents the output of a DeleteDeploymentPipeline operation.
type DeleteDeploymentPipelineOutput struct {
	_ struct{}
}

// DeleteDeploymentPipeline deletes the specified deployment pipeline resource. https://api.ilert.com/api-docs/#tag/deployment-pipelines/delete/deployment-pipelines/{id}
func (c *Client) DeleteDeploymentPipeline(input *DeleteDeploymentPipelineInput) (*DeleteDeploymentPipelineOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.DeploymentPipelineID == nil {
		return nil, errors.New("deployment pipeline id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.deploymentPipelines, *input.DeploymentPipelineID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteDeploymentPipelineOutput{}, nil
}
