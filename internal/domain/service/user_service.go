package domain

import (
	"context"

	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user []*entity.User) ([]*entity.User, error) {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *UserService) UpdateIsActive(ctx context.Context, userID string, isActive bool) error {
	return s.repo.UpdateIsActive(ctx, userID, isActive)
}

func (s *UserService) FindByTeam(ctx context.Context, teamName string) ([]*entity.User, error) {
	return s.repo.FindByTeam(ctx, teamName)
}
