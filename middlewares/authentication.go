package middlewares

import (
	"auth/auth"
	"auth/enums"
	"auth/repositories"
	"net/http"
)

func AuthenticationMiddleware(authenticationController auth.IAuthenticationController) func(http.Handler, bool) http.Handler {
	return func(next http.Handler, required bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			status, user := authenticationController.Login(ctx, r)
			if status == enums.Ok && user != nil || !required {
				ctx = repositories.SetUserToContext(ctx, user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
}
