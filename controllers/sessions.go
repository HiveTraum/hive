package controllers

import (
	"auth/enums"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (controller *Controller) CreateSession(ctx context.Context, userID uuid.UUID, fingerprint, userAgent string) *models.Session {
	secret := controller.GetActualSecret(ctx)
	session := controller.store.CreateSession(ctx, userID, secret.Id, fingerprint, userAgent)
	user := controller.GetUserView(ctx, userID)
	session.AccessToken = controller.accessTokenEncoder(ctx, user.Id, user.Roles, secret, session.Expires)
	return session
}

func (controller *Controller) UpdateSession(ctx context.Context, id uuid.UUID, fingerprint, userAgent string) (int, *models.Session) {
	oldSession := controller.store.DeleteSession(ctx, id)
	if oldSession == nil ||
		oldSession.Fingerprint != fingerprint ||
		oldSession.Expires <= time.Now().Unix() {
		return enums.SessionNotFound, nil
	}

	return enums.Ok, controller.CreateSession(ctx, oldSession.UserID, fingerprint, userAgent)
}
