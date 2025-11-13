package api //nolint:testpackage

import (
	"reflect"
	"testing"
	"time"

	"github.com/kxddry/avito-backend-internship-2025/internal/api/generated"
	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

func TestToGeneratedPullRequest(t *testing.T) {
	now := time.Now()
	mergedAt := now.Add(time.Hour)

	tests := []struct {
		name string
		pr   *domain.PullRequest
		want generated.PullRequest
	}{
		{
			name: "full pull request",
			pr: &domain.PullRequest{
				ID:                "pr-1",
				Name:              "Test PR",
				AuthorID:          "author-1",
				Status:            domain.PullRequestStatusOpen,
				AssignedReviewers: []string{"reviewer-1", "reviewer-2"},
				CreatedAt:         &now,
				MergedAt:          &mergedAt,
			},
			want: generated.PullRequest{
				PullRequestId:     "pr-1",
				PullRequestName:   "Test PR",
				AuthorId:          "author-1",
				Status:            generated.PullRequestStatusOPEN,
				AssignedReviewers: []string{"reviewer-1", "reviewer-2"},
				CreatedAt:         &now,
				MergedAt:          &mergedAt,
			},
		},
		{
			name: "pull request without reviewers",
			pr: &domain.PullRequest{
				ID:                "pr-2",
				Name:              "Empty PR",
				AuthorID:          "author-2",
				Status:            domain.PullRequestStatusMerged,
				AssignedReviewers: []string{},
				CreatedAt:         &now,
			},
			want: generated.PullRequest{
				PullRequestId:     "pr-2",
				PullRequestName:   "Empty PR",
				AuthorId:          "author-2",
				Status:            generated.PullRequestStatusMERGED,
				AssignedReviewers: []string{},
				CreatedAt:         &now,
			},
		},
		{
			name: "pull request with nil times",
			pr: &domain.PullRequest{
				ID:                "pr-3",
				Name:              "Nil Times PR",
				AuthorID:          "author-3",
				Status:            domain.PullRequestStatusOpen,
				AssignedReviewers: []string{"reviewer-1"},
			},
			want: generated.PullRequest{
				PullRequestId:     "pr-3",
				PullRequestName:   "Nil Times PR",
				AuthorId:          "author-3",
				Status:            generated.PullRequestStatusOPEN,
				AssignedReviewers: []string{"reviewer-1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toGeneratedPullRequest(tt.pr)

			if got.PullRequestId != tt.want.PullRequestId {
				t.Errorf("PullRequestId = %v, want %v", got.PullRequestId, tt.want.PullRequestId)
			}
			if got.PullRequestName != tt.want.PullRequestName {
				t.Errorf("PullRequestName = %v, want %v", got.PullRequestName, tt.want.PullRequestName)
			}
			if got.AuthorId != tt.want.AuthorId {
				t.Errorf("AuthorId = %v, want %v", got.AuthorId, tt.want.AuthorId)
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
			if !reflect.DeepEqual(got.AssignedReviewers, tt.want.AssignedReviewers) {
				t.Errorf("AssignedReviewers = %v, want %v", got.AssignedReviewers, tt.want.AssignedReviewers)
			}
			if !equalTimePointers(got.CreatedAt, tt.want.CreatedAt) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}
			if !equalTimePointers(got.MergedAt, tt.want.MergedAt) {
				t.Errorf("MergedAt = %v, want %v", got.MergedAt, tt.want.MergedAt)
			}
		})
	}
}

func TestToGeneratedPullRequest_ReviewersCopy(t *testing.T) {
	// Test that reviewers are copied, not referenced
	reviewers := []string{"r1", "r2"}
	pr := &domain.PullRequest{
		ID:                "pr-1",
		Name:              "Test",
		AuthorID:          "a1",
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: reviewers,
	}

	got := toGeneratedPullRequest(pr)

	// Modify original
	reviewers[0] = "modified"

	// Check that generated is not affected
	if got.AssignedReviewers[0] == "modified" {
		t.Error("toGeneratedPullRequest should copy reviewers, not reference them")
	}
}

func TestToGeneratedPullRequestShort(t *testing.T) {
	tests := []struct {
		name string
		pr   domain.PullRequestShort
		want generated.PullRequestShort
	}{
		{
			name: "open PR",
			pr: domain.PullRequestShort{
				ID:       "pr-1",
				Name:     "Test PR",
				AuthorID: "author-1",
				Status:   domain.PullRequestStatusOpen,
			},
			want: generated.PullRequestShort{
				PullRequestId:   "pr-1",
				PullRequestName: "Test PR",
				AuthorId:        "author-1",
				Status:          generated.PullRequestShortStatusOPEN,
			},
		},
		{
			name: "merged PR",
			pr: domain.PullRequestShort{
				ID:       "pr-2",
				Name:     "Merged PR",
				AuthorID: "author-2",
				Status:   domain.PullRequestStatusMerged,
			},
			want: generated.PullRequestShort{
				PullRequestId:   "pr-2",
				PullRequestName: "Merged PR",
				AuthorId:        "author-2",
				Status:          generated.PullRequestShortStatusMERGED,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toGeneratedPullRequestShort(tt.pr)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toGeneratedPullRequestShort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToGeneratedPullRequestShortList(t *testing.T) {
	input := []domain.PullRequestShort{
		{ID: "pr-1", Name: "PR 1", AuthorID: "a1", Status: domain.PullRequestStatusOpen},
		{ID: "pr-2", Name: "PR 2", AuthorID: "a2", Status: domain.PullRequestStatusMerged},
	}

	got := toGeneratedPullRequestShortList(input)

	if len(got) != 2 {
		t.Fatalf("toGeneratedPullRequestShortList() length = %d, want 2", len(got))
	}

	if got[0].PullRequestId != "pr-1" {
		t.Errorf("toGeneratedPullRequestShortList()[0].PullRequestId = %s, want pr-1", got[0].PullRequestId)
	}
	if got[1].PullRequestId != "pr-2" {
		t.Errorf("toGeneratedPullRequestShortList()[1].PullRequestId = %s, want pr-2", got[1].PullRequestId)
	}
}

func TestToGeneratedTeam(t *testing.T) {
	team := &domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: false},
		},
	}

	got := toGeneratedTeam(team)

	if got.TeamName != "backend" {
		t.Errorf("TeamName = %s, want backend", got.TeamName)
	}

	if len(got.Members) != 2 {
		t.Fatalf("Members length = %d, want 2", len(got.Members))
	}

	if got.Members[0].UserId != "u1" {
		t.Errorf("Members[0].UserId = %s, want u1", got.Members[0].UserId)
	}
	if got.Members[0].Username != "Alice" {
		t.Errorf("Members[0].Username = %s, want Alice", got.Members[0].Username)
	}
	if !got.Members[0].IsActive {
		t.Error("Members[0].IsActive = false, want true")
	}

	if got.Members[1].UserId != "u2" {
		t.Errorf("Members[1].UserId = %s, want u2", got.Members[1].UserId)
	}
	if got.Members[1].IsActive {
		t.Error("Members[1].IsActive = true, want false")
	}
}

func TestFromGeneratedTeam(t *testing.T) {
	genTeam := generated.Team{
		TeamName: "frontend",
		Members: []generated.TeamMember{
			{UserId: "u1", Username: "Charlie", IsActive: true},
			{UserId: "u2", Username: "David", IsActive: false},
		},
	}

	got := fromGeneratedTeam(genTeam)

	if got.Name != "frontend" {
		t.Errorf("Name = %s, want frontend", got.Name)
	}

	if len(got.Members) != 2 {
		t.Fatalf("Members length = %d, want 2", len(got.Members))
	}

	if got.Members[0].UserID != "u1" {
		t.Errorf("Members[0].UserID = %s, want u1", got.Members[0].UserID)
	}
	if got.Members[0].Username != "Charlie" {
		t.Errorf("Members[0].Username = %s, want Charlie", got.Members[0].Username)
	}
	if !got.Members[0].IsActive {
		t.Error("Members[0].IsActive = false, want true")
	}

	if got.Members[1].UserID != "u2" {
		t.Errorf("Members[1].UserID = %s, want u2", got.Members[1].UserID)
	}
	if got.Members[1].IsActive {
		t.Error("Members[1].IsActive = true, want false")
	}
}

func TestToGeneratedUser(t *testing.T) {
	user := &domain.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	got := toGeneratedUser(user)

	if got.UserId != "u1" {
		t.Errorf("UserId = %s, want u1", got.UserId)
	}
	if got.Username != "Alice" {
		t.Errorf("Username = %s, want Alice", got.Username)
	}
	if got.TeamName != "backend" {
		t.Errorf("TeamName = %s, want backend", got.TeamName)
	}
	if !got.IsActive {
		t.Error("IsActive = false, want true")
	}
}

func TestRoundTrip_Team(t *testing.T) {
	// Test that fromGenerated(toGenerated(x)) == x
	original := &domain.Team{
		Name: "test-team",
		Members: []domain.TeamMember{
			{UserID: "u1", Username: "User1", IsActive: true},
			{UserID: "u2", Username: "User2", IsActive: false},
		},
	}

	gen := toGeneratedTeam(original)
	roundTrip := fromGeneratedTeam(gen)

	if roundTrip.Name != original.Name {
		t.Errorf("Round trip Name = %s, want %s", roundTrip.Name, original.Name)
	}

	if len(roundTrip.Members) != len(original.Members) {
		t.Fatalf("Round trip Members length = %d, want %d", len(roundTrip.Members), len(original.Members))
	}

	for i := range original.Members {
		if !reflect.DeepEqual(roundTrip.Members[i], original.Members[i]) {
			t.Errorf("Round trip Members[%d] = %+v, want %+v", i, roundTrip.Members[i], original.Members[i])
		}
	}
}

func TestUnwrapDomainError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr *domain.Error
		wantNil bool
	}{
		{
			name:    "nil error",
			err:     nil,
			wantErr: nil,
			wantNil: true,
		},
		{
			name:    "domain error",
			err:     domain.ErrResourceNotFound,
			wantErr: domain.ErrResourceNotFound,
			wantNil: true,
		},
		{
			name: "wrapped domain error",
			err: domain.NewError(
				500,
				"TEST",
				"test",
				domain.ErrResourceNotFound,
			),
			wantErr: domain.NewError(500, "TEST", "test", domain.ErrResourceNotFound),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotNilErr := unwrapDomainError(tt.err)

			if tt.wantNil && gotNilErr != nil {
				t.Errorf("unwrapDomainError() error = %v, want nil", gotNilErr)
			}

			if tt.wantErr == nil && gotErr != nil {
				t.Errorf("unwrapDomainError() gotErr = %v, want nil", gotErr)
			}

			if tt.wantErr != nil && gotErr == nil {
				t.Error("unwrapDomainError() gotErr = nil, want non-nil")
			}

			if tt.wantErr != nil && gotErr != nil {
				if gotErr.Code != tt.wantErr.Code {
					t.Errorf("unwrapDomainError() Code = %s, want %s", gotErr.Code, tt.wantErr.Code)
				}
			}
		})
	}
}

// Helper function to compare time pointers.
func equalTimePointers(a, b *time.Time) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equal(*b)
}
