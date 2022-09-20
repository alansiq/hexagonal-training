package rest

import (
	"net/http"
)

//go:generate go run github.com/golang/mock/mockgen -destination mock_dao/doer.go . HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
