package models

import "errors"

var (
	ErrFailedToParsePage  = errors.New("page number is invalid")
	ErrFailedToParseLimit = errors.New("limit number is invalid")
	ErrFailedToParseID    = errors.New("task id is invalid")
	ErrNotFound           = errors.New("nothing was found")
)
