package application

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetActiveUserDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserPRsDTO struct {
	UserID       string                 `json:"user_id"`
	PullRequests []*PullRequestShortDTO `json:"pull_requests"`
}
