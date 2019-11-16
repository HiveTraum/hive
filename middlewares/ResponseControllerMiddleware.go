package middlewares

import (
	"auth/functools"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"net/http"
	"reflect"
)

type ResponseControllerHandler func(*functools.Request) (int, proto.Message)

func ResponseControllerMiddleware(next ResponseControllerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := &functools.Request{
			Request: r,
		}

		status, data := next(request)

		var err error

		var bytes []byte

		if !reflect.ValueOf(data).IsNil() {
			if request.IsProto() {
				bytes, err = proto.Marshal(data)
			} else {
				bytes, err = json.Marshal(data)
			}
		}

		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(status)
			_, _ = w.Write(bytes)
		}
	}
}
