package errs

import "github.com/solher/snakepit"

type (
	ErrDatabase      struct{ snakepit.InternalError }
	ErrNotFound      struct{ snakepit.InternalError }
	ErrInvalidFilter struct{ snakepit.InternalError }
	ErrSeedsNotSync  struct{ snakepit.InternalError }
	ErrService       struct{ snakepit.InternalError }
)

var (
	Database = ErrDatabase{
		snakepit.NewInternalError("undefined database error"),
	}
	NotFound = ErrNotFound{
		snakepit.NewInternalError("the specified resource was not found or insufficient permissions"),
	}
	InvalidFilter = ErrInvalidFilter{
		snakepit.NewInternalError("the given query filter is invalid"),
	}
	SeedsNotSync = ErrSeedsNotSync{
		snakepit.NewInternalError("local and distant seeds does not match"),
	}
	ServiceError = ErrService{
		snakepit.NewInternalError("an internal service failed to respond"),
	}
)
