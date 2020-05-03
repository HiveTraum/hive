package repositories

import (
	"auth/config"
	"auth/enums"
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

func GetAuthorizationCookie(r *http.Request, environment *config.Environment) string {
	cookie, err := r.Cookie(environment.RefreshTokenCookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}
