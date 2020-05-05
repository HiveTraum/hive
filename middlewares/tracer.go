package middlewares

import (
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func TracerMiddleware(tracer opentracing.Tracer) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return gorilla.Middleware(tracer, handler)
	}
}
