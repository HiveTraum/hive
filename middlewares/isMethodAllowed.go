package middlewares

import (
	"auth/functools"
	"net/http"
)

func IsMethodAllowedMiddleware(methods []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !functools.Contains(r.Method, methods) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			next(w, r)
		}
	}
}
