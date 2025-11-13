package domain

import "time"

// PullRequestStatus is the status of a pull request.
type PullRequestStatus string

// PullRequest statuses.
const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

// PullRequest is a pull request.
type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PullRequestStatus
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}

// PullRequestShort is a short pull request.
type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   PullRequestStatus
}

// Team is a team.
type Team struct {
	Name    string
	Members []TeamMember
}

// TeamMember is a team member.
type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

// User is a user.
type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

// CreatePullRequestInput is the input for creating a pull request.
type CreatePullRequestInput struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
}

// MergePullRequestInput is the input for merging a pull request.
type MergePullRequestInput struct {
	PullRequestID string
}

// ReassignPullRequestInput is the input for reassigning a pull request.
type ReassignPullRequestInput struct {
	PullRequestID string
	OldUserID     string
}

// ReassignPullRequestResult is the result of reassigning a pull request.
type ReassignPullRequestResult struct {
	PullRequest PullRequest
	ReplacedBy  string
}

// ReviewerAssignments is the assignments of reviewers to pull requests.
type ReviewerAssignments struct {
	UserID       string
	PullRequests []PullRequestShort
}

// SetUserIsActiveInput is the input for setting the active status of a user.
type SetUserIsActiveInput struct {
	UserID   string
	IsActive bool
}
