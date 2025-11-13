package api //nolint:testpackage

import (
	"net/http"
	"testing"

	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

func TestMakeErrorPayload(t *testing.T) {
	tests := []struct {
		name        string
		err         *domain.Error
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{
			name:        "error with message",
			err:         domain.ErrResourceNotFound,
			wantStatus:  http.StatusNotFound,
			wantCode:    string(domain.ErrorCodeNotFound),
			wantMessage: string(domain.ErrorMessageResourceNotFound),
		},
		{
			name:        "error with different message",
			err:         domain.ErrPRExists,
			wantStatus:  http.StatusConflict,
			wantCode:    string(domain.ErrorCodePullRequestExists),
			wantMessage: string(domain.ErrorMessagePullRequestExists),
		},
		{
			name: "error without message uses code",
			err: &domain.Error{
				Status: http.StatusBadRequest,
				Code:   "TEST_CODE",
			},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "TEST_CODE",
			wantMessage: "TEST_CODE",
		},
		{
			name:        "team exists error",
			err:         domain.ErrTeamExists,
			wantStatus:  http.StatusBadRequest,
			wantCode:    string(domain.ErrorCodeTeamExists),
			wantMessage: string(domain.ErrorMessageTeamExists),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := makeErrorPayload(tt.err)

			if payload.status != tt.wantStatus {
				t.Errorf("makeErrorPayload().status = %d, want %d", payload.status, tt.wantStatus)
			}

			if string(payload.body.Error.Code) != tt.wantCode {
				t.Errorf("makeErrorPayload().body.Error.Code = %s, want %s",
					payload.body.Error.Code, tt.wantCode)
			}

			if payload.body.Error.Message != tt.wantMessage {
				t.Errorf("makeErrorPayload().body.Error.Message = %s, want %s",
					payload.body.Error.Message, tt.wantMessage)
			}
		})
	}
}

func TestMakeErrorPayload_AllPredefinedErrors(t *testing.T) {
	errors := []*domain.Error{
		domain.ErrPRExists,
		domain.ErrResourceNotFound,
		domain.ErrTeamExists,
		domain.ErrInternal,
		domain.ErrReassignOnMerged,
		domain.ErrReviewerMissing,
		domain.ErrNoCandidate,
	}

	for _, err := range errors {
		t.Run(string(err.Code), func(t *testing.T) {
			payload := makeErrorPayload(err)

			// Check that status is set
			if payload.status == 0 {
				t.Error("makeErrorPayload().status is 0")
			}

			// Check that code is set
			if payload.body.Error.Code == "" {
				t.Error("makeErrorPayload().body.Error.Code is empty")
			}

			// Check that message is set
			if payload.body.Error.Message == "" {
				t.Error("makeErrorPayload().body.Error.Message is empty")
			}

			// Verify code matches
			if string(payload.body.Error.Code) != string(err.Code) {
				t.Errorf("makeErrorPayload().body.Error.Code = %s, want %s",
					payload.body.Error.Code, err.Code)
			}
		})
	}
}

func TestErrorPayload_Structure(t *testing.T) {
	err := domain.ErrResourceNotFound
	payload := makeErrorPayload(err)

	// Verify the structure is correct for JSON serialization
	if payload.body.Error.Code == "" {
		t.Error("Error.Code should not be empty")
	}

	if payload.body.Error.Message == "" {
		t.Error("Error.Message should not be empty")
	}
}
