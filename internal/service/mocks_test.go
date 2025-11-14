package service

import (
	"context"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

// MockTxManager is a mock implementation of storage.TxManager.
type MockTxManager struct {
	DoFunc func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error
}

func (m *MockTxManager) Do(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
	if m.DoFunc != nil {
		return m.DoFunc(ctx, fn)
	}
	return nil
}

// MockTx is a mock implementation of storage.Tx.
type MockTx struct {
	PullRequestRepoFunc func() storage.PullRequestRepository
	TeamRepoFunc        func() storage.TeamRepository
	UserRepoFunc        func() storage.UserRepository
	CommitFunc          func() error
	RollbackFunc        func() error
}

func (m *MockTx) PullRequestRepo() storage.PullRequestRepository {
	if m.PullRequestRepoFunc != nil {
		return m.PullRequestRepoFunc()
	}
	return &MockPullRequestRepository{}
}

func (m *MockTx) TeamRepo() storage.TeamRepository {
	if m.TeamRepoFunc != nil {
		return m.TeamRepoFunc()
	}
	return &MockTeamRepository{}
}

func (m *MockTx) UserRepo() storage.UserRepository {
	if m.UserRepoFunc != nil {
		return m.UserRepoFunc()
	}
	return &MockUserRepository{}
}

func (m *MockTx) Commit() error {
	if m.CommitFunc != nil {
		return m.CommitFunc()
	}
	return nil
}

func (m *MockTx) Rollback() error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc()
	}
	return nil
}

// MockPullRequestRepository is a mock implementation of storage.PullRequestRepository.
type MockPullRequestRepository struct {
	CreateFunc           func(ctx context.Context, pr *domain.PullRequest) error
	GetByIDFunc          func(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	GetPRAssignmentsFunc func(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
	UpdateFunc           func(ctx context.Context, pr *domain.PullRequest) error
	GetStatsFunc         func(ctx context.Context) (*domain.StatsPRs, error)
}

func (m *MockPullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, pr)
	}
	return nil
}

func (m *MockPullRequestRepository) GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, pullRequestID)
	}
	return domain.PullRequest{}, storage.ErrNotFound
}

func (m *MockPullRequestRepository) GetPRAssignments(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) { //nolint:lll
	if m.GetPRAssignmentsFunc != nil {
		return m.GetPRAssignmentsFunc(ctx, reviewerID)
	}
	return nil, storage.ErrNotFound
}

func (m *MockPullRequestRepository) Update(ctx context.Context, pr *domain.PullRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, pr)
	}
	return nil
}

func (m *MockPullRequestRepository) GetStats(ctx context.Context) (*domain.StatsPRs, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx)
	}
	return &domain.StatsPRs{}, nil
}

// MockTeamRepository is a mock implementation of storage.TeamRepository.
type MockTeamRepository struct {
	CreateFunc    func(ctx context.Context, team *domain.Team) error
	GetByNameFunc func(ctx context.Context, teamName string) (domain.Team, error)
	GetStatsFunc  func(ctx context.Context) (*domain.StatsTeams, error)
}

func (m *MockTeamRepository) Create(ctx context.Context, team *domain.Team) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, team)
	}
	return nil
}

func (m *MockTeamRepository) GetByName(ctx context.Context, teamName string) (domain.Team, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, teamName)
	}
	return domain.Team{}, storage.ErrNotFound
}

func (m *MockTeamRepository) GetStats(ctx context.Context) (*domain.StatsTeams, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx)
	}
	return &domain.StatsTeams{}, nil
}

// MockUserRepository is a mock implementation of storage.UserRepository.
type MockUserRepository struct {
	GetByIDFunc     func(ctx context.Context, userID string) (domain.User, error)
	UpdateFunc      func(ctx context.Context, user *domain.User) error
	UpsertBatchFunc func(ctx context.Context, users []domain.User) error
	GetStatsFunc    func(ctx context.Context) (*domain.StatsUsers, error)
}

func (m *MockUserRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, userID)
	}
	return domain.User{}, storage.ErrNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) UpsertBatch(ctx context.Context, users []domain.User) error {
	if m.UpsertBatchFunc != nil {
		return m.UpsertBatchFunc(ctx, users)
	}
	return nil
}

func (m *MockUserRepository) GetStats(ctx context.Context) (*domain.StatsUsers, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx)
	}
	return &domain.StatsUsers{}, nil
}
