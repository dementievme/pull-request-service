package application

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName    string           `json:"team_name"`
	TeamMembers []*TeamMemberDTO `json:"members"`
}
