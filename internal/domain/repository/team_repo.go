package domain

import (
	"context"

	domain "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByName(ctx context.Context, teamName string) (*domain.Team, error)
	FindActiveUsers(ctx context.Context, teamName string) ([]*domain.User, error)
}
