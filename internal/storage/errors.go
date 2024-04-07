package storage

import "errors"

var (
	ErrNotFound      = errors.New("no such item")
	ErrInsertInCache = errors.New("cant insert in cache")
)
