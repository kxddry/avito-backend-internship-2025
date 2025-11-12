package txmanager

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

type TxManager struct {
	pool  *pgxpool.Pool
	repos *Repositories
}

type Repositories struct {
	// repos here
}

func New(ctx context.Context, dsn string) (*TxManager, error) {
	panic("not implemented")
}

func (m *TxManager) Close() {
	m.pool.Close()
}

type tx struct {
	ctx   context.Context
	repos *Repositories
	tx    pgx.Tx
}

// Rollback rolls back the transaction.
func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// Commit commits the transaction.
func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

// Do executes a function within a transaction.
func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context, tx *tx) error) error {
	pgTx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	t := &tx{
		ctx:   ctx,
		repos: m.repos,
		tx:    pgTx,
	}

	if err := fn(ctx, t); err != nil {
		_ = t.Rollback(ctx)
		return err
	}

	return t.Commit(ctx)
}

// DoWith starts a transaction, injects it into the context, and runs the provided function.
// It commits on success and rolls back on error.
//
//nolint:gocognit
func (m *TxManager) DoWith(ctx context.Context, fn func(ctx context.Context) error) error {
	pgTx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctxWithTx := storage.WithTx(ctx, pgTx)

	if err := fn(ctxWithTx); err != nil {
		_ = pgTx.Rollback(ctx)
		return err
	}

	return pgTx.Commit(ctx)
}
