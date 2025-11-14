//go:build integration
// +build integration

package repos

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage"
	prrepo "github.com/kxddry/avito-backend-internship-2025/internal/storage/repos/pullrequests"
	teamsrepo "github.com/kxddry/avito-backend-internship-2025/internal/storage/repos/teams"
	usersrepo "github.com/kxddry/avito-backend-internship-2025/internal/storage/repos/users"
)

var (
	testPool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Get database URL from environment or use default test database
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	}

	var err error
	testPool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	// Ping database to ensure connection
	if err := testPool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		testPool.Close()
		os.Exit(1)
	}

	// Run migrations if needed
	// In real scenario, you would run migrations here

	code := m.Run()

	testPool.Close()
	os.Exit(code)
}

func setupTestData(t *testing.T) {
	ctx := context.Background()

	// Clean up all tables before each test
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE pull_requests, users, teams CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
}

func TestTeamRepository_Integration(t *testing.T) {
	setupTestData(t)
	ctx := context.Background()

	repo := teamsrepo.New(testPool)

	t.Run("Create and Get Team", func(t *testing.T) {
		team := &domain.Team{
			Name: "test-team",
			Members: []domain.TeamMember{
				{UserID: "u1", Username: "User1", IsActive: true},
				{UserID: "u2", Username: "User2", IsActive: false},
			},
		}

		// Create team
		err := repo.Create(ctx, team)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Get team
		retrieved, err := repo.GetByName(ctx, "test-team")
		if err != nil {
			t.Fatalf("GetByName() error = %v", err)
		}

		if retrieved.Name != team.Name {
			t.Errorf("GetByName() Name = %s, want %s", retrieved.Name, team.Name)
		}
	})

	t.Run("Create Duplicate Team", func(t *testing.T) {
		team := &domain.Team{
			Name:    "duplicate-team",
			Members: []domain.TeamMember{},
		}

		err := repo.Create(ctx, team)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		err = repo.Create(ctx, team)
		if err == nil {
			t.Error("Create() should return error for duplicate team")
		}

		if err != storage.ErrAlreadyExists {
			t.Errorf("Create() error = %v, want ErrAlreadyExists", err)
		}
	})

	t.Run("Get Non-Existing Team", func(t *testing.T) {
		_, err := repo.GetByName(ctx, "non-existing")
		if err == nil {
			t.Error("GetByName() should return error for non-existing team")
		}

		if err != storage.ErrNotFound {
			t.Errorf("GetByName() error = %v, want ErrNotFound", err)
		}
	})
}

func TestUserRepository_Integration(t *testing.T) {
	setupTestData(t)
	ctx := context.Background()

	teamRepo := teamsrepo.New(testPool)
	userRepo := usersrepo.New(testPool)

	// Create a team first
	team := &domain.Team{
		Name:    "user-test-team",
		Members: []domain.TeamMember{},
	}
	if err := teamRepo.Create(ctx, team); err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}

	t.Run("UpsertBatch and GetByID", func(t *testing.T) {
		users := []domain.User{
			{UserID: "u1", Username: "User1", TeamName: "user-test-team", IsActive: true},
			{UserID: "u2", Username: "User2", TeamName: "user-test-team", IsActive: false},
		}

		err := userRepo.UpsertBatch(ctx, users)
		if err != nil {
			t.Fatalf("UpsertBatch() error = %v", err)
		}

		// Get first user
		user, err := userRepo.GetByID(ctx, "u1")
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if user.UserID != "u1" {
			t.Errorf("GetByID() UserID = %s, want u1", user.UserID)
		}
		if user.Username != "User1" {
			t.Errorf("GetByID() Username = %s, want User1", user.Username)
		}
		if !user.IsActive {
			t.Error("GetByID() IsActive = false, want true")
		}
	})

	t.Run("Update User", func(t *testing.T) {
		user := domain.User{
			UserID:   "u3",
			Username: "User3",
			TeamName: "user-test-team",
			IsActive: true,
		}

		err := userRepo.UpsertBatch(ctx, []domain.User{user})
		if err != nil {
			t.Fatalf("UpsertBatch() error = %v", err)
		}

		// Update user
		user.IsActive = false
		err = userRepo.Update(ctx, &user)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Verify update
		updated, err := userRepo.GetByID(ctx, "u3")
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if updated.IsActive {
			t.Error("Update() did not update IsActive")
		}
	})

	t.Run("Update Non-Existing User", func(t *testing.T) {
		user := domain.User{
			UserID:   "non-existing",
			Username: "NonExisting",
			TeamName: "user-test-team",
			IsActive: true,
		}

		err := userRepo.Update(ctx, &user)
		if err == nil {
			t.Error("Update() should return error for non-existing user")
		}

		if err != storage.ErrNotFound {
			t.Errorf("Update() error = %v, want ErrNotFound", err)
		}
	})
}

func TestPullRequestRepository_Integration(t *testing.T) {
	setupTestData(t)
	ctx := context.Background()

	teamRepo := teamsrepo.New(testPool)
	userRepo := usersrepo.New(testPool)
	prRepo := prrepo.New(testPool)

	// Setup test data
	team := &domain.Team{
		Name:    "pr-test-team",
		Members: []domain.TeamMember{},
	}
	if err := teamRepo.Create(ctx, team); err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}

	users := []domain.User{
		{UserID: "author1", Username: "Author1", TeamName: "pr-test-team", IsActive: true},
		{UserID: "reviewer1", Username: "Reviewer1", TeamName: "pr-test-team", IsActive: true},
		{UserID: "reviewer2", Username: "Reviewer2", TeamName: "pr-test-team", IsActive: true},
	}
	if err := userRepo.UpsertBatch(ctx, users); err != nil {
		t.Fatalf("Failed to create users: %v", err)
	}

	t.Run("CreateOld and Get PR", func(t *testing.T) {
		pr := &domain.PullRequest{
			ID:                "pr-1",
			Name:              "Test PR",
			AuthorID:          "author1",
			Status:            domain.PullRequestStatusOpen,
			AssignedReviewers: []string{"reviewer1", "reviewer2"},
		}

		err := prRepo.CreateOld(ctx, pr)
		if err != nil {
			t.Fatalf("CreateOld() error = %v", err)
		}

		// Verify CreatedAt was set
		if pr.CreatedAt == nil {
			t.Error("CreateOld() did not set CreatedAt")
		}

		// Get PR
		retrieved, err := prRepo.GetByID(ctx, "pr-1")
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if retrieved.ID != pr.ID {
			t.Errorf("GetByID() ID = %s, want %s", retrieved.ID, pr.ID)
		}
		if retrieved.Name != pr.Name {
			t.Errorf("GetByID() Name = %s, want %s", retrieved.Name, pr.Name)
		}
		if retrieved.Status != pr.Status {
			t.Errorf("GetByID() Status = %s, want %s", retrieved.Status, pr.Status)
		}
		if len(retrieved.AssignedReviewers) != len(pr.AssignedReviewers) {
			t.Errorf("GetByID() AssignedReviewers length = %d, want %d",
				len(retrieved.AssignedReviewers), len(pr.AssignedReviewers))
		}
	})

	t.Run("Update PR", func(t *testing.T) {
		pr := &domain.PullRequest{
			ID:                "pr-2",
			Name:              "Test PR 2",
			AuthorID:          "author1",
			Status:            domain.PullRequestStatusOpen,
			AssignedReviewers: []string{"reviewer1"},
		}

		err := prRepo.CreateOld(ctx, pr)
		if err != nil {
			t.Fatalf("CreateOld() error = %v", err)
		}

		// Update to merged
		now := time.Now()
		pr.Status = domain.PullRequestStatusMerged
		pr.MergedAt = &now

		err = prRepo.Update(ctx, pr)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Verify update
		retrieved, err := prRepo.GetByID(ctx, "pr-2")
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if retrieved.Status != domain.PullRequestStatusMerged {
			t.Errorf("Update() Status = %s, want %s", retrieved.Status, domain.PullRequestStatusMerged)
		}
		if retrieved.MergedAt == nil {
			t.Error("Update() MergedAt is nil")
		}
	})

	t.Run("GetPRAssignments", func(t *testing.T) {
		// Create multiple PRs assigned to reviewer1
		prs := []*domain.PullRequest{
			{
				ID:                "pr-3",
				Name:              "PR 3",
				AuthorID:          "author1",
				Status:            domain.PullRequestStatusOpen,
				AssignedReviewers: []string{"reviewer1"},
			},
			{
				ID:                "pr-4",
				Name:              "PR 4",
				AuthorID:          "author1",
				Status:            domain.PullRequestStatusOpen,
				AssignedReviewers: []string{"reviewer1", "reviewer2"},
			},
		}

		for _, pr := range prs {
			if err := prRepo.Create(ctx, pr); err != nil {
				t.Fatalf("Create() error = %v", err)
			if err := prRepo.CreateOld(ctx, pr); err != nil {
				t.Fatalf("CreateOld() error = %v", err)
			}
		}

		// Get assignments for reviewer1
		assignments, err := prRepo.GetPRAssignments(ctx, "reviewer1")
		if err != nil {
			t.Fatalf("GetPRAssignments() error = %v", err)
		}

		if len(assignments) < 2 {
			t.Errorf("GetPRAssignments() returned %d PRs, want at least 2", len(assignments))
		}
	})

	t.Run("CreateOld Duplicate PR", func(t *testing.T) {
		pr := &domain.PullRequest{
			ID:                "pr-duplicate",
			Name:              "Duplicate PR",
			AuthorID:          "author1",
			Status:            domain.PullRequestStatusOpen,
			AssignedReviewers: []string{},
		}

		err := prRepo.CreateOld(ctx, pr)
		if err != nil {
			t.Fatalf("CreateOld() error = %v", err)
		}

		err = prRepo.CreateOld(ctx, pr)
		if err == nil {
			t.Error("CreateOld() should return error for duplicate PR")
		}

		if err != storage.ErrAlreadyExists {
			t.Errorf("CreateOld() error = %v, want ErrAlreadyExists", err)
		}
	})
}
