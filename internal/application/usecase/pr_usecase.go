package application

import (
	"context"
	"slices"
	"time"

	dto "github.com/dementievme/pull-request-service/internal/application/dto"
	errors "github.com/dementievme/pull-request-service/internal/application/errors"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type PullRequestUseCase struct {
	repo domain.PullRequestRepository
}

func NewPullRequestUseCase(prRepo domain.PullRequestRepository) *PullRequestUseCase {
	return &PullRequestUseCase{repo: prRepo}
}

func (u *PullRequestUseCase) GetByID(ctx context.Context, prID string) (*dto.PullRequestShortDTO, error) {
	pr, err := u.repo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	return &dto.PullRequestShortDTO{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}, nil
}

func (u *PullRequestUseCase) CreatePR(ctx context.Context, prDTO *dto.CreatePullRequestDTO, team *dto.TeamDTO) (*dto.PullRequestDTO, error) {
	reviewerIDs := make([]string, 0, 2)

	for _, m := range team.TeamMembers {
		if len(reviewerIDs) == 2 {
			break
		}
		if m.UserID == prDTO.AuthorID {
			continue
		}
		reviewerIDs = append(reviewerIDs, m.UserID)
	}

	now := time.Now()

	pr := &entity.PullRequest{
		ID:                  prDTO.PullRequestID,
		Name:                prDTO.PullRequestName,
		AuthorID:            prDTO.AuthorID,
		Status:              entity.OPEN,
		AssignedReviewerIDs: reviewerIDs,
		CreatedAt:           &now,
	}

	if err := u.repo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return &dto.PullRequestDTO{
		PullRequestID:       pr.ID,
		PullRequestName:     pr.Name,
		AuthorID:            pr.AuthorID,
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.AssignedReviewerIDs,
		CreatedAt:           pr.CreatedAt,
	}, nil
}

func (u *PullRequestUseCase) Merge(ctx context.Context, prID string) (*dto.PullRequestDTO, error) {
	pr, err := u.repo.Merge(ctx, prID)
	if err != nil {
		return nil, err
	}

	return &dto.PullRequestDTO{
		PullRequestID:       pr.ID,
		PullRequestName:     pr.Name,
		AuthorID:            pr.AuthorID,
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.AssignedReviewerIDs,
		CreatedAt:           pr.CreatedAt,
		MergedAt:            pr.MergedAt,
	}, nil
}

func (u *PullRequestUseCase) Reassign(ctx context.Context, prDTO *dto.ReassignPullRequestDTO, activeMembers []*dto.TeamMemberDTO) (*dto.PullRequestDTO, string, error) {
	pr, err := u.repo.GetByID(ctx, prDTO.PullRequestID)
	if err != nil {
		return nil, "", errors.ErrNotFound
	}

	if pr.Status == entity.MERGED {
		return nil, "", errors.ErrPRMerged
	}

	if len(pr.AssignedReviewerIDs) == 0 {
		return nil, "", errors.ErrNoCandidate
	}

	if !slices.Contains(pr.AssignedReviewerIDs, prDTO.OldReviewerID) {
		return nil, "", errors.ErrNotAssigned
	}

	replacedBy := ""

	for _, activeMember := range activeMembers {
		if slices.Contains(pr.AssignedReviewerIDs, activeMember.UserID) {
			continue
		}

		index := slices.Index(pr.AssignedReviewerIDs, prDTO.OldReviewerID)
		pr.AssignedReviewerIDs[index] = activeMember.UserID
		replacedBy = activeMember.UserID
		break
	}

	if replacedBy == "" {
		return nil, "", errors.ErrNoCandidate
	}

	if err := u.repo.Update(ctx, pr); err != nil {
		return nil, "", err
	}

	return &dto.PullRequestDTO{
		PullRequestID:       pr.ID,
		PullRequestName:     pr.Name,
		AuthorID:            pr.AuthorID,
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.AssignedReviewerIDs,
		CreatedAt:           pr.CreatedAt,
		MergedAt:            pr.MergedAt,
	}, replacedBy, nil
}

func (u *PullRequestUseCase) GetForReviewer(ctx context.Context, reviewerID string) ([]*dto.PullRequestShortDTO, error) {
	prList, err := u.repo.FindByReviewer(ctx, reviewerID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.PullRequestShortDTO, 0, len(prList))

	for _, pr := range prList {
		result = append(result, &dto.PullRequestShortDTO{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	return result, nil
}
