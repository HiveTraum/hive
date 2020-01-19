package middlewares

import (
	"auth/functools"
	"net/http"
)

func ContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := functools.Request{
			Request: r,
		}

		if !request.IsContentTypeAllowed(nil) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Set("content-type", request.GetContentType())

		next(w, r)
	}
}
