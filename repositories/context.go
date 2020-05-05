package repositories

import (
	"hive/models"
	"context"
)

const (
	AuthenticatedUser string = "authenticatedUser"
)

func GetUserFromContext(ctx context.Context) models.IAuthenticationBackendUser {
	user, _ := ctx.Value(AuthenticatedUser).(models.IAuthenticationBackendUser)
	return user
}

func SetUserToContext(ctx context.Context, user models.IAuthenticationBackendUser) context.Context {
	return context.WithValue(ctx, AuthenticatedUser, user)
}
