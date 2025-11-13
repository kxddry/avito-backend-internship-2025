package helpers

import (
	"math/rand"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	"github.com/kxddry/avito-backend-internship-2025/pkg/algo"
)

// PickReviewers picks two reviewers from the candidates.
func PickReviewers(candidates []domain.TeamMember, skip algo.Set[string]) []string {
	n := len(candidates)
	if n == 0 {
		return nil
	}

	const maxAttempts = 1024

	first := -1
	second := -1

	for attempts := 0; attempts < maxAttempts && (first == -1 || second == -1); attempts++ {
		i := rand.Intn(n) //nolint:gosec
		c := candidates[i]

		if i == first || !c.IsActive || skip.Has(c.UserID) {
			continue
		}

		if first == -1 {
			first = i
		} else {
			second = i
		}
	}

	switch {
	case first == -1:
		return []string{} // no valid candidates found
	case second == -1:
		return []string{candidates[first].UserID}
	default:
		return []string{candidates[first].UserID, candidates[second].UserID}
	}
}

// ReplaceReviewer replaces a reviewer from the candidates.
func ReplaceReviewer(candidates []domain.TeamMember, skip algo.Set[string]) (out string, ok bool) {
	filtered := make([]domain.TeamMember, 0, len(candidates))
	for _, cand := range candidates {
		if !cand.IsActive || skip.Has(cand.UserID) {
			continue
		}
		filtered = append(filtered, cand)
	}

	if len(filtered) == 0 {
		return "", false
	}

	i := rand.Intn(len(filtered)) //nolint:gosec
	return filtered[i].UserID, true
}
