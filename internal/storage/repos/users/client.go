package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

// Repository is the repository for users.
type Repository struct {
	pool *pgxpool.Pool
}

var _ storage.UserRepository = (*Repository)(nil)

// New creates a new repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetByID gets a user by ID.
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

// Update updates a user.
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

// UpsertBatch upserts a batch of users.
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

// GetStats returns user statistics.
func (r *Repository) GetStats(ctx context.Context) (*domain.StatsUsers, error) {
	q := r.getQuerier(ctx)

	rows, err := q.Query(ctx, getUserStatsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := &domain.StatsUsers{
		ByUser: make([]domain.StatsUserEntry, 0),
	}

	for rows.Next() {
		var entry domain.StatsUserEntry
		if err := rows.Scan(&entry.UserID, &entry.UserName, &entry.Team, &entry.IsActive,
			&entry.AssignedReviewsTotal, &entry.OpenReviews, &entry.MergedReviews); err != nil {
			return nil, err
		}
		stats.ByUser = append(stats.ByUser, entry)

		if entry.IsActive {
			stats.Active++
		} else {
			stats.Inactive++
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	stats.Total = stats.Active + stats.Inactive
	return stats, nil
}
