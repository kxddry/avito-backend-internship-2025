package service

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/helpers"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
	"github.com/kxddry/avito-backend-internship-2025/pkg/algo"
)

// Dependencies are the dependencies for the service.
type Dependencies struct {
	TransactionManager storage.TxManager
}

// Service is the service for assigning pull requests to reviewers.
type Service struct {
	txmgr storage.TxManager
}

// New creates a new service.
func New(deps Dependencies) *Service {
	txmgr := deps.TransactionManager
	if txmgr == nil {
		panic("New Service: deps.TransactionManager is nil")
	}

	return &Service{
		txmgr: txmgr,
	}
}

var _ domain.AssignmentService = (*Service)(nil)

// formatError formats an error.
func (s *Service) formatError(ctx context.Context, op string, err error) error {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return domain.ErrResourceNotFound
	case domain.IsDomainError(err):
		return err
	case errors.Is(err, ctx.Err()):
		return ctx.Err()
	default:
		log.Error().Err(err).Str("operation", op).Msg("operation failed")
		return domain.ErrInternal
	}
}

// CreatePullRequest creates a new pull request.
func (s *Service) CreatePullRequest(outerCtx context.Context, input *domain.CreatePullRequestInput) (*domain.PullRequest, error) { //nolint:lll
	const op = "service.CreatePullRequest"
	var pr *domain.PullRequest
	err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		prRepo := tx.PullRequestRepo()
		_, err := prRepo.GetByID(ctx, input.PullRequestID)
		if err == nil {
			return storage.ErrAlreadyExists
		}
		if !errors.Is(err, storage.ErrNotFound) {
			return err
		}

		user, err := tx.UserRepo().GetByID(ctx, input.AuthorID)
		if err != nil {
			return err
		}

		team, err := tx.TeamRepo().GetByName(ctx, user.TeamName)
		if err != nil {
			return err
		}

		reviewers := helpers.PickReviewers(team.Members, algo.SetFrom(input.AuthorID))

		pr = &domain.PullRequest{
			ID:                input.PullRequestID,
			Name:              input.PullRequestName,
			AuthorID:          input.AuthorID,
			AssignedReviewers: reviewers,
			Status:            domain.PullRequestStatusOpen,
		}

		return prRepo.Create(ctx, pr)
	})

	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, domain.ErrPRExists
		}
		return nil, s.formatError(outerCtx, op, err)
	}
	return pr, nil
}

// MergePullRequest merges a pull request.
func (s *Service) MergePullRequest(outerCtx context.Context, input *domain.MergePullRequestInput) (*domain.PullRequest, error) { //nolint:lll
	const op = "service.MergePullRequest"
	var pr *domain.PullRequest
	err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		prRepo := tx.PullRequestRepo()
		pr2, err := prRepo.GetByID(ctx, input.PullRequestID)
		if err != nil {
			return err
		}
		if pr2.Status == domain.PullRequestStatusMerged {
			pr = &pr2
			return nil
		}
		now := time.Now().UTC()
		pr2.MergedAt = &now
		pr2.Status = domain.PullRequestStatusMerged

		err = prRepo.Update(ctx, &pr2)
		if err != nil {
			return err
		}
		pr = &pr2
		return nil
	})

	if err != nil {
		return nil, s.formatError(outerCtx, op, err)
	}
	return pr, nil
}

// ReassignPullRequest reassigns a pull request.
func (s *Service) ReassignPullRequest(outerCtx context.Context, input *domain.ReassignPullRequestInput) (*domain.ReassignPullRequestResult, error) { //nolint:lll
	const op = "service.ReassignPullRequest"
	var result *domain.ReassignPullRequestResult
	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		oldUser, err := tx.UserRepo().GetByID(ctx, input.OldUserID)
		if err != nil {
			return err
		}

		team, err := tx.TeamRepo().GetByName(ctx, oldUser.TeamName)
		if err != nil {
			return err
		}

		prRepo := tx.PullRequestRepo()
		pr, err := prRepo.GetByID(ctx, input.PullRequestID)
		if err != nil {
			return err
		}

		if !slices.Contains(pr.AssignedReviewers, oldUser.UserID) {
			return domain.ErrReviewerMissing
		}

		if pr.Status == domain.PullRequestStatusMerged {
			return domain.ErrReassignOnMerged
		}

		excludeSet := algo.SetFrom(pr.AssignedReviewers...)
		excludeSet.Add(pr.AuthorID)
		newReviewerID, ok := helpers.ReplaceReviewer(team.Members, excludeSet)
		if !ok {
			return domain.ErrNoCandidate
		}

		if !algo.ReplaceOnce(pr.AssignedReviewers, oldUser.UserID, newReviewerID) {
			return domain.ErrInternal
		}

		err = prRepo.Update(ctx, &pr)
		if err != nil {
			return err
		}
		result = &domain.ReassignPullRequestResult{
			PullRequest: pr,
			ReplacedBy:  newReviewerID,
		}
		return nil

	}); err != nil {
		return nil, s.formatError(outerCtx, op, err)
	}
	return result, nil
}

// transformMembersToUsers transforms a list of team members to a list of users.
func transformMembersToUsers(teamName string, members []domain.TeamMember) []domain.User {
	users := make([]domain.User, len(members))
	for i, member := range members {
		users[i] = domain.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: teamName,
			IsActive: member.IsActive,
		}
	}
	return users
}

// CreateTeam creates a new team.
func (s *Service) CreateTeam(outerCtx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "service.CreateTeam"
	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		if err := tx.TeamRepo().Create(ctx, team); err != nil {
			return err
		}
		if err := tx.UserRepo().UpsertBatch(ctx, transformMembersToUsers(team.Name, team.Members)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, domain.ErrTeamExists
		}
		return nil, s.formatError(outerCtx, op, err)
	}

	return team, nil
}

// GetTeam gets a team by name.
func (s *Service) GetTeam(outerCtx context.Context, teamName string) (*domain.Team, error) {
	const op = "service.GetTeam"
	var result *domain.Team
	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		team, err := tx.TeamRepo().GetByName(ctx, teamName)
		if err != nil {
			return err
		}
		result = &team
		return nil
	}); err != nil {
		return nil, s.formatError(outerCtx, op, err)
	}
	return result, nil
}

// GetReviewerAssignments gets the reviewer assignments for a user.
func (s *Service) GetReviewerAssignments(outerCtx context.Context, userID string) (*domain.ReviewerAssignments, error) {
	const op = "service.GetReviewerAssignments"
	result := &domain.ReviewerAssignments{UserID: userID, PullRequests: []domain.PullRequestShort{}}

	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) (err error) {
		_, err = tx.UserRepo().GetByID(ctx, userID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return domain.ErrResourceNotFound
			}
			return err
		}

		result.PullRequests, err = tx.PullRequestRepo().GetPRAssignments(ctx, userID)
		if errors.Is(err, storage.ErrNotFound) {
			result.PullRequests = []domain.PullRequestShort{}
			return nil
		}
		return err
	}); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return result, nil
		}
		return nil, s.formatError(outerCtx, op, err)
	}
	return result, nil
}

// SetUserIsActive sets the active status of a user.
func (s *Service) SetUserIsActive(outerCtx context.Context, input *domain.SetUserIsActiveInput) (*domain.User, error) {
	const op = "service.SetUserIsActive"
	var user *domain.User
	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		userRepo := tx.UserRepo()
		user2, err := userRepo.GetByID(ctx, input.UserID)
		if err != nil {
			return err
		}

		user2.IsActive = input.IsActive
		if err := userRepo.Update(ctx, &user2); err != nil {
			return err
		}
		user = &user2
		return nil
	}); err != nil {
		return nil, s.formatError(outerCtx, op, err)
	}
	return user, nil
}

// DeactivateTeam deactivates all users in a team.
func (s *Service) DeactivateTeam(outerCtx context.Context, teamName string) (int, error) {
	const op = "service.DeactivateTeam"
	var count int
	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		team, err := tx.TeamRepo().GetByName(ctx, teamName)
		if err != nil {
			return err
		}

		users := transformMembersToUsers(team.Name, team.Members)
		for i := range users {
			users[i].IsActive = false
		}

		if err := tx.UserRepo().UpsertBatch(ctx, users); err != nil {
			return err
		}

		count = len(users)
		return nil
	}); err != nil {
		return 0, s.formatError(outerCtx, op, err)
	}
	return count, nil
}

// SafeReassignPR safely reassigns inactive reviewers on an open PR.
func (s *Service) SafeReassignPR(outerCtx context.Context, prID string) (*domain.PullRequest, error) {
	const op = "service.SafeReassignPR"
	var pr *domain.PullRequest
	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		prRepo := tx.PullRequestRepo()
		pr2, err := prRepo.GetByID(ctx, prID)
		if err != nil {
			return err
		}

		if pr2.Status != domain.PullRequestStatusOpen {
			pr = &pr2
			return nil
		}

		author, err := tx.UserRepo().GetByID(ctx, pr2.AuthorID)
		if err != nil {
			return err
		}

		team, err := tx.TeamRepo().GetByName(ctx, author.TeamName)
		if err != nil {
			return err
		}

		needUpdate := false
		newReviewers := make([]string, 0, len(pr2.AssignedReviewers))

		for _, reviewerID := range pr2.AssignedReviewers {
			reviewer, err := tx.UserRepo().GetByID(ctx, reviewerID)
			if err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					needUpdate = true
					continue
				}
				return err
			}

			if !reviewer.IsActive {
				needUpdate = true
				excludeSet := algo.SetFrom(pr2.AssignedReviewers...)
				excludeSet.Add(pr2.AuthorID)
				if newReviewerID, ok := helpers.ReplaceReviewer(team.Members, excludeSet); ok {
					newReviewers = append(newReviewers, newReviewerID)
					excludeSet.Add(newReviewerID)
				}
			} else {
				newReviewers = append(newReviewers, reviewerID)
			}
		}

		if needUpdate {
			pr2.AssignedReviewers = newReviewers
			if err := prRepo.Update(ctx, &pr2); err != nil {
				return err
			}
		}

		pr = &pr2
		return nil
	}); err != nil {
		return nil, s.formatError(outerCtx, op, err)
	}
	return pr, nil
}

// GetStats returns statistics about the system.
func (s *Service) GetStats(outerCtx context.Context) (*domain.Stats, error) {
	const op = "service.GetStats"
	stats := &domain.Stats{}

	if err := s.txmgr.Do(outerCtx, func(ctx context.Context, tx storage.Tx) error {
		userStats, err := tx.UserRepo().GetStats(ctx)
		if err != nil {
			return err
		}
		stats.Users = *userStats

		prStats, err := tx.PullRequestRepo().GetStats(ctx)
		if err != nil {
			return err
		}
		stats.PRs = *prStats

		teamStats, err := tx.TeamRepo().GetStats(ctx)
		if err != nil {
			return err
		}
		stats.Teams = *teamStats

		return nil
	}); err != nil {
		return nil, s.formatError(outerCtx, op, err)
	}

	return stats, nil
}
