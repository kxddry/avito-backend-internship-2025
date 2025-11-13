package domain

import (
	"errors"
	"net/http"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "error with Err field",
			err: &Error{
				Status:  http.StatusInternalServerError,
				Code:    "TEST_CODE",
				Message: "test message",
				Err:     errors.New("underlying error"),
			},
			want: "underlying error",
		},
		{
			name: "error with Message field",
			err: &Error{
				Status:  http.StatusBadRequest,
				Code:    "TEST_CODE",
				Message: "test message",
			},
			want: "test message",
		},
		{
			name: "error with only Code",
			err: &Error{
				Status: http.StatusNotFound,
				Code:   "NOT_FOUND",
			},
			want: "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")

	err := &Error{
		Status:  http.StatusInternalServerError,
		Code:    "TEST_CODE",
		Message: "test message",
		Err:     underlyingErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != underlyingErr {
		t.Errorf("Error.Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestError_Unwrap_Nil(t *testing.T) {
	err := &Error{
		Status:  http.StatusBadRequest,
		Code:    "TEST_CODE",
		Message: "test message",
	}

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Errorf("Error.Unwrap() = %v, want nil", unwrapped)
	}
}

func TestNewError(t *testing.T) {
	underlyingErr := errors.New("underlying")
	err := NewError(http.StatusBadRequest, "TEST_CODE", "test message", underlyingErr)

	if err.Status != http.StatusBadRequest {
		t.Errorf("NewError().Status = %d, want %d", err.Status, http.StatusBadRequest)
	}
	if err.Code != "TEST_CODE" {
		t.Errorf("NewError().Code = %s, want TEST_CODE", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("NewError().Message = %s, want test message", err.Message)
	}
	if err.Err != underlyingErr {
		t.Errorf("NewError().Err = %v, want %v", err.Err, underlyingErr)
	}
}

func TestIsDomainError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is domain error",
			err:  ErrResourceNotFound,
			want: true,
		},
		{
			name: "is wrapped domain error",
			err:  NewError(http.StatusBadRequest, "TEST", "test", nil),
			want: true,
		},
		{
			name: "is not domain error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDomainError(tt.err); got != tt.want {
				t.Errorf("IsDomainError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         *Error
		wantStatus  int
		wantCode    ErrorCode
		wantMessage string
	}{
		{
			name:        "ErrPRExists",
			err:         ErrPRExists,
			wantStatus:  http.StatusConflict,
			wantCode:    ErrorCodePullRequestExists,
			wantMessage: string(ErrorMessagePullRequestExists),
		},
		{
			name:        "ErrResourceNotFound",
			err:         ErrResourceNotFound,
			wantStatus:  http.StatusNotFound,
			wantCode:    ErrorCodeNotFound,
			wantMessage: string(ErrorMessageResourceNotFound),
		},
		{
			name:        "ErrTeamExists",
			err:         ErrTeamExists,
			wantStatus:  http.StatusBadRequest,
			wantCode:    ErrorCodeTeamExists,
			wantMessage: string(ErrorMessageTeamExists),
		},
		{
			name:        "ErrInternal",
			err:         ErrInternal,
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_ERROR",
			wantMessage: "internal server error",
		},
		{
			name:        "ErrReassignOnMerged",
			err:         ErrReassignOnMerged,
			wantStatus:  http.StatusConflict,
			wantCode:    ErrorCodePullRequestMerged,
			wantMessage: string(ErrorMessageReassignOnMerged),
		},
		{
			name:        "ErrReviewerMissing",
			err:         ErrReviewerMissing,
			wantStatus:  http.StatusConflict,
			wantCode:    ErrorCodeReviewerMissing,
			wantMessage: string(ErrorMessageReviewerIsNotAssigned),
		},
		{
			name:        "ErrNoCandidate",
			err:         ErrNoCandidate,
			wantStatus:  http.StatusConflict,
			wantCode:    ErrorCodeNoCandidate,
			wantMessage: string(ErrorMessageNoActiveReplacementCandidate),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Status != tt.wantStatus {
				t.Errorf("%s.Status = %d, want %d", tt.name, tt.err.Status, tt.wantStatus)
			}
			if tt.err.Code != tt.wantCode {
				t.Errorf("%s.Code = %s, want %s", tt.name, tt.err.Code, tt.wantCode)
			}
			if tt.err.Message != tt.wantMessage {
				t.Errorf("%s.Message = %s, want %s", tt.name, tt.err.Message, tt.wantMessage)
			}
		})
	}
}

func TestErrorCodes(t *testing.T) {
	codes := []ErrorCode{
		ErrorCodeTeamExists,
		ErrorCodePullRequestExists,
		ErrorCodePullRequestMerged,
		ErrorCodeReviewerMissing,
		ErrorCodeNoCandidate,
		ErrorCodeNotFound,
	}

	expectedCodes := []string{
		"TEAM_EXISTS",
		"PR_EXISTS",
		"PR_MERGED",
		"NOT_ASSIGNED",
		"NO_CANDIDATE",
		"NOT_FOUND",
	}

	for i, code := range codes {
		if string(code) != expectedCodes[i] {
			t.Errorf("ErrorCode[%d] = %s, want %s", i, code, expectedCodes[i])
		}
	}
}

func TestErrorMessages(t *testing.T) {
	messages := []ErrorMessage{
		ErrorMessagePullRequestExists,
		ErrorMessageResourceNotFound,
		ErrorMessageTeamExists,
		ErrorMessageReassignOnMerged,
		ErrorMessageReviewerIsNotAssigned,
		ErrorMessageNoActiveReplacementCandidate,
	}

	// Just check that messages are not empty
	for i, msg := range messages {
		if string(msg) == "" {
			t.Errorf("ErrorMessage[%d] is empty", i)
		}
	}
}

func TestError_WithErrors_As(t *testing.T) {
	// Test that errors.As works correctly with our Error type
	originalErr := ErrResourceNotFound
	wrappedErr := NewError(http.StatusInternalServerError, "WRAPPER", "wrapped", originalErr)

	var domainErr *Error
	if !errors.As(wrappedErr, &domainErr) {
		t.Error("errors.As should work with domain.Error")
	}

	if domainErr.Code != "WRAPPER" {
		t.Errorf("errors.As returned wrong error, got code %s, want WRAPPER", domainErr.Code)
	}
}

func TestError_WithErrors_Is(t *testing.T) {
	// Test that errors.Is works with wrapped errors
	underlyingErr := errors.New("underlying")
	domainErr := NewError(http.StatusBadRequest, "TEST", "test", underlyingErr)

	if !errors.Is(domainErr, underlyingErr) {
		t.Error("errors.Is should find underlying error")
	}
}

func TestErrNotImplemented(t *testing.T) {
	if ErrNotImplemented == nil {
		t.Error("ErrNotImplemented should not be nil")
	}

	if ErrNotImplemented.Error() != "not implemented" {
		t.Errorf("ErrNotImplemented.Error() = %s, want 'not implemented'", ErrNotImplemented.Error())
	}
}
