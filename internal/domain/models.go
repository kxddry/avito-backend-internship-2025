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

// Stats is the statistics.
type Stats struct {
	Users StatsUsers `json:"users"`
	PRs   StatsPRs   `json:"prs"`
	Teams StatsTeams `json:"teams"`
}

// StatsUsers is the user statistics.
type StatsUsers struct {
	Total    int              `json:"total"`
	Active   int              `json:"active"`
	Inactive int              `json:"inactive"`
	ByUser   []StatsUserEntry `json:"byUser"`
}

// StatsUserEntry is a user statistics entry.
type StatsUserEntry struct {
	UserID               string `json:"userId"`
	UserName             string `json:"userName"`
	Team                 string `json:"team"`
	IsActive             bool   `json:"isActive"`
	AssignedReviewsTotal int    `json:"assignedReviewsTotal"`
	OpenReviews          int    `json:"openReviews"`
	MergedReviews        int    `json:"mergedReviews"`
}

// StatsPRs is the PR statistics.
type StatsPRs struct {
	Total          int `json:"total"`
	Open           int `json:"open"`
	Merged         int `json:"merged"`
	With0Reviewers int `json:"with0Reviewers"`
	With1Reviewer  int `json:"with1Reviewer"`
	With2Reviewers int `json:"with2Reviewers"`
}

// StatsTeams is the team statistics.
type StatsTeams struct {
	Total  int              `json:"total"`
	ByTeam []StatsTeamEntry `json:"byTeam"`
}

// StatsTeamEntry is a team statistics entry.
type StatsTeamEntry struct {
	TeamName        string `json:"teamName"`
	MembersTotal    int    `json:"membersTotal"`
	MembersActive   int    `json:"membersActive"`
	PRsCreatedTotal int    `json:"prsCreatedTotal"`
	PRsOpen         int    `json:"prsOpen"`
}
