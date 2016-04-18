package errs

import "github.com/ansel1/merry"

var (
	NotFound      = merry.New("the specified resource was not found or insufficient permissions")
	InvalidFilter = merry.New("the given query filter is invalid")
	SeedsNotSync  = merry.New("local and distant seeds does not match")
)
