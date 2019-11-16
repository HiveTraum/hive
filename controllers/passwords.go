package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSaltValue(value string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.MinCost)
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return string(hash)
}

func CreatePassword(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, userId int64, value string) (int, *models.Password) {
	value = hashAndSaltValue(value)
	if value == "" {
		return enums.IncorrectPassword, nil
	}

	status, password := store.CreatePassword(ctx, userId, value)
	if status == enums.Ok {
		esb.OnPasswordChanged(userId)
	}

	return status, password
}
