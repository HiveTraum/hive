package api

import (
	"auth/functools"
	"github.com/getsentry/sentry-go"
	"net/http"
	"strconv"
)

func unhandledStatus(r *functools.Request, status int) int {

	request := sentry.Request{}
	request.FromHTTPRequest(r.Request)

	sentry.CaptureEvent(&sentry.Event{
		Level:   sentry.LevelError,
		Message: "Unhandled controller status",
		Tags:    map[string]string{"controller status": strconv.Itoa(status)},
		Request: request,
	})

	return http.StatusInternalServerError
}
