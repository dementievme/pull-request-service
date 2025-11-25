package infrastructure

import (
	"database/sql"

	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type Repositories struct {
	UserRepo        domain.UserRepository
	TeamRepo        domain.TeamRepository
	PullRequestRepo domain.PullRequestRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:        NewPostgreUserRepository(db),
		TeamRepo:        NewPostgreTeamRepository(db),
		PullRequestRepo: NewPostgrePullRequestRepository(db),
	}
}
