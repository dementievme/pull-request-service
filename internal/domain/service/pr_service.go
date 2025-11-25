package domain

import (
	"context"

	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type PullRequestService struct {
	repo domain.PullRequestRepository
}

func NewPullRequestService(repo domain.PullRequestRepository) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (s *PullRequestService) Create(ctx context.Context, pr *entity.PullRequest) error {
	return s.repo.Create(ctx, pr)
}

func (s *PullRequestService) GetByID(ctx context.Context, prID string) (*entity.PullRequest, error) {
	return s.repo.GetByID(ctx, prID)
}

func (s *PullRequestService) Update(ctx context.Context, pr *entity.PullRequest) error {
	return s.repo.Update(ctx, pr)
}

func (s *PullRequestService) Merge(ctx context.Context, prID string) (*entity.PullRequest, error) {
	return s.repo.Merge(ctx, prID)
}

func (s *PullRequestService) FindByReviewer(ctx context.Context, userID string) ([]*entity.PullRequest, error) {
	return s.repo.FindByReviewer(ctx, userID)
}
