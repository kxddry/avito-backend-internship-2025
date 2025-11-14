package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

// txKey is the key for the transaction.
type txKey struct{}

// WithTx adds a transaction to the context.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx gets the transaction from the context.
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// PullRequestRepository is the repository for pull requests.
type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.CreatePullRequestInput) (domain.PullRequest, error)
	GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	GetPRAssignments(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
	Update(ctx context.Context, pr *domain.PullRequest) error
	GetStats(ctx context.Context) (*domain.StatsPRs, error)
}

// TeamRepository is the repository for teams.
type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByName(ctx context.Context, teamName string) (domain.Team, error)
	GetStats(ctx context.Context) (*domain.StatsTeams, error)
}

// UserRepository is the repository for users.
type UserRepository interface {
	GetByID(ctx context.Context, userID string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpsertBatch(ctx context.Context, users []domain.User) error
	GetStats(ctx context.Context) (*domain.StatsUsers, error)
}

// Tx is the transaction.
type Tx interface {
	PullRequestRepo() PullRequestRepository
	TeamRepo() TeamRepository
	UserRepo() UserRepository
	Commit() error
	Rollback() error
}

// TxManager is the transaction manager.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

// Querier is the querier.
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
