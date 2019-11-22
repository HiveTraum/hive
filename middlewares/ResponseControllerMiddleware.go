package middlewares

import (
	"auth/functools"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"net/http"
)

type ResponseControllerHandler func(*functools.Request) (int, proto.Message)

func ResponseControllerMiddleware(next ResponseControllerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		request := &functools.Request{
			Request: r,
		}

		status, data := next(request)

		var err error

		var bytes []byte

		if data != nil {
			if request.IsProto() {
				bytes, err = proto.Marshal(data)
			} else {
				bytes, err = json.Marshal(data)
			}

			if err != nil {
				sentry.CaptureException(err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(status)
				_, _ = w.Write(bytes)
			}
		} else {
			w.WriteHeader(status)
		}
	}
}
