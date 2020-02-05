package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
	"github.com/badoux/checkmail"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func getEmail(email string) string {
	err := checkmail.ValidateFormat(email)

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return email
}

func checkEmailConfirmationCode(ctx context.Context, store infrastructure.StoreInterface, email string, code string) int {
	cachedCode := store.GetEmailConfirmationCode(ctx, email)
	if cachedCode == "" {
		return enums.EmailNotFound
	} else if cachedCode != code {
		return enums.IncorrectEmailCode
	}

	return enums.Ok
}

func validateEmail(ctx context.Context, store infrastructure.StoreInterface, email string, code string) (int, string) {
	email = getEmail(email)

	if email == "" {
		return enums.IncorrectEmail, ""
	}

	status := checkEmailConfirmationCode(ctx, store, email, code)
	return status, email
}

func CreateEmail(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, email string, code string, userId uuid.UUID) (int, *models.Email) {

	status, email := validateEmail(ctx, store, email, code)

	if status != enums.Ok {
		return status, nil
	}

	_, oldEmail := store.GetEmail(ctx, email)

	identifiers := []uuid.UUID{userId}

	if oldEmail != nil {
		identifiers = append(identifiers, oldEmail.UserId)
	}

	status, phoneObject := store.CreateEmail(ctx, userId, email)
	esb.OnEmailChanged(identifiers)
	return status, phoneObject
}

func CreateEmailConfirmation(ctx context.Context, store infrastructure.StoreInterface, esb infrastructure.ESBInterface, email string) (int, *models.EmailConfirmation) {

	email = getEmail(email)

	if email == "" {
		return enums.IncorrectEmail, nil
	}

	code := store.GetRandomCodeForEmailConfirmation()
	emailConfirmation := store.CreateEmailConfirmationCode(ctx, email, code, time.Minute*15)
	esb.OnEmailCodeConfirmationCreated(email, code)
	return enums.Ok, emailConfirmation
}
