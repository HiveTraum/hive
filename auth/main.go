package auth

import (
	"auth/auth/backends"
	"auth/config"
	"auth/enums"
	"auth/models"
	"auth/repositories"
	"context"
	"net/http"
	"strings"
)

type IAuthenticationController interface {
	Login(ctx context.Context, r *http.Request) (int, models.IAuthenticationBackendUser)
}

type AuthenticationController struct {
	backends    map[string]backends.IAuthenticationBackend
	environment *config.Environment
}

func InitAuthController(backends map[string]backends.IAuthenticationBackend, environment *config.Environment) *AuthenticationController {
	return &AuthenticationController{backends: backends, environment: environment}
}

func (controller *AuthenticationController) GetToken(authorizationHeader string) (string, string) {

	if authorizationHeader == "" {
		return "", ""
	}

	parts := strings.Split(authorizationHeader, " ")

	if len(parts) < 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

func (controller *AuthenticationController) Login(ctx context.Context, r *http.Request) (int, models.IAuthenticationBackendUser) {

	authorizationHeader := repositories.GetAuthorizationHeader(r)
	tokenType, token := controller.GetToken(authorizationHeader)
	if tokenType == "" || token == "" {
		return enums.Ok, nil
	}

	backend := controller.backends[tokenType]
	if backend == nil {
		return enums.BackendNotFound, nil
	}

	refreshToken := repositories.GetAuthorizationCookie(r, controller.environment)

	status, user := backend.GetUser(ctx, token, refreshToken)
	return status, user
}
