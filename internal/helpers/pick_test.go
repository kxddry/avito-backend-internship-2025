package helpers

import (
	"math/rand"
	"testing"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/pkg/algo"
)

func TestPickReviewers_EmptyCandidates(t *testing.T) {
	candidates := []domain.TeamMember{}
	skip := algo.SetFrom[string]()

	result := PickReviewers(candidates, skip)

	if result != nil {
		t.Errorf("PickReviewers() with empty candidates = %v, want nil", result)
	}
}

func TestPickReviewers_AllInactive(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: false},
		{UserID: "u2", Username: "User2", IsActive: false},
		{UserID: "u3", Username: "User3", IsActive: false},
	}
	skip := algo.SetFrom[string]()

	result := PickReviewers(candidates, skip)

	if len(result) != 0 {
		t.Errorf("PickReviewers() with all inactive = %v, want empty slice", result)
	}
}

func TestPickReviewers_AllSkipped(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
	}
	skip := algo.SetFrom("u1", "u2")

	result := PickReviewers(candidates, skip)

	if len(result) != 0 {
		t.Errorf("PickReviewers() with all skipped = %v, want empty slice", result)
	}
}

func TestPickReviewers_OneValidCandidate(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: false},
		{UserID: "u3", Username: "User3", IsActive: false},
	}
	skip := algo.SetFrom[string]()

	// Run multiple times to ensure consistency due to randomness
	for i := range 10 {
		result := PickReviewers(candidates, skip)

		if len(result) != 1 {
			t.Errorf("PickReviewers() iteration %d returned %d reviewers, want 1", i, len(result))
			continue
		}

		if result[0] != "u1" {
			t.Errorf("PickReviewers() iteration %d = %v, want [u1]", i, result)
		}
	}
}

func TestPickReviewers_TwoValidCandidates(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
	}
	skip := algo.SetFrom[string]()

	// Run multiple times to check randomness
	for i := range 10 {
		result := PickReviewers(candidates, skip)

		if len(result) != 2 {
			t.Errorf("PickReviewers() iteration %d returned %d reviewers, want 2", i, len(result))
			continue
		}

		// Check that both are from valid candidates
		validIDs := map[string]bool{"u1": true, "u2": true}
		for _, id := range result {
			if !validIDs[id] {
				t.Errorf("PickReviewers() iteration %d returned invalid ID %s", i, id)
			}
		}

		// Check no duplicates
		if result[0] == result[1] {
			t.Errorf("PickReviewers() iteration %d returned duplicate reviewers: %v", i, result)
		}
	}
}

func TestPickReviewers_MultipleValidCandidates(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
		{UserID: "u3", Username: "User3", IsActive: true},
		{UserID: "u4", Username: "User4", IsActive: true},
	}
	skip := algo.SetFrom[string]()

	for i := range 20 {
		result := PickReviewers(candidates, skip)

		if len(result) != 2 {
			t.Errorf("PickReviewers() iteration %d returned %d reviewers, want 2", i, len(result))
			continue
		}

		// Check all are valid
		validIDs := map[string]bool{"u1": true, "u2": true, "u3": true, "u4": true}
		for _, id := range result {
			if !validIDs[id] {
				t.Errorf("PickReviewers() iteration %d returned invalid ID %s", i, id)
			}
		}

		// Check no duplicates
		if result[0] == result[1] {
			t.Errorf("PickReviewers() iteration %d returned duplicate reviewers: %v", i, result)
		}
	}
}

func TestPickReviewers_WithSkip(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "author", Username: "Author", IsActive: true},
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
	}
	skip := algo.SetFrom("author")

	for i := range 10 {
		result := PickReviewers(candidates, skip)

		if len(result) != 2 {
			t.Errorf("PickReviewers() iteration %d returned %d reviewers, want 2", i, len(result))
			continue
		}

		// Check author is not in result
		for _, id := range result {
			if id == "author" {
				t.Errorf("PickReviewers() iteration %d included skipped author", i)
			}
		}

		// Check no duplicates
		if result[0] == result[1] {
			t.Errorf("PickReviewers() iteration %d returned duplicate reviewers: %v", i, result)
		}
	}
}

func TestPickReviewers_MixedActiveInactive(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: false},
		{UserID: "u3", Username: "User3", IsActive: true},
		{UserID: "u4", Username: "User4", IsActive: false},
	}
	skip := algo.SetFrom[string]()

	for i := range 10 {
		result := PickReviewers(candidates, skip)

		if len(result) != 2 {
			t.Errorf("PickReviewers() iteration %d returned %d reviewers, want 2", i, len(result))
			continue
		}

		// Check all are active
		for _, id := range result {
			if id != "u1" && id != "u3" {
				t.Errorf("PickReviewers() iteration %d returned inactive user %s", i, id)
			}
		}

		// Check no duplicates
		if result[0] == result[1] {
			t.Errorf("PickReviewers() iteration %d returned duplicate reviewers: %v", i, result)
		}
	}
}

func TestReplaceReviewer_NoValidCandidates(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: false},
		{UserID: "u2", Username: "User2", IsActive: false},
	}
	skip := algo.SetFrom[string]()

	result, ok := ReplaceReviewer(candidates, skip)

	if ok {
		t.Errorf("ReplaceReviewer() with no valid candidates returned ok=true")
	}
	if result != "" {
		t.Errorf("ReplaceReviewer() with no valid candidates = %q, want empty string", result)
	}
}

func TestReplaceReviewer_AllSkipped(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
	}
	skip := algo.SetFrom("u1", "u2")

	result, ok := ReplaceReviewer(candidates, skip)

	if ok {
		t.Errorf("ReplaceReviewer() with all skipped returned ok=true")
	}
	if result != "" {
		t.Errorf("ReplaceReviewer() with all skipped = %q, want empty string", result)
	}
}

func TestReplaceReviewer_OneValidCandidate(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: false},
	}
	skip := algo.SetFrom[string]()

	for i := range 10 {
		result, ok := ReplaceReviewer(candidates, skip)

		if !ok {
			t.Errorf("ReplaceReviewer() iteration %d returned ok=false, want true", i)
		}
		if result != "u1" {
			t.Errorf("ReplaceReviewer() iteration %d = %q, want u1", i, result)
		}
	}
}

func TestReplaceReviewer_MultipleValidCandidates(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
		{UserID: "u3", Username: "User3", IsActive: true},
	}
	skip := algo.SetFrom[string]()

	validIDs := map[string]bool{"u1": true, "u2": true, "u3": true}

	for i := range 20 {
		result, ok := ReplaceReviewer(candidates, skip)

		if !ok {
			t.Errorf("ReplaceReviewer() iteration %d returned ok=false, want true", i)
		}
		if !validIDs[result] {
			t.Errorf("ReplaceReviewer() iteration %d = %q, want one of u1, u2, u3", i, result)
		}
	}
}

func TestReplaceReviewer_WithSkip(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "old", Username: "Old", IsActive: true},
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
	}
	skip := algo.SetFrom("old")

	for i := range 20 {
		result, ok := ReplaceReviewer(candidates, skip)

		if !ok {
			t.Errorf("ReplaceReviewer() iteration %d returned ok=false, want true", i)
		}
		if result == "old" {
			t.Errorf("ReplaceReviewer() iteration %d = %q, should not return skipped user", i, result)
		}
		if result != "u1" && result != "u2" {
			t.Errorf("ReplaceReviewer() iteration %d = %q, want u1 or u2", i, result)
		}
	}
}

func TestReplaceReviewer_InactiveSkipped(t *testing.T) {
	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: false},
		{UserID: "u2", Username: "User2", IsActive: true},
		{UserID: "u3", Username: "User3", IsActive: false},
	}
	skip := algo.SetFrom[string]()

	for i := range 10 {
		result, ok := ReplaceReviewer(candidates, skip)

		if !ok {
			t.Errorf("ReplaceReviewer() iteration %d returned ok=false, want true", i)
		}
		if result != "u2" {
			t.Errorf("ReplaceReviewer() iteration %d = %q, want u2", i, result)
		}
	}
}

// TestPickReviewers_Randomness verifies that the function produces different results.
func TestPickReviewers_Randomness(t *testing.T) {
	rand.Seed(1)

	candidates := []domain.TeamMember{
		{UserID: "u1", Username: "User1", IsActive: true},
		{UserID: "u2", Username: "User2", IsActive: true},
		{UserID: "u3", Username: "User3", IsActive: true},
		{UserID: "u4", Username: "User4", IsActive: true},
		{UserID: "u5", Username: "User5", IsActive: true},
	}
	skip := algo.SetFrom[string]()

	results := make(map[string]bool)
	for range 50 {
		result := PickReviewers(candidates, skip)
		key := result[0] + "," + result[1]
		results[key] = true
	}

	// We should get at least some variety (not all the same pair)
	if len(results) < 2 {
		t.Errorf("PickReviewers() seems not random enough, only got %d different pairs out of 50 runs", len(results))
	}
}
