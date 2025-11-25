package domain

import (
	"context"

	domain "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	Update(ctx context.Context, pr *domain.PullRequest) error
	Merge(ctx context.Context, prID string) (*domain.PullRequest, error)
	FindByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}
