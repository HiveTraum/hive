package controllers

import (
	"hive/enums"
	"hive/models"
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

func (controller *Controller) checkEmailConfirmationCode(ctx context.Context, email string, code string) int {
	cachedCode := controller.store.GetEmailConfirmationCode(ctx, email)
	if code == "" {
		return enums.EmailConfirmationCodeRequired
	} else if cachedCode == "" {
		return enums.EmailConfirmationCodeNotFound
	} else if cachedCode != code {
		return enums.IncorrectEmailCode
	}

	return enums.Ok
}

func (controller *Controller) validateEmail(ctx context.Context, email string, code string) (int, string) {
	email = getEmail(email)

	if email == "" {
		return enums.IncorrectEmail, ""
	}

	status := controller.checkEmailConfirmationCode(ctx, email, code)
	return status, email
}

func (controller *Controller) CreateEmail(ctx context.Context, email string, code string, userId uuid.UUID) (int, *models.Email) {

	status, email := controller.validateEmail(ctx, email, code)

	if status != enums.Ok {
		return status, nil
	}

	_, oldEmail := controller.store.GetEmail(ctx, email)

	identifiers := []uuid.UUID{userId}

	if oldEmail != nil {
		identifiers = append(identifiers, oldEmail.UserId)
	}

	status, phoneObject := controller.store.CreateEmail(ctx, userId, email)
	controller.OnEmailChanged(identifiers)
	return status, phoneObject
}

func (controller *Controller) CreateEmailConfirmation(ctx context.Context, email string) (int, *models.EmailConfirmation) {

	email = getEmail(email)

	if email == "" {
		return enums.IncorrectEmail, nil
	}

	code := controller.store.GetRandomCodeForEmailConfirmation()
	emailConfirmation := controller.store.CreateEmailConfirmationCode(ctx, email, code, time.Minute*15)
	controller.OnEmailCodeConfirmationCreated(email, code)
	return enums.Ok, emailConfirmation
}
