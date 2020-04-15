package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"context"
	"strings"
)

type LoginController struct {
	Backends              map[string]infrastructure.AuthenticationBackend
	RequestContextUserKey string
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

func (controller *LoginController) Login(ctx context.Context, authorizationHeader string) (int, infrastructure.AuthenticationBackendUser, context.Context) {

	if user, ok := ctx.Value(controller.RequestContextUserKey).(infrastructure.AuthenticationBackendUser); ok {
		return enums.Ok, user, ctx
	}

	tokenType, token := controller.GetToken(authorizationHeader)

	if tokenType == "" || token == "" {
		return enums.Ok, nil, ctx
	}

	backend := controller.Backends[tokenType]
	if backend == nil {
		return enums.BackendNotFound, nil, ctx
	}

	status, user := backend.GetUser(ctx, token)
	if status == enums.Ok || user != nil {
		ctx = context.WithValue(ctx, controller.RequestContextUserKey, user)
	}

	return status, user, ctx
}
