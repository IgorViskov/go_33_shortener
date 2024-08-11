package app

import (
	"net/http"
)

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
