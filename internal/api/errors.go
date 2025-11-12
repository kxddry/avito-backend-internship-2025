package api

import (
	"github.com/kxddry/avito-backend-internship-2025/internal/api/generated"
	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

type errorPayload struct {
	status int
	body   generated.ErrorResponse
}

func makeErrorPayload(err *domain.Error) errorPayload {
	message := err.Message
	if message == "" {
		message = string(err.Code)
	}

	payload := generated.ErrorResponse{
		Error: struct {
			Code    generated.ErrorResponseErrorCode `json:"code"`
			Message string                           `json:"message"`
		}{
			Code:    generated.ErrorResponseErrorCode(err.Code),
			Message: message,
		},
	}

	return errorPayload{
		status: err.Status,
		body:   payload,
	}
}
