package domain

import "time"

type PRStatus string

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)

type PullRequest struct {
	ID                  string
	Name                string
	AuthorID            int64
	Status              PRStatus
	AssignedReviewerIDs []int64
	CreatedAt           *time.Time
	MergedAt            *time.Time
}
