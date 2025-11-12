package service

import (
	"context"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) error
	GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	Update(ctx context.Context, pr domain.PullRequest) error
	Delete(ctx context.Context, pullRequestID string) error
}

type TeamRepository interface {
	Create(ctx context.Context, team domain.Team) error
	GetByName(ctx context.Context, teamName string) (domain.Team, error)
	Update(ctx context.Context, team domain.Team) error
	Delete(ctx context.Context, teamName string) error
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, userID string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, userID string) error
}

type RepositoryProvider interface {
	PullRequests() PullRequestRepository
	Teams() TeamRepository
	Users() UserRepository
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context, RepositoryProvider) error) error
	WithROTx(ctx context.Context, fn func(context.Context, RepositoryProvider) error) error
}

type Storage struct {
	PullRequestRepo PullRequestRepository
	TeamRepo        TeamRepository
	UserRepo        UserRepository
}

func (s Storage) PullRequests() PullRequestRepository {
	return s.PullRequestRepo
}

func (s Storage) Teams() TeamRepository {
	return s.TeamRepo
}

func (s Storage) Users() UserRepository {
	return s.UserRepo
}

type Dependencies struct {
	Repositories       RepositoryProvider
	TransactionManager TxManager
}

type Service struct {
	repos RepositoryProvider
	tx    TxManager
}

func New(deps Dependencies) *Service {
	repos := deps.Repositories
	if repos == nil {
		repos = Storage{}
	}

	tx := deps.TransactionManager
	if tx == nil {
		tx = noopTxManager{repos: repos}
	}

	return &Service{
		repos: repos,
		tx:    tx,
	}
}

var _ domain.AssignmentService = (*Service)(nil)

func (s *Service) CreatePullRequest(ctx context.Context, input domain.CreatePullRequestInput) (domain.PullRequest, error) {
	var result domain.PullRequest
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

func (s *Service) MergePullRequest(ctx context.Context, input domain.MergePullRequestInput) (domain.PullRequest, error) {
	var result domain.PullRequest
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

func (s *Service) ReassignPullRequest(ctx context.Context, input domain.ReassignPullRequestInput) (domain.ReassignPullRequestResult, error) {
	var result domain.ReassignPullRequestResult
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

func (s *Service) UpsertTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	var result domain.Team
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	var result domain.Team
	err := s.tx.WithROTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

func (s *Service) GetReviewerAssignments(ctx context.Context, userID string) (domain.ReviewerAssignments, error) {
	var result domain.ReviewerAssignments
	err := s.tx.WithROTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

func (s *Service) SetUserIsActive(ctx context.Context, input domain.SetUserIsActiveInput) (domain.User, error) {
	var result domain.User
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos RepositoryProvider) error {
		return domain.ErrNotImplemented
	})
	return result, err
}

type noopTxManager struct {
	repos RepositoryProvider
}

func (n noopTxManager) WithTx(ctx context.Context, fn func(context.Context, RepositoryProvider) error) error {
	if fn == nil {
		return nil
	}
	return fn(ctx, n.repos)
}

func (n noopTxManager) WithROTx(ctx context.Context, fn func(context.Context, RepositoryProvider) error) error {
	return n.WithTx(ctx, fn)
}
