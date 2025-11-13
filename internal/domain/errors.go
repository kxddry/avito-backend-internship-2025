package domain

import (
	"errors"
	"net/http"
)

// ErrorCode is a code for the error.
type ErrorCode string

// Error codes.
const (
	ErrorCodeTeamExists        ErrorCode = "TEAM_EXISTS"
	ErrorCodePullRequestExists ErrorCode = "PR_EXISTS"
	ErrorCodePullRequestMerged ErrorCode = "PR_MERGED"
	ErrorCodeReviewerMissing   ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate       ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
)

// ErrorMessage is a message for the error.
type ErrorMessage string

// Error messages.
const (
	ErrorMessagePullRequestExists            ErrorMessage = "PR id already exists"
	ErrorMessageResourceNotFound             ErrorMessage = "resource not found"
	ErrorMessageTeamExists                   ErrorMessage = "team_name already exists"
	ErrorMessageReassignOnMerged             ErrorMessage = "cannot reassign on merged PR"
	ErrorMessageReviewerIsNotAssigned        ErrorMessage = "reviewer is not assigned to this PR"
	ErrorMessageNoActiveReplacementCandidate ErrorMessage = "no active replacement candidate in team"
)

// Error is a domain error.
type Error struct {
	Status  int
	Code    ErrorCode
	Message string
	Err     error
}

// I did not want to hardcode this, but your specification forces me to do so.
var (
	ErrPRExists         = NewError(http.StatusConflict, ErrorCodePullRequestExists, string(ErrorMessagePullRequestExists), nil)
	ErrResourceNotFound = NewError(http.StatusNotFound, ErrorCodeNotFound, string(ErrorMessageResourceNotFound), nil)
	ErrTeamExists       = NewError(http.StatusBadRequest, ErrorCodeTeamExists, string(ErrorMessageTeamExists), nil)
	ErrInternal         = NewError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil)
	ErrReassignOnMerged = NewError(http.StatusConflict, ErrorCodePullRequestMerged, string(ErrorMessageReassignOnMerged), nil)
	ErrReviewerMissing  = NewError(http.StatusConflict, ErrorCodeReviewerMissing, string(ErrorMessageReviewerIsNotAssigned), nil)
	ErrNoCandidate      = NewError(http.StatusConflict, ErrorCodeNoCandidate, string(ErrorMessageNoActiveReplacementCandidate), nil)
)

// IsDomainError checks if the error is a domain error.
func IsDomainError(err error) bool {
	var e *Error
	ok := errors.As(err, &e)
	return ok
}

// Error returns the error message.
func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new domain error.
func NewError(status int, code ErrorCode, message string, err error) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ErrNotImplemented is a not implemented error.
var ErrNotImplemented = errors.New("not implemented")
