package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

type txKey struct{}

// WithTx adds a transaction to the context.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	GetPRAssignments(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
	Update(ctx context.Context, pr *domain.PullRequest) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByName(ctx context.Context, teamName string) (domain.Team, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, userID string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpsertBatch(ctx context.Context, users []domain.User) error
}

type Tx interface {
	PullRequestRepo() PullRequestRepository
	TeamRepo() TeamRepository
	UserRepo() UserRepository
	Commit() error
	Rollback() error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
