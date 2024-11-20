package utils

import (
	"errors"
	"jwt-auth-app/types"
)

var (
	ErrUserExists         = errors.New("USER_EXISTS")
	ErrUserNotFound       = errors.New("USER_NOT_FOUND")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")
	ErrInternalServer     = errors.New("INTERNAL_SERVER_ERROR")
	ErrUnauthorized       = errors.New("UNAUTHORIZED")
	ErrMissingAuthHeader  = errors.New("MISSING_AUTH_HEADER")
	ErrInvalidAuthHeader  = errors.New("INVALID_AUTH_HEADER")
	ErrInvalidToken       = errors.New("INVALID_TOKEN")
)

func GetErrorResponse(err error) (int, types.ErrorResponse) {
	switch err {
	case ErrUserExists:
		return 409, types.ErrorResponse{
			Code:    "USER_EXISTS",
			Message: "User with this email already exists",
		}
	case ErrUserNotFound:
		return 404, types.ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		}
	case ErrInvalidCredentials:
		return 401, types.ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		}
	case ErrUnauthorized:
		return 401, types.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Unauthorized",
		}
	case ErrMissingAuthHeader:
		return 401, types.ErrorResponse{
			Code:    "MISSING_AUTH_HEADER",
			Message: "Authorization header is missing",
		}
	case ErrInvalidAuthHeader:
		return 401, types.ErrorResponse{
			Code:    "INVALID_AUTH_HEADER",
			Message: "Invalid authorization header",
		}
	case ErrInvalidToken:
		return 401, types.ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: "Invalid token",
		}
	case ErrInternalServer:
		return 500, types.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		}
	default:
		return 500, types.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		}
	}
}
