package service

import (
	"context"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

type Dependencies struct{}

type Service struct {
	deps Dependencies
}

func New(deps Dependencies) *Service {
	return &Service{deps: deps}
}

var _ domain.AssignmentService = (*Service)(nil)

func (s *Service) CreatePullRequest(ctx context.Context, input domain.CreatePullRequestInput) (domain.PullRequest, error) {
	return domain.PullRequest{}, domain.ErrNotImplemented
}

func (s *Service) MergePullRequest(ctx context.Context, input domain.MergePullRequestInput) (domain.PullRequest, error) {
	return domain.PullRequest{}, domain.ErrNotImplemented
}

func (s *Service) ReassignPullRequest(ctx context.Context, input domain.ReassignPullRequestInput) (domain.ReassignPullRequestResult, error) {
	return domain.ReassignPullRequestResult{}, domain.ErrNotImplemented
}

func (s *Service) UpsertTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	return domain.Team{}, domain.ErrNotImplemented
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	return domain.Team{}, domain.ErrNotImplemented
}

func (s *Service) GetReviewerAssignments(ctx context.Context, userID string) (domain.ReviewerAssignments, error) {
	return domain.ReviewerAssignments{}, domain.ErrNotImplemented
}

func (s *Service) SetUserIsActive(ctx context.Context, input domain.SetUserIsActiveInput) (domain.User, error) {
	return domain.User{}, domain.ErrNotImplemented
}
