package errs

import "github.com/solher/snakepit"

var (
	APIInternal     = snakepit.APIInternal
	APIUnauthorized = snakepit.APIError{
		Description: "Authorization Required.",
		ErrorCode:   "AUTHORIZATION_REQUIRED",
	}
	APIForbidden = snakepit.APIError{
		Description: "The specified resource was not found or you don't have sufficient permissions.",
		ErrorCode:   "FORBIDDEN",
	}
	APIFilterDecoding = snakepit.APIError{
		Description: "Could not decode the given filter.",
		ErrorCode:   "FILTER_DECODING_ERROR",
	}
	APIBodyDecoding = snakepit.APIError{
		Description: "Could not decode the JSON request.",
		ErrorCode:   "BODY_DECODING_ERROR",
	}
	APIValidation = snakepit.APIError{
		Description: "The model validation failed.",
		ErrorCode:   "VALIDATION_ERROR",
	}
	APIInvalidFilter = snakepit.APIError{
		Description: "The given filter is invalid.",
		ErrorCode:   "INVALID_FILTER",
	}
)
