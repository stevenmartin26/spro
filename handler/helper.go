package handler

import (
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
)

func constructErrorResponse(err *common.CustomError) (int, generated.ErrorResponse) {
	response := generated.ErrorResponse{
		Message: err.Message,
	}

	if len(err.Details) > 0 {
		var errDetails []struct {
			Error string
		}

		for _, detail := range err.Details {
			errDetails = append(errDetails, struct{ Error string }{detail})
		}

		response.Details = (*[]struct {
			Error string "json:\"error\""
		})(&errDetails)
	}

	var statusCode int
	switch err.ErrType {
	case common.ErrInvalidInput:
		statusCode = http.StatusBadRequest
	case common.ErrUnauthorized:
		statusCode = http.StatusForbidden
	case common.ErrEntityAlreadyExists:
		statusCode = http.StatusConflict
	case common.ErrTooManyAttempts:
		statusCode = http.StatusTooManyRequests
	default:
		statusCode = http.StatusInternalServerError
	}

	return statusCode, response
}

func extractAccessToken(authParam string) (string, *generated.ErrorResponse) {
	if authParam == "" || !strings.HasPrefix(authParam, "Bearer ") {
		response := generated.ErrorResponse{
			Message: "invalid access token",
		}
		return "", &response
	}

	return strings.TrimPrefix(authParam, "Bearer "), nil
}
