package domain

import "context"

// AssignmentService is the service for assigning pull requests to reviewers.
type AssignmentService interface {
	CreatePullRequest(ctx context.Context, input *CreatePullRequestInput) (*PullRequest, error)
	MergePullRequest(ctx context.Context, input *MergePullRequestInput) (*PullRequest, error)
	ReassignPullRequest(ctx context.Context, input *ReassignPullRequestInput) (*ReassignPullRequestResult, error)
	CreateTeam(ctx context.Context, team *Team) (*Team, error)
	GetTeam(ctx context.Context, teamName string) (*Team, error)
	GetReviewerAssignments(ctx context.Context, userID string) (*ReviewerAssignments, error)
	SetUserIsActive(ctx context.Context, input *SetUserIsActiveInput) (*User, error)
	DeactivateTeam(ctx context.Context, teamName string) (int, error)
	SafeReassignPR(ctx context.Context, prID string) (*PullRequest, error)
	GetStats(ctx context.Context) (*Stats, error)
}
