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

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	GetPRAssignments(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) // must return {} if none found instead of nil
	Update(ctx context.Context, pr *domain.PullRequest) error
	Delete(ctx context.Context, pullRequestID string) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByName(ctx context.Context, teamName string) (domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, teamName string) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, userID string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpsertBatch(ctx context.Context, users []domain.User) error
	Delete(ctx context.Context, userID string) error
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

type Dependencies struct {
	TransactionManager TxManager
}

type Service struct {
	txmgr TxManager
}

func New(deps Dependencies) *Service {
	txmgr := deps.TransactionManager
	if txmgr == nil {
		panic("разраб ты че решил мне нил давать?")
	}

	return &Service{
		txmgr: txmgr,
	}
}

var _ domain.AssignmentService = (*Service)(nil)

func (s *Service) formatError(op string, err error) error {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return domain.ErrorResourceNotFound
	case domain.IsDomainError(err):
		return err
	default:
		log.Error().Err(err).Str("operation", op).Msg("operation failed")
		return domain.ErrorInternal
	}
}

func (s *Service) CreatePullRequest(ctx context.Context, input *domain.CreatePullRequestInput) (*domain.PullRequest, error) {
	const op = "service.CreatePullRequest"
	var pr *domain.PullRequest
	err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
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
			return nil, domain.ErrorPullRequestExists
		}
		return nil, s.formatError(op, err)
	}
	return pr, nil
}

func (s *Service) MergePullRequest(ctx context.Context, input *domain.MergePullRequestInput) (*domain.PullRequest, error) {
	const op = "service.MergePullRequest"
	var pr *domain.PullRequest
	err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
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
		return nil, s.formatError(op, err)
	}
	return pr, nil
}

func (s *Service) ReassignPullRequest(ctx context.Context, input *domain.ReassignPullRequestInput) (*domain.ReassignPullRequestResult, error) {
	const op = "service.ReassignPullRequest"
	var result *domain.ReassignPullRequestResult
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
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
			return domain.ErrorReviewerIsNotAssigned
		}

		if pr.Status == domain.PullRequestStatusMerged {
			return domain.ErrorCantReassignOnMergedPr
		}

		excludeSet := algo.SetFrom(pr.AssignedReviewers...)
		excludeSet.Add(pr.AuthorID)
		newReviewerID, ok := helpers.ReplaceReviewer(team.Members, excludeSet)
		if !ok {
			return domain.ErrorNoCandidate
		}

		if !algo.ReplaceOnce(pr.AssignedReviewers, oldUser.UserID, newReviewerID) {
			return domain.ErrorInternal
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
		return nil, s.formatError(op, err)
	}
	return result, nil
}

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

func (s *Service) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "service.CreateTeam"
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		if err := tx.TeamRepo().Create(ctx, team); err != nil {
			return err
		}
		if err := tx.UserRepo().UpsertBatch(ctx, transformMembersToUsers(team.Name, team.Members)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, domain.ErrorTeamExists
		}
		return nil, s.formatError(op, err)
	}

	return team, nil
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	const op = "service.GetTeam"
	var result *domain.Team
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		team, err := tx.TeamRepo().GetByName(ctx, teamName)
		if err != nil {
			return err
		}
		result = &team
		return nil
	}); err != nil {
		return nil, s.formatError(op, err)
	}
	return result, nil
}

func (s *Service) GetReviewerAssignments(ctx context.Context, userID string) (*domain.ReviewerAssignments, error) {
	const op = "service.GetReviewerAssignments"
	result := &domain.ReviewerAssignments{UserID: userID, PullRequests: []domain.PullRequestShort{}}

	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) (err error) {
		_, err = tx.UserRepo().GetByID(ctx, userID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return domain.ErrorResourceNotFound
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
		return nil, s.formatError(op, err)
	}
	return result, nil
}

func (s *Service) SetUserIsActive(ctx context.Context, input *domain.SetUserIsActiveInput) (*domain.User, error) {
	const op = "service.SetUserIsActive"
	var user *domain.User
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
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
		return nil, s.formatError(op, err)
	}
	return user, nil
}
