package domain

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeTeamExists        ErrorCode = "TEAM_EXISTS"
	ErrorCodePullRequestExists ErrorCode = "PR_EXISTS"
	ErrorCodePullRequestMerged ErrorCode = "PR_MERGED"
	ErrorCodeReviewerMissing   ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate       ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
)

type ErrorMessage string

const (
	ErrorMessagePullRequestExists            ErrorMessage = "PR id already exists"
	ErrorMessageResourceNotFound             ErrorMessage = "resource not found"
	ErrorMessageTeamExists                   ErrorMessage = "team_name already exists"
	ErrorMessageReassignOnMerged             ErrorMessage = "cannot reassign on merged PR"
	ErrorMessageReviewerIsNotAssigned        ErrorMessage = "reviewer is not assigned to this PR"
	ErrorMessageNoActiveReplacementCandidate ErrorMessage = "no active replacement candidate in team"
)

type Error struct {
	Status  int
	Code    ErrorCode
	Message string
	Err     error
}

// ... i did not want to hardcode this. but your specification forces me to do so ...
var (
	ErrorPullRequestExists      = NewError(http.StatusConflict, ErrorCodePullRequestExists, string(ErrorMessagePullRequestExists), nil)
	ErrorResourceNotFound       = NewError(http.StatusNotFound, ErrorCodeNotFound, string(ErrorMessageResourceNotFound), nil)
	ErrorTeamExists             = NewError(http.StatusBadRequest, ErrorCodeTeamExists, string(ErrorMessageTeamExists), nil)
	ErrorInternal               = NewError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil)
	ErrorCantReassignOnMergedPr = NewError(http.StatusConflict, ErrorCodePullRequestMerged, string(ErrorMessageReassignOnMerged), nil)
	ErrorReviewerIsNotAssigned  = NewError(http.StatusConflict, ErrorCodeReviewerMissing, string(ErrorMessageReviewerIsNotAssigned), nil)
	ErrorNoCandidate            = NewError(http.StatusConflict, ErrorCodeNoCandidate, string(ErrorMessageNoActiveReplacementCandidate), nil)
)

func IsDomainError(err error) bool {
	var e *Error
	ok := errors.As(err, &e)
	return ok
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("%s", e.Code)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewError(status int, code ErrorCode, message string, err error) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var ErrNotImplemented = errors.New("not implemented")
