package errs

import "github.com/palantir/stacktrace"

const (
	EcodeDatabase = stacktrace.ErrorCode(iota)
	EcodeNotFound
	EcodeInvalidFilter
	EcodeSeedsNotSync
	EcodeServiceError
)

var (
	ErrDatabase = stacktrace.NewMessageWithCode(
		EcodeDatabase,
		"undefined database error",
	)
	ErrNotFound = stacktrace.NewMessageWithCode(
		EcodeNotFound,
		"the specified resource was not found or you do not have sufficient permissions",
	)
	ErrInvalidFilter = stacktrace.NewMessageWithCode(
		EcodeInvalidFilter,
		"the given query filter is invalid",
	)
	ErrSeedsNotSync = stacktrace.NewMessageWithCode(
		EcodeSeedsNotSync,
		"local and distant seeds does not match",
	)
	ErrServiceError = stacktrace.NewMessageWithCode(
		EcodeServiceError,
		"an internal service failed to respond",
	)
)
