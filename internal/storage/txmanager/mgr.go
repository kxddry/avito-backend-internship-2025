package txmanager

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
	prrepo "github.com/kxddry/avito-backend-internship-2025/internal/storage/repos/pullrequests"
	teamsrepo "github.com/kxddry/avito-backend-internship-2025/internal/storage/repos/teams"
	usersrepo "github.com/kxddry/avito-backend-internship-2025/internal/storage/repos/users"
)

// Repositories is a collection of repositories.
type Repositories struct {
	PullRequests *prrepo.Repository
	Teams        *teamsrepo.Repository
	Users        *usersrepo.Repository
}

// TxManager is a transaction manager.
type TxManager struct {
	pool  *pgxpool.Pool
	repos *Repositories
}

var _ storage.TxManager = (*TxManager)(nil)

// New creates a new TxManager.
func New(ctx context.Context, dsn string) (*TxManager, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	repos := &Repositories{
		PullRequests: prrepo.New(pool),
		Teams:        teamsrepo.New(pool),
		Users:        usersrepo.New(pool),
	}

	return &TxManager{
		pool:  pool,
		repos: repos,
	}, nil
}

// Close closes the database connection pool.
func (m *TxManager) Close() {
	m.pool.Close()
}
