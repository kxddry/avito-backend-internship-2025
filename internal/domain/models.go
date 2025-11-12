package domain

import "time"

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PullRequestStatus
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   PullRequestStatus
}

type Team struct {
	Name    string
	Members []TeamMember
}

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

type CreatePullRequestInput struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
}

type MergePullRequestInput struct {
	PullRequestID string
}

type ReassignPullRequestInput struct {
	PullRequestID string
	OldUserID     string
}

type ReassignPullRequestResult struct {
	PullRequest PullRequest
	ReplacedBy  string
}

type ReviewerAssignments struct {
	UserID       string
	PullRequests []PullRequestShort
}

type SetUserIsActiveInput struct {
	UserID   string
	IsActive bool
}
