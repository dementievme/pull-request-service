package domain

import (
	"context"

	domain "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	Update(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	ListByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error)
}
