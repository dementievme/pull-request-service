package infrastructure

import (
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type Repositories struct {
	UserRepo        domain.UserRepository
	TeamRepo        domain.TeamRepository
	PullRequestRepo domain.PullRequestRepository
}

// TODO [TASK #6] Realize interfaces for CRUD
// func NewRepositories(db *sql.DB) *Repositories {
// 	return &Repositories{
// 		UserRepo:        NewUserRepository(db),
// 		TeamRepo:        NewTeamRepository(db),
// 		PullRequestRepo: NewPullRequestRepository(db),
// 	}
// }
