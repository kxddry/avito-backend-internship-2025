package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
)

func TestNew(t *testing.T) {
	txmgr := &MockTxManager{}
	service := New(Dependencies{TransactionManager: txmgr})

	if service == nil {
		t.Fatal("New() returned nil")
	}
	if service.txmgr != txmgr {
		t.Error("New() did not set txmgr correctly")
	}
}

func TestNew_PanicsOnNilTxManager(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("New() should panic when TransactionManager is nil")
		}
	}()

	New(Dependencies{TransactionManager: nil})
}

func TestService_CreatePullRequest_Success(t *testing.T) {
	ctx := t.Context()
	input := &domain.CreatePullRequestInput{
		PullRequestID:   "pr-1",
		PullRequestName: "Test PR",
		AuthorID:        "author-1",
	}

	user := domain.User{
		UserID:   "author-1",
		Username: "Author",
		TeamName: "backend",
		IsActive: true,
	}

	team := domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "author-1", Username: "Author", IsActive: true},
			{UserID: "reviewer-1", Username: "Reviewer1", IsActive: true},
			{UserID: "reviewer-2", Username: "Reviewer2", IsActive: true},
		},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return domain.PullRequest{}, storage.ErrNotFound
						},
						CreateFunc: func(ctx context.Context, pr *domain.PullRequest) error {
							now := time.Now()
							pr.CreatedAt = &now
							return nil
						},
					}
				},
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return user, nil
						},
					}
				},
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						GetByNameFunc: func(ctx context.Context, teamName string) (domain.Team, error) {
							return team, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	pr, err := service.CreatePullRequest(ctx, input)

	if err != nil {
		t.Fatalf("CreatePullRequest() error = %v, want nil", err)
	}

	if pr == nil {
		t.Fatal("CreatePullRequest() returned nil PR")
	}

	if pr.ID != input.PullRequestID {
		t.Errorf("PR.ID = %s, want %s", pr.ID, input.PullRequestID)
	}
	if pr.Name != input.PullRequestName {
		t.Errorf("PR.Name = %s, want %s", pr.Name, input.PullRequestName)
	}
	if pr.AuthorID != input.AuthorID {
		t.Errorf("PR.AuthorID = %s, want %s", pr.AuthorID, input.AuthorID)
	}
	if pr.Status != domain.PullRequestStatusOpen {
		t.Errorf("PR.Status = %s, want %s", pr.Status, domain.PullRequestStatusOpen)
	}
}

func TestService_CreatePullRequest_AlreadyExists(t *testing.T) {
	ctx := t.Context()
	input := &domain.CreatePullRequestInput{
		PullRequestID:   "pr-1",
		PullRequestName: "Test PR",
		AuthorID:        "author-1",
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return domain.PullRequest{ID: pullRequestID}, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.CreatePullRequest(ctx, input)

	if err == nil {
		t.Fatal("CreatePullRequest() should return error when PR exists")
	}

	if !errors.Is(err, domain.ErrPRExists) {
		t.Errorf("CreatePullRequest() error = %v, want ErrPRExists", err)
	}
}

func TestService_CreatePullRequest_AuthorNotFound(t *testing.T) {
	ctx := t.Context()
	input := &domain.CreatePullRequestInput{
		PullRequestID:   "pr-1",
		PullRequestName: "Test PR",
		AuthorID:        "nonexistent",
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return domain.PullRequest{}, storage.ErrNotFound
						},
					}
				},
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return domain.User{}, storage.ErrNotFound
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.CreatePullRequest(ctx, input)

	if err == nil {
		t.Fatal("CreatePullRequest() should return error when author not found")
	}

	if !errors.Is(err, domain.ErrResourceNotFound) {
		t.Errorf("CreatePullRequest() error = %v, want ErrResourceNotFound", err)
	}
}

func TestService_MergePullRequest_Success(t *testing.T) {
	ctx := t.Context()
	input := &domain.MergePullRequestInput{
		PullRequestID: "pr-1",
	}

	existingPR := domain.PullRequest{
		ID:                "pr-1",
		Name:              "Test PR",
		AuthorID:          "author-1",
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: []string{"r1", "r2"},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return existingPR, nil
						},
						UpdateFunc: func(ctx context.Context, pr *domain.PullRequest) error {
							return nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	pr, err := service.MergePullRequest(ctx, input)

	if err != nil {
		t.Fatalf("MergePullRequest() error = %v, want nil", err)
	}

	if pr.Status != domain.PullRequestStatusMerged {
		t.Errorf("PR.Status = %s, want %s", pr.Status, domain.PullRequestStatusMerged)
	}

	if pr.MergedAt == nil {
		t.Error("PR.MergedAt should be set")
	}
}

func TestService_MergePullRequest_Idempotent(t *testing.T) {
	ctx := t.Context()
	input := &domain.MergePullRequestInput{
		PullRequestID: "pr-1",
	}

	mergedAt := time.Now()
	existingPR := domain.PullRequest{
		ID:                "pr-1",
		Name:              "Test PR",
		AuthorID:          "author-1",
		Status:            domain.PullRequestStatusMerged,
		AssignedReviewers: []string{"r1", "r2"},
		MergedAt:          &mergedAt,
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return existingPR, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	pr, err := service.MergePullRequest(ctx, input)

	if err != nil {
		t.Fatalf("MergePullRequest() error = %v, want nil", err)
	}

	if pr.Status != domain.PullRequestStatusMerged {
		t.Errorf("PR.Status = %s, want %s", pr.Status, domain.PullRequestStatusMerged)
	}

	if !pr.MergedAt.Equal(mergedAt) {
		t.Error("MergedAt should not change on idempotent call")
	}
}

func TestService_MergePullRequest_NotFound(t *testing.T) {
	ctx := t.Context()
	input := &domain.MergePullRequestInput{
		PullRequestID: "nonexistent",
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return domain.PullRequest{}, storage.ErrNotFound
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.MergePullRequest(ctx, input)

	if err == nil {
		t.Fatal("MergePullRequest() should return error when PR not found")
	}

	if !errors.Is(err, domain.ErrResourceNotFound) {
		t.Errorf("MergePullRequest() error = %v, want ErrResourceNotFound", err)
	}
}

func TestService_ReassignPullRequest_Success(t *testing.T) {
	ctx := t.Context()
	input := &domain.ReassignPullRequestInput{
		PullRequestID: "pr-1",
		OldUserID:     "reviewer-1",
	}

	oldUser := domain.User{
		UserID:   "reviewer-1",
		Username: "Reviewer1",
		TeamName: "backend",
		IsActive: true,
	}

	team := domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "author-1", Username: "Author", IsActive: true},
			{UserID: "reviewer-1", Username: "Reviewer1", IsActive: true},
			{UserID: "reviewer-2", Username: "Reviewer2", IsActive: true},
			{UserID: "reviewer-3", Username: "Reviewer3", IsActive: true},
		},
	}

	existingPR := domain.PullRequest{
		ID:                "pr-1",
		Name:              "Test PR",
		AuthorID:          "author-1",
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: []string{"reviewer-1", "reviewer-2"},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return oldUser, nil
						},
					}
				},
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						GetByNameFunc: func(ctx context.Context, teamName string) (domain.Team, error) {
							return team, nil
						},
					}
				},
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return existingPR, nil
						},
						UpdateFunc: func(ctx context.Context, pr *domain.PullRequest) error {
							return nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	result, err := service.ReassignPullRequest(ctx, input)

	if err != nil {
		t.Fatalf("ReassignPullRequest() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("ReassignPullRequest() returned nil result")
	}

	if result.ReplacedBy == "" {
		t.Error("ReassignPullRequest() ReplacedBy is empty")
	}

	if result.ReplacedBy == "reviewer-1" || result.ReplacedBy == "reviewer-2" {
		t.Errorf("ReassignPullRequest() should not assign to old reviewer or current reviewer")
	}
}

func TestService_ReassignPullRequest_NotAssigned(t *testing.T) {
	ctx := t.Context()
	input := &domain.ReassignPullRequestInput{
		PullRequestID: "pr-1",
		OldUserID:     "not-assigned",
	}

	oldUser := domain.User{
		UserID:   "not-assigned",
		Username: "NotAssigned",
		TeamName: "backend",
		IsActive: true,
	}

	team := domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "author-1", Username: "Author", IsActive: true},
			{UserID: "reviewer-1", Username: "Reviewer1", IsActive: true},
		},
	}

	existingPR := domain.PullRequest{
		ID:                "pr-1",
		Name:              "Test PR",
		AuthorID:          "author-1",
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: []string{"reviewer-1"},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return oldUser, nil
						},
					}
				},
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						GetByNameFunc: func(ctx context.Context, teamName string) (domain.Team, error) {
							return team, nil
						},
					}
				},
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return existingPR, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.ReassignPullRequest(ctx, input)

	if err == nil {
		t.Fatal("ReassignPullRequest() should return error when reviewer not assigned")
	}

	if !errors.Is(err, domain.ErrReviewerMissing) {
		t.Errorf("ReassignPullRequest() error = %v, want ErrReviewerMissing", err)
	}
}

func TestService_ReassignPullRequest_MergedPR(t *testing.T) {
	ctx := t.Context()
	input := &domain.ReassignPullRequestInput{
		PullRequestID: "pr-1",
		OldUserID:     "reviewer-1",
	}

	oldUser := domain.User{
		UserID:   "reviewer-1",
		Username: "Reviewer1",
		TeamName: "backend",
		IsActive: true,
	}

	team := domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "author-1", Username: "Author", IsActive: true},
			{UserID: "reviewer-1", Username: "Reviewer1", IsActive: true},
		},
	}

	mergedAt := time.Now()
	existingPR := domain.PullRequest{
		ID:                "pr-1",
		Name:              "Test PR",
		AuthorID:          "author-1",
		Status:            domain.PullRequestStatusMerged,
		AssignedReviewers: []string{"reviewer-1"},
		MergedAt:          &mergedAt,
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return oldUser, nil
						},
					}
				},
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						GetByNameFunc: func(ctx context.Context, teamName string) (domain.Team, error) {
							return team, nil
						},
					}
				},
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetByIDFunc: func(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
							return existingPR, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.ReassignPullRequest(ctx, input)

	if err == nil {
		t.Fatal("ReassignPullRequest() should return error when PR is merged")
	}

	if !errors.Is(err, domain.ErrReassignOnMerged) {
		t.Errorf("ReassignPullRequest() error = %v, want ErrReassignOnMerged", err)
	}
}

func TestService_CreateTeam_Success(t *testing.T) {
	ctx := t.Context()
	team := &domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "u1", Username: "User1", IsActive: true},
			{UserID: "u2", Username: "User2", IsActive: true},
		},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						CreateFunc: func(ctx context.Context, t *domain.Team) error {
							return nil
						},
					}
				},
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						UpsertBatchFunc: func(ctx context.Context, users []domain.User) error {
							if len(users) != 2 {
								t.Errorf("UpsertBatch called with %d users, want 2", len(users))
							}
							return nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	result, err := service.CreateTeam(ctx, team)

	if err != nil {
		t.Fatalf("CreateTeam() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("CreateTeam() returned nil")
	}

	if result.Name != team.Name {
		t.Errorf("CreateTeam() Name = %s, want %s", result.Name, team.Name)
	}
}

func TestService_CreateTeam_AlreadyExists(t *testing.T) {
	ctx := t.Context()
	team := &domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "u1", Username: "User1", IsActive: true},
		},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						CreateFunc: func(ctx context.Context, t *domain.Team) error {
							return storage.ErrAlreadyExists
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.CreateTeam(ctx, team)

	if err == nil {
		t.Fatal("CreateTeam() should return error when team exists")
	}

	if !errors.Is(err, domain.ErrTeamExists) {
		t.Errorf("CreateTeam() error = %v, want ErrTeamExists", err)
	}
}

func TestService_GetTeam_Success(t *testing.T) {
	ctx := t.Context()
	teamName := "backend"

	expectedTeam := domain.Team{
		Name: teamName,
		Members: []domain.TeamMember{
			{UserID: "u1", Username: "User1", IsActive: true},
		},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				TeamRepoFunc: func() storage.TeamRepository {
					return &MockTeamRepository{
						GetByNameFunc: func(ctx context.Context, name string) (domain.Team, error) {
							return expectedTeam, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	result, err := service.GetTeam(ctx, teamName)

	if err != nil {
		t.Fatalf("GetTeam() error = %v, want nil", err)
	}

	if result.Name != expectedTeam.Name {
		t.Errorf("GetTeam() Name = %s, want %s", result.Name, expectedTeam.Name)
	}
}

func TestService_GetReviewerAssignments_Success(t *testing.T) {
	ctx := t.Context()
	userID := "reviewer-1"

	user := domain.User{
		UserID:   userID,
		Username: "Reviewer1",
		TeamName: "backend",
		IsActive: true,
	}

	prs := []domain.PullRequestShort{
		{ID: "pr-1", Name: "PR 1", AuthorID: "a1", Status: domain.PullRequestStatusOpen},
		{ID: "pr-2", Name: "PR 2", AuthorID: "a2", Status: domain.PullRequestStatusOpen},
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, id string) (domain.User, error) {
							return user, nil
						},
					}
				},
				PullRequestRepoFunc: func() storage.PullRequestRepository {
					return &MockPullRequestRepository{
						GetPRAssignmentsFunc: func(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
							return prs, nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	result, err := service.GetReviewerAssignments(ctx, userID)

	if err != nil {
		t.Fatalf("GetReviewerAssignments() error = %v, want nil", err)
	}

	if result.UserID != userID {
		t.Errorf("GetReviewerAssignments() UserID = %s, want %s", result.UserID, userID)
	}

	if len(result.PullRequests) != 2 {
		t.Errorf("GetReviewerAssignments() returned %d PRs, want 2", len(result.PullRequests))
	}
}

func TestService_GetReviewerAssignments_UserNotFound(t *testing.T) {
	ctx := t.Context()
	userID := "nonexistent"

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, id string) (domain.User, error) {
							return domain.User{}, storage.ErrNotFound
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.GetReviewerAssignments(ctx, userID)

	if err == nil {
		t.Fatal("GetReviewerAssignments() should return error when user not found")
	}

	if !errors.Is(err, domain.ErrResourceNotFound) {
		t.Errorf("GetReviewerAssignments() error = %v, want ErrResourceNotFound", err)
	}
}

func TestService_SetUserIsActive_Success(t *testing.T) {
	ctx := t.Context()
	input := &domain.SetUserIsActiveInput{
		UserID:   "u1",
		IsActive: false,
	}

	user := domain.User{
		UserID:   "u1",
		Username: "User1",
		TeamName: "backend",
		IsActive: true,
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return user, nil
						},
						UpdateFunc: func(ctx context.Context, u *domain.User) error {
							return nil
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	result, err := service.SetUserIsActive(ctx, input)

	if err != nil {
		t.Fatalf("SetUserIsActive() error = %v, want nil", err)
	}

	if result.IsActive != false {
		t.Errorf("SetUserIsActive() IsActive = %v, want false", result.IsActive)
	}
}

func TestService_SetUserIsActive_UserNotFound(t *testing.T) {
	ctx := t.Context()
	input := &domain.SetUserIsActiveInput{
		UserID:   "nonexistent",
		IsActive: false,
	}

	txmgr := &MockTxManager{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context, tx storage.Tx) error) error {
			mockTx := &MockTx{
				UserRepoFunc: func() storage.UserRepository {
					return &MockUserRepository{
						GetByIDFunc: func(ctx context.Context, userID string) (domain.User, error) {
							return domain.User{}, storage.ErrNotFound
						},
					}
				},
			}
			return fn(ctx, mockTx)
		},
	}

	service := New(Dependencies{TransactionManager: txmgr})
	_, err := service.SetUserIsActive(ctx, input)

	if err == nil {
		t.Fatal("SetUserIsActive() should return error when user not found")
	}

	if !errors.Is(err, domain.ErrResourceNotFound) {
		t.Errorf("SetUserIsActive() error = %v, want ErrResourceNotFound", err)
	}
}
