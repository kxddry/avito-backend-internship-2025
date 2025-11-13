package txmanager

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

type tx struct {
	ctx   context.Context
	repos *Repositories
	tx    pgx.Tx
}

var _ storage.Tx = (*tx)(nil)

func (t *tx) PullRequestRepo() storage.PullRequestRepository {
	return t.repos.PullRequests
}

func (t *tx) TeamRepo() storage.TeamRepository {
	return t.repos.Teams
}

func (t *tx) UserRepo() storage.UserRepository {
	return t.repos.Users
}

func (t *tx) Rollback() error {
	return t.tx.Rollback(t.ctx)
}

func (t *tx) Commit() error {
	return t.tx.Commit(t.ctx)
}

// Do executes fn() within a transaction.
func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
	pgTx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctxWithTx := storage.WithTx(ctx, pgTx)

	t := &tx{
		ctx:   ctxWithTx,
		repos: m.repos,
		tx:    pgTx,
	}

	if err := fn(ctxWithTx, t); err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
}
