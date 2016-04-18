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

// An internal error occured. Please retry later.
// swagger:response InternalResponse
type internalResponse struct {
	// in: body
	Body snakepit.APIError
}

// The specified resource was not found or you do not have sufficient permissions.
// swagger:response UnauthorizedResponse
type unauthorizedResponse struct {
	// in: body
	Body snakepit.APIError
}

// Could not decode the given filter.
// swagger:response FilterDecodingResponse
type filterDecodingResponse struct {
	// in: body
	Body snakepit.APIError
}

// Could not decode the JSON request.
// swagger:response BodyDecodingResponse
type bodyDecodingResponse struct {
	// in: body
	Body snakepit.APIError
}

// The model validation failed.
// swagger:response ValidationResponse
type validationResponse struct {
	// in: body
	Body snakepit.APIError
}

// The given filter is invalid.
// swagger:response InvalidFilterResponse
type invalidFilterResponse struct {
	// in: body
	Body snakepit.APIError
}
