package app

import (
	"net/http"
)

type Controller interface {
	Get(w http.ResponseWriter, req *http.Request)
	Post(w http.ResponseWriter, req *http.Request)
	GetPath() string
}
