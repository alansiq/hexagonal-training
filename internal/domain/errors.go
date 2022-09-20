package domain

import "errors"

var (
	ErrApiException = errors.New("api error")
	ErrNotFound     = errors.New("not found")
)
