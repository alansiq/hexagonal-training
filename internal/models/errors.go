package models

import "errors"

var (
	ErrApiException = errors.New("api error")
	ErrNotFound     = errors.New("not found")
)
