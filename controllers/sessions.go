package controllers

import (
	"auth/config"
	"auth/enums"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (controller *Controller) CreateSession(ctx context.Context, userID uuid.UUID, userAgent string, fingerprint string) (int, *models.Session) {
	secret := controller.GetActualSecret(ctx)
	status, session := controller.store.CreateSession(ctx, fingerprint, userID, secret.Id, userAgent)
	if status != enums.Ok {
		return status, nil
	}
	env := config.GetEnvironment()
	expires := time.Now().Add(time.Minute * time.Duration(env.AccessTokenLifetime))
	session.Expires = expires.Unix()
	return status, session
}
