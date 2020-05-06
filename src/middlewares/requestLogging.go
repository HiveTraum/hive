package middlewares

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"hive/config"
	"net/http"
	"os"
	"time"
)

func RequestLoggingMiddleware(environment *config.Environment) mux.MiddlewareFunc {

	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Str("service", environment.ServiceName).
		Str("instance", environment.Instance).
		Logger()

	chain := alice.New()
	chain = chain.Append(hlog.NewHandler(logger))
	chain = chain.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Send()
	}))
	chain = chain.Append(hlog.RemoteAddrHandler("ip"))
	chain = chain.Append(hlog.UserAgentHandler("user_agent"))
	chain = chain.Append(hlog.RefererHandler("referer"))
	chain = chain.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	return chain.Then
}
