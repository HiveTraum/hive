package controllers

import (
	"auth/config"
	"auth/enums"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

func CreateSession(
	store infrastructure.StoreInterface,
	loginController infrastructure.LoginControllerInterface,
	ctx context.Context,
	body inout.CreateSessionRequestV1) (int, *models.Session) {

	status, userID := loginController.Login(ctx, body)
	if status != enums.Ok {
		return status, nil
	}
	if userID == uuid.Nil {
		return enums.CredentialsNotProvided, nil
	}

	if store.GetUser(ctx, userID) == nil {
		return enums.UserNotFound, nil
	}

	tokens := body.GetTokens()

	if tokens != nil {
		session := store.GetSession(ctx, body.Fingerprint, tokens.RefreshToken, userID)
		if session == nil {
			return enums.SessionNotFound, nil
		}
	}

	secret := store.GetActualSecret(ctx)
	status, session := store.CreateSession(ctx, body.Fingerprint, userID, secret.Id, body.UserAgent)
	if status != enums.Ok {
		return status, nil
	}
	userView := store.GetUserView(ctx, userID)
	if userView == nil {
		return enums.UserNotFound, nil
	}
	env := config.GetEnvironment()
	session.AccessToken = loginController.EncodeAccessToken(ctx, userID, userView.Roles, secret, time.Now().Add(time.Minute*time.Duration(env.AccessTokenLifetime)))

	return status, session
}
