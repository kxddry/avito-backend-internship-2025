//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kxddry/avito-backend-internship-2025/internal/api"
	"github.com/kxddry/avito-backend-internship-2025/internal/api/generated"
	"github.com/kxddry/avito-backend-internship-2025/internal/service"
	"github.com/kxddry/avito-backend-internship-2025/internal/storage/txmanager"
	"github.com/kxddry/avito-backend-internship-2025/pkg/logging"
)

var (
	testServer *httptest.Server
	testClient *http.Client
	baseURL    string
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	logging.SetupLogger(false)

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	}

	txmgr, err := txmanager.New(ctx, dsn)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create transaction manager: %v\n", err)
		os.Exit(1)
	}
	defer txmgr.Close()

	svc := service.New(service.Dependencies{
		TransactionManager: txmgr,
	})

	apiServer := api.NewServer(svc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	strictHandler := generated.NewStrictHandler(apiServer, nil)
	generated.RegisterHandlers(e, strictHandler)

	testServer = httptest.NewServer(e)
	baseURL = testServer.URL
	testClient = testServer.Client()

	code := m.Run()

	testServer.Close()
	os.Exit(code)
}

func cleanupTestData(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, "TRUNCATE TABLE pull_requests, users, teams CASCADE")
	if err != nil {
		t.Fatalf("Failed to cleanup test data: %v", err)
	}
}

func TestE2E_TeamFlow(t *testing.T) {
	cleanupTestData(t)

	t.Run("CreateOld Team", func(t *testing.T) {
		team := generated.Team{
			TeamName: "backend-team",
			Members: []generated.TeamMember{
				{UserId: "user1", Username: "Alice", IsActive: true},
				{UserId: "user2", Username: "Bob", IsActive: true},
				{UserId: "user3", Username: "Charlie", IsActive: false},
			},
		}

		body, _ := json.Marshal(team)
		resp, err := testClient.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("CreateOld team status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}

		var result struct {
			Team generated.Team `json:"team"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Team.TeamName != "backend-team" {
			t.Errorf("Team name = %s, want backend-team", result.Team.TeamName)
		}
		if len(result.Team.Members) != 3 {
			t.Errorf("Team members count = %d, want 3", len(result.Team.Members))
		}
	})

	t.Run("Get Team", func(t *testing.T) {
		resp, err := testClient.Get(baseURL + "/team/get?team_name=backend-team")
		if err != nil {
			t.Fatalf("Failed to get team: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Get team status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var team generated.Team
		if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if team.TeamName != "backend-team" {
			t.Errorf("Team name = %s, want backend-team", team.TeamName)
		}
	})

	t.Run("CreateOld Duplicate Team", func(t *testing.T) {
		team := generated.Team{
			TeamName: "backend-team",
			Members:  []generated.TeamMember{},
		}

		body, _ := json.Marshal(team)
		resp, err := testClient.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("CreateOld duplicate team status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})
}

func TestE2E_PullRequestFlow(t *testing.T) {
	cleanupTestData(t)

	team := generated.Team{
		TeamName: "dev-team",
		Members: []generated.TeamMember{
			{UserId: "author", Username: "Author", IsActive: true},
			{UserId: "reviewer1", Username: "Reviewer1", IsActive: true},
			{UserId: "reviewer2", Username: "Reviewer2", IsActive: true},
			{UserId: "reviewer3", Username: "Reviewer3", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	resp, err := testClient.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	_ = resp.Body.Close()

	t.Run("CreateOld Pull Request", func(t *testing.T) {
		prInput := map[string]string{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Add feature X",
			"author_id":         "author",
		}

		body, _ := json.Marshal(prInput)
		resp, err := testClient.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("CreateOld PR status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}

		var result struct {
			Pr generated.PullRequest `json:"pr"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Pr.PullRequestId != "pr-1" {
			t.Errorf("PR ID = %s, want pr-1", result.Pr.PullRequestId)
		}
		if result.Pr.Status != generated.PullRequestStatusOPEN {
			t.Errorf("PR Status = %s, want OPEN", result.Pr.Status)
		}
		if len(result.Pr.AssignedReviewers) == 0 {
			t.Error("PR should have assigned reviewers")
		}

		for _, reviewer := range result.Pr.AssignedReviewers {
			if reviewer == "author" {
				t.Error("Author should not be assigned as reviewer")
			}
		}
	})

	t.Run("Get Reviewer Assignments", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)

		resp, err := testClient.Get(baseURL + "/users/getReview?user_id=reviewer1")
		if err != nil {
			t.Fatalf("Failed to get assignments: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Get assignments status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var result struct {
			UserId       string                       `json:"user_id"`
			PullRequests []generated.PullRequestShort `json:"pull_requests"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.UserId != "reviewer1" {
			t.Errorf("User ID = %s, want reviewer1", result.UserId)
		}
	})

	t.Run("Merge Pull Request", func(t *testing.T) {
		mergeInput := map[string]string{
			"pull_request_id": "pr-1",
		}

		body, _ := json.Marshal(mergeInput)
		resp, err := testClient.Post(baseURL+"/pullRequest/merge", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to merge PR: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Merge PR status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var result struct {
			Pr generated.PullRequest `json:"pr"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Pr.Status != generated.PullRequestStatusMERGED {
			t.Errorf("PR Status = %s, want MERGED", result.Pr.Status)
		}
		if result.Pr.MergedAt == nil {
			t.Error("PR MergedAt should be set")
		}
	})

	t.Run("Merge Idempotency", func(t *testing.T) {
		mergeInput := map[string]string{
			"pull_request_id": "pr-1",
		}

		body, _ := json.Marshal(mergeInput)
		resp, err := testClient.Post(baseURL+"/pullRequest/merge", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to merge PR: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Merge PR (idempotent) status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("Reassign on Merged PR Should Fail", func(t *testing.T) {
		reassignInput := map[string]string{
			"pull_request_id": "pr-1",
			"old_user_id":     "reviewer1",
		}

		body, _ := json.Marshal(reassignInput)
		resp, err := testClient.Post(baseURL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to reassign PR: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Reassign merged PR status = %d, want %d", resp.StatusCode, http.StatusConflict)
		}
	})
}

func TestE2E_ReassignFlow(t *testing.T) {
	cleanupTestData(t)

	team := generated.Team{
		TeamName: "reassign-team",
		Members: []generated.TeamMember{
			{UserId: "author", Username: "Author", IsActive: true},
			{UserId: "rev1", Username: "Reviewer1", IsActive: true},
			{UserId: "rev2", Username: "Reviewer2", IsActive: true},
			{UserId: "rev3", Username: "Reviewer3", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	resp, err := testClient.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	_ = resp.Body.Close()

	prInput := map[string]string{
		"pull_request_id":   "pr-reassign",
		"pull_request_name": "Test Reassign",
		"author_id":         "author",
	}

	body, _ = json.Marshal(prInput)
	resp, err = testClient.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create PR: %v", err)
	}

	var createResult struct {
		Pr generated.PullRequest `json:"pr"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResult); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	_ = resp.Body.Close()

	if len(createResult.Pr.AssignedReviewers) == 0 {
		t.Skip("No reviewers assigned, cannot test reassignment")
	}

	oldReviewer := createResult.Pr.AssignedReviewers[0]

	t.Run("Reassign Reviewer", func(t *testing.T) {
		reassignInput := map[string]string{
			"pull_request_id": "pr-reassign",
			"old_user_id":     oldReviewer,
		}

		body, _ := json.Marshal(reassignInput)
		resp, err := testClient.Post(baseURL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to reassign: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Reassign status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var result struct {
			Pr         generated.PullRequest `json:"pr"`
			ReplacedBy string                `json:"replaced_by"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.ReplacedBy == "" {
			t.Error("ReplacedBy should not be empty")
		}

		if result.ReplacedBy == oldReviewer {
			t.Error("New reviewer should not be the same as old reviewer")
		}

		found := false
		for _, reviewer := range result.Pr.AssignedReviewers {
			if reviewer == oldReviewer {
				found = true
				break
			}
		}
		if found {
			t.Error("Old reviewer should not be in assigned reviewers")
		}
	})
}

func TestE2E_SetUserIsActive(t *testing.T) {
	cleanupTestData(t)

	team := generated.Team{
		TeamName: "active-team",
		Members: []generated.TeamMember{
			{UserId: "test-user", Username: "TestUser", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	resp, err := testClient.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	_ = resp.Body.Close()

	t.Run("Deactivate User", func(t *testing.T) {
		input := map[string]interface{}{
			"user_id":   "test-user",
			"is_active": false,
		}

		body, _ := json.Marshal(input)
		resp, err := testClient.Post(baseURL+"/users/setIsActive", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to set user active: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Set user active status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var result struct {
			User generated.User `json:"user"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.User.IsActive {
			t.Error("User should be inactive")
		}
	})

	t.Run("Reactivate User", func(t *testing.T) {
		input := map[string]interface{}{
			"user_id":   "test-user",
			"is_active": true,
		}

		body, _ := json.Marshal(input)
		resp, err := testClient.Post(baseURL+"/users/setIsActive", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to set user active: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Set user active status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var result struct {
			User generated.User `json:"user"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !result.User.IsActive {
			t.Error("User should be active")
		}
	})
}

func TestE2E_ErrorCases(t *testing.T) {
	cleanupTestData(t)

	t.Run("Get Non-Existing Team", func(t *testing.T) {
		resp, err := testClient.Get(baseURL + "/team/get?team_name=non-existing")
		if err != nil {
			t.Fatalf("Failed to get team: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Get non-existing team status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("CreateOld PR with Non-Existing Author", func(t *testing.T) {
		prInput := map[string]string{
			"pull_request_id":   "pr-bad",
			"pull_request_name": "Bad PR",
			"author_id":         "non-existing-author",
		}

		body, _ := json.Marshal(prInput)
		resp, err := testClient.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("CreateOld PR with bad author status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("Set Active for Non-Existing User", func(t *testing.T) {
		input := map[string]interface{}{
			"user_id":   "non-existing-user",
			"is_active": false,
		}

		body, _ := json.Marshal(input)
		resp, err := testClient.Post(baseURL+"/users/setIsActive", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to set user active: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Set active for non-existing user status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})
}
