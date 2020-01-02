package controllers

import (
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/models"
	"context"
	"github.com/badoux/checkmail"
	"github.com/getsentry/sentry-go"
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

func getRandomStringCode() string {
	return functools.GetRandomString(6)
}

func checkEmailConfirmationCode(store infrastructure.StoreInterface, email string, code string) int {
	cachedCode := store.GetEmailConfirmationCode(email)
	if cachedCode == "" {
		return enums.EmailNotFound
	} else if cachedCode != code {
		return enums.IncorrectEmailCode
	}

	return enums.Ok
}

func validateEmail(store infrastructure.StoreInterface, email string, code string) (int, string) {
	email = getEmail(email)

	if email == "" {
		return enums.IncorrectEmail, ""
	}

	status := checkEmailConfirmationCode(store, email, code)
	return status, email
}

func CreateEmail(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, email string, code string, userId models.UserID) (int, *models.Email) {

	status, email := validateEmail(store, email, code)

	if status != enums.Ok {
		return status, nil
	}

	_, oldEmail := store.GetEmail(ctx, email)

	identifiers := []models.UserID{userId}

	if oldEmail != nil {
		identifiers = append(identifiers, oldEmail.UserId)
	}

	status, phoneObject := store.CreateEmail(ctx, userId, email)
	esb.OnEmailChanged(identifiers)
	return status, phoneObject
}

func CreateEmailConfirmation(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, email string) (int, *models.EmailConfirmation) {

	email = getEmail(email)

	if email == "" {
		return enums.IncorrectEmail, nil
	}

	code := getRandomStringCode()
	emailConfirmation := store.CreateEmailConfirmationCode(email, code, time.Minute*15)
	esb.OnEmailCodeConfirmationCreated(email, code)
	return enums.Ok, emailConfirmation
}
