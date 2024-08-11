package app

import (
	"net/http"
)

func methodHandler(c Controller) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.Get(w, r)
		case http.MethodPost:
			c.Post(w, r)
		default:
			pageNotFound(w, r)
		}
	})
}
