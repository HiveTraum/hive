package repositories

import (
	"hive/config"
	"hive/enums"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

const (
	ContentType   = "content-type"
	Authorization = "authorization"
)

func GetContentTypeHeader(r *http.Request) enums.ContentType {
	return enums.ContentType(r.Header.Get(ContentType))
}

func GetAuthorizationHeader(r *http.Request) string {
	return r.Header.Get(Authorization)
}

func GetRefreshTokenCookie(r *http.Request, environment *config.Environment) *uuid.UUID {
	cookie, err := r.Cookie(environment.RefreshTokenCookieName)
	if err != nil {
		return nil
	}
	value, err := uuid.FromString(cookie.Value)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	return &value
}
