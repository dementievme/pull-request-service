package domain

import (
	"context"

	domain "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user []*domain.User) error
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateIsActive(ctx context.Context, userID string, isActive bool) error
	FindByTeam(ctx context.Context, teamName string) ([]*domain.User, error)
}
