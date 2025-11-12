package domain

import "context"

type AssignmentService interface {
	CreatePullRequest(ctx context.Context, input CreatePullRequestInput) (PullRequest, error)
	MergePullRequest(ctx context.Context, input MergePullRequestInput) (PullRequest, error)
	ReassignPullRequest(ctx context.Context, input ReassignPullRequestInput) (ReassignPullRequestResult, error)
	UpsertTeam(ctx context.Context, team Team) (Team, error)
	GetTeam(ctx context.Context, teamName string) (Team, error)
	GetReviewerAssignments(ctx context.Context, userID string) (ReviewerAssignments, error)
	SetUserIsActive(ctx context.Context, input SetUserIsActiveInput) (User, error)
}
