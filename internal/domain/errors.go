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
	ErrorMessagePullRequestExists ErrorMessage = "PR id already exists"
	ErrorMessageResourceNotFound  ErrorMessage = "resource not found"
	ErrorMessageTeamExists        ErrorMessage = "team_name already exists"
)

type Error struct {
	Status  int
	Code    ErrorCode
	Message string
	Err     error
}

// ... i did not want to hardcode this. but your specification forces me to do so ...
var (
	ErrorPullRequestExists = NewError(http.StatusConflict, ErrorCodePullRequestExists, string(ErrorMessagePullRequestExists), nil)
	ErrorResourceNotFound  = NewError(http.StatusNotFound, ErrorCodeNotFound, string(ErrorMessageResourceNotFound), nil)
	ErrorTeamExists        = NewError(http.StatusBadRequest, ErrorCodeTeamExists, string(ErrorMessageTeamExists), nil)
)

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

func NewTeamExistsError(message string, err error) *Error {
	return NewError(http.StatusBadRequest, ErrorCodeTeamExists, message, err)
}

func NewPullRequestExistsError(message string, err error) *Error {
	return NewError(http.StatusConflict, ErrorCodePullRequestExists, message, err)
}

func NewPullRequestMergedError(message string, err error) *Error {
	return NewError(http.StatusConflict, ErrorCodePullRequestMerged, message, err)
}

func NewReviewerMissingError(message string, err error) *Error {
	return NewError(http.StatusConflict, ErrorCodeReviewerMissing, message, err)
}

func NewNoCandidateError(message string, err error) *Error {
	return NewError(http.StatusConflict, ErrorCodeNoCandidate, message, err)
}

func NewNotFoundError(message string, err error) *Error {
	return NewError(http.StatusNotFound, ErrorCodeNotFound, message, err)
}

func NewUnauthorizedError(message string, err error) *Error {
	return NewError(http.StatusUnauthorized, ErrorCodeNotFound, message, err)
}

var ErrNotImplemented = errors.New("not implemented")
