package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
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
	status, session := store.CreateSession(ctx, body.Fingerprint, user.Id, secret.Id)
	userView := store.GetUserView(ctx, user.Id)
	token := loginController.EncodeAccessToken(ctx, user.Id, userView.Roles, secret.Value)

	return status, session, token
}
