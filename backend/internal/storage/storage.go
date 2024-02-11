package storage

import "errors"

var (
	ErrExpressionNotFound   = errors.New("expression not found")
	ErrExpressionInProgress = errors.New("expression in progress")
)
