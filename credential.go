package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Credential definition, an early access resource behind /api/beta/credentials.
// Its contract may change or be withdrawn while it is in early access.
//
// A credential holds a secret that ilert uses against another system, for
// example the MCP servers of an agent. It cannot be changed once it is created:
// the API rejects every update, so a changed credential has to be deleted and
// created again. Names are unique within an account.
type Credential struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	TypeHint string `json:"typeHint,omitempty"`

	// SecretValue is only sent when the credential is created. The API stores it
	// encrypted and never returns it. How it is used depends on TypeHint:
	//   - NONE takes no secret: HostURL carries the whole URL, and the API keeps
	//     its path and query as the secret and returns only the scheme, host and
	//     port as HostURL
	//   - BASIC requires Username and stores base64(username:secret)
	//   - BEARER and API_KEY store the secret as it is, BEARER takes no Username
	//   - OAUTH2 takes the client secret, or the refresh token with the
	//     REFRESH_TOKEN method, and no Username
	//   - AWS_SIGV4 takes the secret access key, with the access key id as
	//     Username unless the credential assumes RoleArn
	SecretValue string `json:"secretValue,omitempty"`
	Username    string `json:"username,omitempty"`

	// HostURL is the target the secret is used against. A credential in the MCP
	// category needs a public https host.
	HostURL string `json:"hostUrl,omitempty"`

	// IdpURL, IdpScopes, IdpClient, Audience and OAuthMethod configure the token
	// endpoint of an OAUTH2 credential. IdpURL is required for one, IdpClient
	// with the CLIENT_SECRET_BASIC and CLIENT_SECRET_POST methods.
	IdpURL      string `json:"idpUrl,omitempty"`
	IdpScopes   string `json:"idpScopes,omitempty"` // space separated
	IdpClient   string `json:"idpClient,omitempty"`
	Audience    string `json:"audience,omitempty"`
	OAuthMethod string `json:"oauthMethod,omitempty"`

	// Category groups credentials by what uses them. Agents only take
	// credentials in the MCP category for their MCP servers.
	Category    string `json:"category,omitempty"`
	SubCategory string `json:"subCategory,omitempty"`

	// APIKeyHeaderName is the header an API_KEY credential is sent in
	APIKeyHeaderName string `json:"apiKeyHeaderName,omitempty"`

	// RoleArn and Region are used by AWS_SIGV4 credentials
	RoleArn string `json:"roleArn,omitempty"`
	Region  string `json:"region,omitempty"`

	Teams []TeamShort `json:"teams,omitempty"`
}

// CredentialType defines how a credential is presented to the system it is used against
var CredentialType = struct {
	None     string
	Unknown  string
	Bearer   string
	APIKey   string
	Basic    string
	OAuth2   string
	AWSSigV4 string
}{
	None:     "NONE",
	Unknown:  "UNKNOWN",
	Bearer:   "BEARER",
	APIKey:   "API_KEY",
	Basic:    "BASIC",
	OAuth2:   "OAUTH2",
	AWSSigV4: "AWS_SIGV4",
}

// CredentialTypeAll defines credential type list
var CredentialTypeAll = []string{
	CredentialType.None,
	CredentialType.Unknown,
	CredentialType.Bearer,
	CredentialType.APIKey,
	CredentialType.Basic,
	CredentialType.OAuth2,
	CredentialType.AWSSigV4,
}

// CredentialOAuthMethod defines the client authentication at the token endpoint
// of an OAUTH2 credential. CLIENT_SECRET_BASIC and CLIENT_SECRET_POST use the
// client credentials grant, REFRESH_TOKEN the refresh token grant.
var CredentialOAuthMethod = struct {
	ClientSecretBasic string
	ClientSecretPost  string
	RefreshToken      string
}{
	ClientSecretBasic: "CLIENT_SECRET_BASIC",
	ClientSecretPost:  "CLIENT_SECRET_POST",
	RefreshToken:      "REFRESH_TOKEN",
}

// CredentialOAuthMethodAll defines credential oauth method list
var CredentialOAuthMethodAll = []string{
	CredentialOAuthMethod.ClientSecretBasic,
	CredentialOAuthMethod.ClientSecretPost,
	CredentialOAuthMethod.RefreshToken,
}

// CreateCredentialInput represents the input of a CreateCredential operation.
type CreateCredentialInput struct {
	_          struct{}
	Credential *Credential
}

// CreateCredentialOutput represents the output of a CreateCredential operation.
type CreateCredentialOutput struct {
	_          struct{}
	Credential *Credential
}

// CreateCredential creates a new credential resource.
func (c *Client) CreateCredential(input *CreateCredentialInput) (*CreateCredentialOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Credential == nil {
		return nil, errors.New("credential input is required")
	}

	resp, err := c.betaRequest().SetBody(input.Credential).Post(apiRoutes.credentials)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	credential := &Credential{}
	err = json.Unmarshal(resp.Body(), credential)
	if err != nil {
		return nil, err
	}

	return &CreateCredentialOutput{Credential: credential}, nil
}

// GetCredentialInput represents the input of a GetCredential operation.
type GetCredentialInput struct {
	_            struct{}
	CredentialID *int64
}

// GetCredentialOutput represents the output of a GetCredential operation.
type GetCredentialOutput struct {
	_          struct{}
	Credential *Credential
}

// GetCredential gets the credential resource with specified id.
func (c *Client) GetCredential(input *GetCredentialInput) (*GetCredentialOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CredentialID == nil {
		return nil, errors.New("credential id is required")
	}

	resp, err := c.betaRequest().Get(fmt.Sprintf("%s/%d", apiRoutes.credentials, *input.CredentialID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	credential := &Credential{}
	err = json.Unmarshal(resp.Body(), credential)
	if err != nil {
		return nil, err
	}

	return &GetCredentialOutput{Credential: credential}, nil
}

// GetCredentialsInput represents the input of a GetCredentials operation.
type GetCredentialsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	MaxResults *int
}

// GetCredentialsOutput represents the output of a GetCredentials operation.
type GetCredentialsOutput struct {
	_           struct{}
	Credentials []*Credential
}

// GetCredentials lists existing credential resources.
func (c *Client) GetCredentials(input *GetCredentialsInput) (*GetCredentialsOutput, error) {
	if input == nil {
		input = &GetCredentialsInput{}
	}

	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.betaRequest().Get(fmt.Sprintf("%s?%s", apiRoutes.credentials, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	credentials := make([]*Credential, 0)
	err = json.Unmarshal(resp.Body(), &credentials)
	if err != nil {
		return nil, err
	}

	return &GetCredentialsOutput{Credentials: credentials}, nil
}

// SearchCredentialInput represents the input of a SearchCredential operation.
type SearchCredentialInput struct {
	_              struct{}
	CredentialName *string
}

// SearchCredentialOutput represents the output of a SearchCredential operation.
type SearchCredentialOutput struct {
	_          struct{}
	Credential *Credential
}

// SearchCredential gets the credential resource with specified name.
func (c *Client) SearchCredential(input *SearchCredentialInput) (*SearchCredentialOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CredentialName == nil {
		return nil, errors.New("credential name is required")
	}

	resp, err := c.betaRequest().Get(fmt.Sprintf("%s/name/%s", apiRoutes.credentials, url.PathEscape(*input.CredentialName)))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	credential := &Credential{}
	err = json.Unmarshal(resp.Body(), credential)
	if err != nil {
		return nil, err
	}

	return &SearchCredentialOutput{Credential: credential}, nil
}

// DeleteCredentialInput represents the input of a DeleteCredential operation.
type DeleteCredentialInput struct {
	_            struct{}
	CredentialID *int64
}

// DeleteCredentialOutput represents the output of a DeleteCredential operation.
type DeleteCredentialOutput struct {
	_ struct{}
}

// DeleteCredential deletes the specified credential resource. The API does not
// check whether an agent still uses the credential.
func (c *Client) DeleteCredential(input *DeleteCredentialInput) (*DeleteCredentialOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.CredentialID == nil {
		return nil, errors.New("credential id is required")
	}

	resp, err := c.betaRequest().Delete(fmt.Sprintf("%s/%d", apiRoutes.credentials, *input.CredentialID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteCredentialOutput{}, nil
}
