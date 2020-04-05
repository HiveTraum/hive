package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"context"
	"strings"
)

type LoginController struct {
	Backends map[string]infrastructure.AuthenticationBackend
}

func (controller *LoginController) GetToken(authorizationHeader string) (string, string) {

	if authorizationHeader == "" {
		return "", ""
	}

	parts := strings.Split(authorizationHeader, " ")

	if len(parts) < 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

// Login

func (controller *LoginController) Login(ctx context.Context, authorizationHeader string) (int, infrastructure.AuthenticationBackendUser) {

	tokenType, token := controller.GetToken(authorizationHeader)

	if tokenType == "" || token == "" {
		return enums.Ok, nil
	}

	backend := controller.Backends[tokenType]
	if backend == nil {
		return enums.BackendNotFound, nil
	}

	return backend.GetUser(ctx, token)
}
