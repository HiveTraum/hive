package middlewares

import (
	"auth/auth"
	"auth/enums"
	"auth/repositories"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthenticationMiddleware(authenticationController auth.IAuthenticationController) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			status, user := authenticationController.Login(ctx, r)
			if status == enums.Ok || user == nil {
				ctx = repositories.SetUserToContext(ctx, user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
}
