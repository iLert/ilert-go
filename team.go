package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Team definition https://api.ilert.com/api-docs/#tag/Teams
type Team struct {
	ID         int64        `json:"id"`
	Name       string       `json:"name"`
	Visibility string       `json:"visibility"`
	Members    []TeamMember `json:"members"`
}

// TeamMember definition
type TeamMember struct {
	User User   `json:"user"`
	Role string `json:"role"` // "ADMIN" or "USER" or "RESPONDER" or "STAKEHOLDER"
}

// TeamShort definition
type TeamShort struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
}

// TeamMemberRoles defines team member roles
var TeamMemberRoles = struct {
	Admin       string
	User        string
	Responder   string
	Stakeholder string
}{
	Admin:       "ADMIN",
	User:        "USER",
	Responder:   "RESPONDER",
	Stakeholder: "STAKEHOLDER",
}

// TeamMemberRolesAll defines team member roles list
var TeamMemberRolesAll = []string{
	TeamMemberRoles.Admin,
	TeamMemberRoles.User,
	TeamMemberRoles.Responder,
	TeamMemberRoles.Stakeholder,
}

// TeamVisibility defines team visibility
var TeamVisibility = struct {
	Public  string
	Private string
}{
	Public:  "PUBLIC",
	Private: "PRIVATE",
}

// TeamVisibilityAll defines team visibility list
var TeamVisibilityAll = []string{
	TeamVisibility.Public,
	TeamVisibility.Private,
}

// Ownership definition
type Ownership struct {
	TeamID int64 `json:"teamId"`
}

// CreateTeamInput represents the input of a CreateTeam operation.
type CreateTeamInput struct {
	_    struct{}
	Team *Team
}

// CreateTeamOutput represents the output of a CreateTeam operation.
type CreateTeamOutput struct {
	_    struct{}
	Team *Team
}

// CreateTeam creates a new team. https://api.ilert.com/api-docs/#tag/Teams/paths/~1teams/posts
func (c *Client) CreateTeam(input *CreateTeamInput) (*CreateTeamOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Team == nil {
		return nil, errors.New("Team input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Team).Post(apiRoutes.teams)
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	return &CreateTeamOutput{Team: team}, nil
}

// GetTeamInput represents the input of a GetTeam operation.
type GetTeamInput struct {
	_      struct{}
	TeamID *int64
}

// GetTeamOutput represents the output of a GetTeam operation.
type GetTeamOutput struct {
	_    struct{}
	Team *Team
}

// GetTeam gets the team with specified id. https://api.ilert.com/api-docs/#tag/Teams/paths/~1teams~1{id}/get
func (c *Client) GetTeam(input *GetTeamInput) (*GetTeamOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TeamID == nil {
		return nil, errors.New("Team id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.teams, *input.TeamID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	return &GetTeamOutput{Team: team}, nil
}

// GetTeamsInput represents the input of a GetTeams operation.
type GetTeamsInput struct {
	_ struct{}

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetTeamsOutput represents the output of a GetTeams operation.
type GetTeamsOutput struct {
	_     struct{}
	Teams []*Team
}

// GetTeams gets list teams. https://api.ilert.com/api-docs/#tag/Teams/paths/~1teams/get
func (c *Client) GetTeams(input *GetTeamsInput) (*GetTeamsOutput, error) {
	q := url.Values{}
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s?%s", apiRoutes.teams, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	teams := make([]*Team, 0)
	err = json.Unmarshal(resp.Body(), &teams)
	if err != nil {
		return nil, err
	}

	return &GetTeamsOutput{Teams: teams}, nil
}

// SearchTeamInput represents the input of a SearchTeam operation.
type SearchTeamInput struct {
	_        struct{}
	TeamName *string
}

// SearchTeamOutput represents the output of a SearchTeam operation.
type SearchTeamOutput struct {
	_    struct{}
	Team *Team
}

// SearchTeam gets the team with specified name.
func (c *Client) SearchTeam(input *SearchTeamInput) (*SearchTeamOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TeamName == nil {
		return nil, errors.New("team name is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/name/%s", apiRoutes.teams, *input.TeamName))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	return &SearchTeamOutput{Team: team}, nil
}

// UpdateTeamInput represents the input of a UpdateTeam operation.
type UpdateTeamInput struct {
	_      struct{}
	TeamID *int64
	Team   *Team
}

// UpdateTeamOutput represents the output of a UpdateTeam operation.
type UpdateTeamOutput struct {
	_    struct{}
	Team *Team
}

// UpdateTeam updates an existing team. https://api.ilert.com/api-docs/#tag/Teams/paths/~1teams~1{id}/put
func (c *Client) UpdateTeam(input *UpdateTeamInput) (*UpdateTeamOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.Team == nil {
		return nil, errors.New("Team input is required")
	}
	if input.TeamID == nil {
		return nil, errors.New("Team id is required")
	}

	resp, err := c.httpClient.R().SetBody(input.Team).Put(fmt.Sprintf("%s/%d", apiRoutes.teams, *input.TeamID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	return &UpdateTeamOutput{Team: team}, nil
}

// DeleteTeamInput represents the input of a DeleteTeam operation.
type DeleteTeamInput struct {
	_      struct{}
	TeamID *int64
}

// DeleteTeamOutput represents the output of a DeleteTeam operation.
type DeleteTeamOutput struct {
	_ struct{}
}

// DeleteTeam deletes the specified alert source. https://api.ilert.com/api-docs/#tag/Teams/paths/~1teams~1{id}/delete
func (c *Client) DeleteTeam(input *DeleteTeamInput) (*DeleteTeamOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TeamID == nil {
		return nil, errors.New("Team id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.teams, *input.TeamID))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 204); apiErr != nil {
		return nil, apiErr
	}

	return &DeleteTeamOutput{}, nil
}
