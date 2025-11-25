package application

import "time"

type PullRequestDTO struct {
	PullRequestID       string     `json:"pull_request_id"`
	PullRequestName     string     `json:"pull_request_name"`
	AuthorID            string     `json:"author_id"`
	Status              string     `json:"status"`
	AssignedReviewerIDs []string   `json:"assigned_reviewers"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	MergedAt            *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type CreatePullRequestDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePullRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPullRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_user_id"`
}

type ReassignPullRequestResponseDTO struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}
