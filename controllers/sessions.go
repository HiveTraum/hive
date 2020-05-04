package controllers

import (
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
	expires := time.Now().Add(time.Minute * time.Duration(controller.environment.AccessTokenLifetime))
	session.Expires = expires.Unix()
	return status, session
}
