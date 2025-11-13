package teams

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

var _ storage.TeamRepository = (*Repository)(nil)

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, team *domain.Team) error {
	q := r.getQuerier(ctx)
	_, err := q.Exec(ctx, createTeamQuery, team.Name)
	if err != nil {
		if storage.IsUniqueViolation(err) {
			return storage.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *Repository) GetByName(ctx context.Context, teamName string) (domain.Team, error) {
	q := r.getQuerier(ctx)
	var name string
	if err := q.QueryRow(ctx, getTeamQuery, teamName).Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, storage.ErrNotFound
		}
		return domain.Team{}, err
	}

	rows, err := q.Query(ctx, getTeamMembersQuery, teamName)
	if err != nil {
		return domain.Team{}, err
	}
	defer rows.Close()

	members := make([]domain.TeamMember, 0)
	for rows.Next() {
		var (
			id       string
			username string
			isActive bool
		)
		if err := rows.Scan(&id, &username, &isActive); err != nil {
			return domain.Team{}, err
		}
		members = append(members, domain.TeamMember{
			UserID:   id,
			Username: username,
			IsActive: isActive,
		})
	}

	if err := rows.Err(); err != nil {
		return domain.Team{}, err
	}

	return domain.Team{
		Name:    name,
		Members: members,
	}, nil
}

func (r *Repository) getQuerier(ctx context.Context) storage.Querier {
	if tx, ok := storage.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}
