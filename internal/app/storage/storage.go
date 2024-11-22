package storage

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrCreateLink = errors.New("failed to create link")
)
