package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	errors "github.com/dementievme/pull-request-service/internal/application/errors"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type PostgrePullRequestRepository struct {
	db *sql.DB
}

func NewPostgrePullRequestRepository(db *sql.DB) *PostgrePullRequestRepository {
	return &PostgrePullRequestRepository{db: db}
}

func (r *PostgrePullRequestRepository) Create(ctx context.Context, pr *entity.PullRequest) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pull_requests (id, name, author_id, status, assigned_reviewers, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		pr.ID, pr.Name, pr.AuthorID, string(pr.Status), pq.Array(pr.AssignedReviewerIDs), pr.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgrePullRequestRepository) GetByID(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr := &entity.PullRequest{}

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, author_id, status, assigned_reviewers, created_at, merged_at
		 FROM pull_requests WHERE id=$1`,
		prID,
	).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, pq.Array(&pr.AssignedReviewerIDs), &pr.CreatedAt, &pr.MergedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return pr, nil
}

func (r *PostgrePullRequestRepository) Update(ctx context.Context, pr *entity.PullRequest) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pull_requests SET name=$1, status=$2, assigned_reviewers=$3, merged_at=$4 WHERE id=$5`,
		pr.Name, string(pr.Status), pq.Array(pr.AssignedReviewerIDs), pr.MergedAt, pr.ID,
	)
	return err
}

func (r *PostgrePullRequestRepository) Merge(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr, err := r.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == entity.MERGED {
		return pr, nil
	}

	now := time.Now()
	_, err = r.db.ExecContext(ctx,
		`UPDATE pull_requests SET status='MERGED', merged_at=$1 WHERE id=$2`,
		now, prID,
	)
	if err != nil {
		return nil, err
	}

	pr.Status = entity.MERGED
	pr.MergedAt = &now
	return pr, nil
}

func (r *PostgrePullRequestRepository) FindByReviewer(ctx context.Context, reviewerID string) ([]*entity.PullRequest, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, author_id, status, assigned_reviewers, created_at, merged_at
		 FROM pull_requests WHERE $1 = ANY(assigned_reviewers)`, reviewerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*entity.PullRequest
	for rows.Next() {
		pr := &entity.PullRequest{}
		var status string
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &status, pq.Array(&pr.AssignedReviewerIDs), &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, rows.Err()
}
