package middlewares

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func IsLocalRequestMiddleware(localAddress string) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.RemoteAddr, localAddress) {
				writer.WriteHeader(http.StatusNotFound)
			} else {
				handler.ServeHTTP(writer, request)
			}
		})
	}
}
