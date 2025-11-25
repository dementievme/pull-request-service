package application

import (
	"context"
	"errors"
	"strconv"
	"time"

	application "github.com/dementievme/pull-request-service/internal/application/dto"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type PullRequestUseCase struct {
	repo domain.PullRequestRepository
}

func NewPullRequestUseCase(repo domain.PullRequestRepository) *PullRequestUseCase {
	return &PullRequestUseCase{repo: repo}
}

func (u *PullRequestUseCase) CreatePR(ctx context.Context, dto *application.CreatePullRequestDTO) (*application.PullRequestDTO, error) {
	pr := &entity.PullRequest{
		ID:        dto.PullRequestID,
		Name:      dto.PullRequestName,
		AuthorID:  parseID(dto.AuthorID),
		Status:    entity.OPEN,
		CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
	}
	if err := u.repo.Create(ctx, pr); err != nil {
		return nil, err
	}
	return &application.PullRequestDTO{
		PullRequestID:       pr.ID,
		PullRequestName:     pr.Name,
		AuthorID:            dto.AuthorID,
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.AssignedReviewerIDs,
		CreatedAt:           pr.CreatedAt,
	}, nil
}

func (u *PullRequestUseCase) Merge(ctx context.Context, prID string) (*application.PullRequestDTO, error) {
	pr, err := u.repo.Merge(ctx, prID)
	if err != nil {
		return nil, err
	}
	return &application.PullRequestDTO{
		PullRequestID:       pr.ID,
		PullRequestName:     pr.Name,
		AuthorID:            formatID(pr.AuthorID),
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.AssignedReviewerIDs,
		CreatedAt:           pr.CreatedAt,
		MergedAt:            pr.MergedAt,
	}, nil
}

func (u *PullRequestUseCase) Reassign(ctx context.Context, dto *application.ReassignPullRequestDTO) (*application.PullRequestDTO, string, error) {
	pr, err := u.repo.GetByID(ctx, dto.PullRequestID)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == entity.MERGED {
		return nil, "", errors.New("cannot reassign on merged PR")
	}
	if len(pr.AssignedReviewerIDs) == 0 {
		return nil, "", errors.New("no reviewer assigned")
	}
	replacedBy := "new_reviewer_id"
	pr.AssignedReviewerIDs[0] = replacedBy
	if err := u.repo.Update(ctx, pr); err != nil {
		return nil, "", err
	}
	return &application.PullRequestDTO{
		PullRequestID:       pr.ID,
		PullRequestName:     pr.Name,
		AuthorID:            formatID(pr.AuthorID),
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.AssignedReviewerIDs,
		CreatedAt:           pr.CreatedAt,
		MergedAt:            pr.MergedAt,
	}, replacedBy, nil
}

func (u *PullRequestUseCase) GetForReviewer(ctx context.Context, userID string) ([]*application.PullRequestShortDTO, error) {
	prs, err := u.repo.FindByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []*application.PullRequestShortDTO
	for _, pr := range prs {
		result = append(result, &application.PullRequestShortDTO{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        formatID(pr.AuthorID),
			Status:          string(pr.Status),
		})
	}
	return result, nil
}

func parseID(id string) int64 {
	v, _ := strconv.ParseInt(id, 10, 64)
	return v
}

func formatID(id int64) string {
	return strconv.FormatInt(id, 10)
}
