package middlewares

import (
	"hive/enums"
	"hive/functools"
	"hive/repositories"
	"net/http"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := string(repositories.GetContentTypeHeader(r))

		if !functools.Contains(contentType, []string{string(enums.JSONContentType), string(enums.BinaryContentType)}) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
		} else {
			w.Header().Set("content-type", contentType)
			next.ServeHTTP(w, r)
		}
	})
}
