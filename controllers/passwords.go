package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
)

func CreatePassword(
	store infrastructure.StoreInterface,
	esb infrastructure.ESBInterface,
	passwordProcessor infrastructure.LoginControllerInterface,
	ctx context.Context,
	userId models.UserID,
	value string) (int, *models.Password) {

	value = passwordProcessor.EncodePassword(ctx, value)
	if value == "" {
		return enums.IncorrectPassword, nil
	}

	status, password := store.CreatePassword(ctx, userId, value)
	if status == enums.Ok {
		esb.OnPasswordChanged(userId)
	}

	return status, password
}
