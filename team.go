package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Team definition https://api.ilert.com/api-docs/#tag/Teams
type Team struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Private bool         `json:"private"`
	Members []TeamMember `json:"members"`
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
		return nil, errors.New("Input is required")
	}
	if input.Team == nil {
		return nil, errors.New("Team input is required")
	}
	resp, err := c.httpClient.R().SetBody(input.Team).Post(fmt.Sprintf("%s", apiRoutes.teams))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 201); err != nil {
		return nil, err
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	output := &CreateTeamOutput{Team: team}

	return output, nil
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
		return nil, errors.New("Input is required")
	}
	if input.TeamID == nil {
		return nil, errors.New("Team id is required")
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d", apiRoutes.teams, *input.TeamID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	output := &GetTeamOutput{
		Team: team,
	}

	return output, nil
}

// GetTeamsInput represents the input of a GetTeams operation.
type GetTeamsInput struct {
	_ struct{}
}

// GetTeamsOutput represents the output of a GetTeams operation.
type GetTeamsOutput struct {
	_     struct{}
	Teams []*Team
}

// GetTeams gets list teams. https://api.ilert.com/api-docs/#tag/Teams/paths/~1teams/get
func (c *Client) GetTeams(input *GetTeamsInput) (*GetTeamsOutput, error) {
	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s", apiRoutes.teams))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	teams := make([]*Team, 0)
	err = json.Unmarshal(resp.Body(), &teams)
	if err != nil {
		return nil, err
	}

	output := &GetTeamsOutput{Teams: teams}

	return output, nil
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
		return nil, errors.New("Input is required")
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
	if err = catchGenericAPIError(resp, 200); err != nil {
		return nil, err
	}

	team := &Team{}
	err = json.Unmarshal(resp.Body(), team)
	if err != nil {
		return nil, err
	}

	output := &UpdateTeamOutput{Team: team}

	return output, nil
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
		return nil, errors.New("Input is required")
	}
	if input.TeamID == nil {
		return nil, errors.New("Team id is required")
	}

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d", apiRoutes.teams, *input.TeamID))
	if err != nil {
		return nil, err
	}
	if err = catchGenericAPIError(resp, 204); err != nil {
		return nil, err
	}

	output := &DeleteTeamOutput{}
	return output, nil
}
