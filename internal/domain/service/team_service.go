package domain

import (
	"context"

	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type TeamService struct {
	repo domain.TeamRepository
}

func NewTeamService(repo domain.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) Create(ctx context.Context, team *entity.Team) error {
	return s.repo.Create(ctx, team)
}

func (s *TeamService) GetByName(ctx context.Context, teamName string) (*entity.Team, error) {
	return s.repo.GetByName(ctx, teamName)
}

func (s *TeamService) FindActiveUsers(ctx context.Context, teamName string) ([]*entity.User, error) {
	return s.repo.FindActiveUsers(ctx, teamName)
}
