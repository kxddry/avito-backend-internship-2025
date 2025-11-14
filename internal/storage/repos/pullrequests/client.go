package pullrequests

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

// Repository is the repository for pull requests.
type Repository struct {
	pool *pgxpool.Pool
}

var _ storage.PullRequestRepository = (*Repository)(nil)

// New creates a new repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new pull request.
func (r *Repository) Create(ctx context.Context, pr *domain.PullRequest) error {
	q := r.getQuerier(ctx)

	var (
		createdAt time.Time
		mergedAt  *time.Time
	)

	if err := q.QueryRow(ctx, createPRQuery, pr.ID, pr.Name, pr.AuthorID, string(pr.Status),
		pr.AssignedReviewers, pr.MergedAt,
	).Scan(&createdAt, &mergedAt); err != nil {
		if storage.IsUniqueViolation(err) {
			return storage.ErrAlreadyExists
		}
		return err
	}

	pr.CreatedAt = &createdAt
	pr.MergedAt = mergedAt
	return nil
}

// GetByID gets a pull request by ID.
func (r *Repository) GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	q := r.getQuerier(ctx)

	var (
		id        string
		name      string
		authorID  string
		status    string
		reviewers []string
		createdAt time.Time
		mergedAt  *time.Time
	)

	err := q.QueryRow(ctx, getPRByIDQuery, pullRequestID).
		Scan(&id, &name, &authorID, &status, &reviewers, &createdAt, &mergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PullRequest{}, storage.ErrNotFound
		}
		return domain.PullRequest{}, err
	}

	return domain.PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            domain.PullRequestStatus(status),
		AssignedReviewers: reviewers,
		CreatedAt:         &createdAt,
		MergedAt:          mergedAt,
	}, nil
}

// GetPRAssignments gets the pull request assignments for a reviewer.
func (r *Repository) GetPRAssignments(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	q := r.getQuerier(ctx)

	rows, err := q.Query(ctx, getAssignmentsQuery, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.PullRequestShort, 0)
	for rows.Next() {
		var (
			id       string
			name     string
			authorID string
			status   string
		)
		if err := rows.Scan(&id, &name, &authorID, &status); err != nil {
			return nil, err
		}
		result = append(result, domain.PullRequestShort{
			ID:       id,
			Name:     name,
			AuthorID: authorID,
			Status:   domain.PullRequestStatus(status),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, storage.ErrNotFound
	}
	return result, nil
}

// Update updates a pull request.
func (r *Repository) Update(ctx context.Context, pr *domain.PullRequest) error {
	q := r.getQuerier(ctx)

	var (
		createdAt time.Time
		mergedAt  *time.Time
	)
	err := q.QueryRow(
		ctx,
		updatePRQuery,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		string(pr.Status),
		pr.AssignedReviewers,
		pr.MergedAt,
	).Scan(&createdAt, &mergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	pr.CreatedAt = &createdAt
	pr.MergedAt = mergedAt
	return nil
}

// getQuerier gets the querier.
func (r *Repository) getQuerier(ctx context.Context) storage.Querier {
	if tx, ok := storage.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

// GetStats returns PR statistics.
func (r *Repository) GetStats(ctx context.Context) (*domain.StatsPRs, error) {
	q := r.getQuerier(ctx)

	var stats domain.StatsPRs
	err := q.QueryRow(ctx, getPRStatsQuery).Scan(&stats.Total, &stats.Open, &stats.Merged, &stats.With0Reviewers, &stats.With1Reviewer, &stats.With2Reviewers)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
