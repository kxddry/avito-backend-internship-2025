package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

type Repository struct {
	pool *pgxpool.Pool
}

var _ storage.UserRepository = (*Repository)(nil)

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	q := r.getQuerier(ctx)
	var (
		id       string
		name     string
		isActive bool
		team     *string
	)

	err := q.QueryRow(ctx, getByIDQuery, userID).Scan(&id, &name, &isActive, &team)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, storage.ErrNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		UserID:   id,
		Username: name,
		IsActive: isActive,
		TeamName: derefString(team),
	}, nil
}

func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	q := r.getQuerier(ctx)
	tag, err := q.Exec(ctx, updateUserQuery, user.UserID, user.Username, user.IsActive, nullableString(user.TeamName))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (r *Repository) UpsertBatch(ctx context.Context, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}

	q := r.getQuerier(ctx)
	for _, user := range users {
		if _, err := q.Exec(ctx, upsertUserQuery, user.UserID, user.Username, user.IsActive, nullableString(user.TeamName)); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) getQuerier(ctx context.Context) storage.Querier {
	if tx, ok := storage.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
