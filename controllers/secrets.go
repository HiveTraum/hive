package controllers

import (
	"auth/infrastructure"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
)

func GetSecret(store infrastructure.StoreInterface, ctx context.Context, id uuid.UUID) *models.Secret {
	secret := store.GetSecret(ctx, id)
	return secret
}
