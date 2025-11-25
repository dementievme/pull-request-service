package domain

import (
	"context"

	domain "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	ListByTeam(ctx context.Context, teamName string) ([]*domain.User, error)
}
