package controllers

import (
	"auth/config"
	"auth/enums"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
	"time"
)

func CreateSession(
	store infrastructure.StoreInterface,
	loginController infrastructure.LoginControllerInterface,
	ctx context.Context,
	body inout.CreateSessionRequestV1) (int, *models.Session, string) {

	status, user := loginController.Login(ctx, body)
	if status != enums.Ok {
		return status, nil, ""
	}

	secret := store.GetActualSecret(ctx)
	status, session := store.CreateSession(ctx, body.Fingerprint, user.Id, secret.Id, body.UserAgent)
	userView := store.GetUserView(ctx, user.Id)
	env := config.GetEnvironment()
	token := loginController.EncodeAccessToken(ctx, user.Id, userView.Roles, secret.Value, time.Now().Add(time.Minute*time.Duration(env.AccessTokenLifetime)))

	return status, session, token
}
