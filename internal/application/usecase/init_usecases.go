package application

import (
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type UseCases struct {
	PullRequestUseCase *PullRequestUseCase
	TeamUseCase        *TeamUseCase
	UserUseCase        *UserUseCase
}

func NewUseCases(pullRequestUseCase domain.PullRequestRepository, teamUseCase domain.TeamRepository, userUseCase domain.UserRepository) *UseCases {
	return &UseCases{
		PullRequestUseCase: NewPullRequestUseCase(pullRequestUseCase),
		TeamUseCase:        NewTeamUseCase(teamUseCase),
		UserUseCase:        NewUserUseCase(userUseCase),
	}
}
