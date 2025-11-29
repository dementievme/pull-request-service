package domain

import domain "github.com/dementievme/pull-request-service/internal/domain/repository"

type Services struct {
	PullRequestService domain.PullRequestRepository
	TeamService        domain.TeamRepository
	UserService        domain.UserRepository
}

func NewServices(pullRequestService domain.PullRequestRepository, teamService domain.TeamRepository, userService domain.UserRepository) *Services {
	return &Services{
		PullRequestService: pullRequestService,
		TeamService:        teamService,
		UserService:        userService,
	}
}
