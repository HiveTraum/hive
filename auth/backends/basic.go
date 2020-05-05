package backends

import (
	"hive/config"
	"hive/enums"
	"hive/functools"
	"hive/models"
	"hive/passwordProcessors"
	"hive/stores"
	"context"
	"encoding/base64"
	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type BasicAuthenticationBackendUser struct {
	IsAdmin bool
	Roles   []string
	UserID  uuid.UUID
}

func (user *BasicAuthenticationBackendUser) GetIsAdmin() bool {
	return user.IsAdmin
}

func (user *BasicAuthenticationBackendUser) GetRoles() []string {
	return user.Roles
}

func (user *BasicAuthenticationBackendUser) GetUserID() uuid.UUID {
	return user.UserID
}

type BasicAuthenticationBackend struct {
	store             stores.IStore
	passwordProcessor passwordProcessors.IPasswordProcessor
	environment       *config.Environment
}

func InitBasicAuthenticationBackend(store stores.IStore, passwordProcessor passwordProcessors.IPasswordProcessor, environment *config.Environment) *BasicAuthenticationBackend {
	return &BasicAuthenticationBackend{
		store:             store,
		passwordProcessor: passwordProcessor,
		environment:       environment,
	}
}

type BasicAuthenticationBackendWithMockedInternals struct {
	Backend           *BasicAuthenticationBackend
	Store             *stores.MockIStore
	PasswordProcessor *passwordProcessors.MockIPasswordProcessor
}

func InitBasicAuthenticationWithMockedInternals(ctrl *gomock.Controller) *BasicAuthenticationBackendWithMockedInternals {
	store := stores.NewMockIStore(ctrl)
	passwordProcessor := passwordProcessors.NewMockIPasswordProcessor(ctrl)
	return &BasicAuthenticationBackendWithMockedInternals{
		Backend:           InitBasicAuthenticationBackend(store, passwordProcessor, config.InitEnvironment()),
		Store:             store,
		PasswordProcessor: passwordProcessor,
	}
}

func (backend *BasicAuthenticationBackend) getUserFromEmailAndCode(ctx context.Context, emailValue string, emailCode string) (int, *uuid.UUID) {

	emailValue = functools.NormalizeEmail(emailValue)
	if emailValue == "" {
		return enums.IncorrectEmail, nil
	}

	code := backend.store.GetEmailConfirmationCode(ctx, emailValue)
	if code == "" {
		return enums.EmailConfirmationCodeNotFound, nil
	} else if code != emailCode {
		return enums.IncorrectEmailCode, nil
	}

	status, email := backend.store.GetEmail(ctx, emailValue)
	if status != enums.Ok {
		return status, nil
	}
	if email == nil {
		return enums.EmailNotFound, nil
	}

	return enums.Ok, &email.UserId
}

func (backend *BasicAuthenticationBackend) getUserFromEmailAndPassword(ctx context.Context, emailValue string, passwordValue string) (int, *uuid.UUID) {

	emailValue = functools.NormalizeEmail(emailValue)
	if emailValue == "" {
		return enums.IncorrectEmail, nil
	}

	status, email := backend.store.GetEmail(ctx, emailValue)
	if status != enums.Ok {
		return status, nil
	}
	if email == nil {
		return enums.EmailNotFound, nil
	}

	status, password := backend.store.GetLatestPassword(ctx, email.UserId)
	if status != enums.Ok {
		return status, nil
	}
	if password == nil {
		return enums.PasswordNotFound, nil
	}

	passwordVerified := backend.passwordProcessor.VerifyPassword(ctx, passwordValue, password.Value)
	if !passwordVerified {
		return enums.IncorrectPassword, nil
	}

	return enums.Ok, &password.UserId
}

func (backend *BasicAuthenticationBackend) getUserFromPhoneAndPassword(ctx context.Context, phoneValue string, passwordValue string) (int, *uuid.UUID) {
	phoneValue = functools.NormalizePhone(phoneValue)
	if phoneValue == "" {
		return enums.IncorrectPhone, nil
	}

	status, phone := backend.store.GetPhone(ctx, phoneValue)
	if status != enums.Ok {
		return status, nil
	}
	if phone == nil {
		return enums.PhoneNotFound, nil
	}

	status, password := backend.store.GetLatestPassword(ctx, phone.UserId)
	if status != enums.Ok {
		return status, nil
	}
	if password == nil {
		return enums.PasswordNotFound, nil
	}

	passwordVerified := backend.passwordProcessor.VerifyPassword(ctx, passwordValue, password.Value)
	if !passwordVerified {
		return enums.IncorrectPassword, nil
	}

	return enums.Ok, &password.UserId
}

func (backend *BasicAuthenticationBackend) getUserFromPhoneAndCode(ctx context.Context, phoneValue string, phoneCode string) (int, *uuid.UUID) {
	phoneValue = functools.NormalizePhone(phoneValue)
	if phoneValue == "" {
		return enums.IncorrectPhone, nil
	}

	code := backend.store.GetPhoneConfirmationCode(ctx, phoneValue)
	if code == "" {
		return enums.PhoneConfirmationCodeNotFound, nil
	} else if code != phoneCode {
		return enums.IncorrectPhoneCode, nil
	}

	status, phone := backend.store.GetPhone(ctx, phoneValue)
	if status != enums.Ok {
		return status, nil
	}
	if phone == nil {
		return enums.PhoneNotFound, nil
	}

	return enums.Ok, &phone.UserId
}

func (backend *BasicAuthenticationBackend) GetUserID(ctx context.Context, first, second string) (int, *uuid.UUID) {
	if functools.NormalizeEmail(first) != "" {
		if len(second) == 6 {
			status, user := backend.getUserFromEmailAndCode(ctx, first, second)
			if status == enums.Ok && user != nil {
				return status, user
			}
		}

		return backend.getUserFromEmailAndPassword(ctx, first, second)
	} else {
		if len(second) == 6 {
			status, user := backend.getUserFromPhoneAndCode(ctx, first, second)
			if status == enums.Ok && user != nil {
				return status, user
			}
		}

		return backend.getUserFromPhoneAndPassword(ctx, first, second)
	}
}

func (backend *BasicAuthenticationBackend) GetUser(ctx context.Context, token string) (int, models.IAuthenticationBackendUser) {
	decodedTokenInBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		sentry.CaptureException(err)
		return enums.IncorrectToken, nil
	}

	decodedToken := string(decodedTokenInBytes)
	parts := strings.Split(decodedToken, ":")
	if len(parts) < 2 {
		return enums.IncorrectToken, nil
	}

	status, userID := backend.GetUserID(ctx, parts[0], parts[1])
	if status != enums.Ok {
		return status, nil
	}

	user := backend.store.GetUserView(ctx, *userID)
	if user == nil {
		return enums.UserNotFound, nil
	}

	return enums.Ok, &BasicAuthenticationBackendUser{
		IsAdmin: functools.Contains(backend.environment.AdminRole, user.Roles),
		Roles:   user.Roles,
		UserID:  user.Id,
	}
}
