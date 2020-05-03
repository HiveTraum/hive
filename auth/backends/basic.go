package backends

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"auth/models"
	"auth/passwordProcessors"
	"auth/stores"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
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
}

func InitBasicAuthenticationBackend(store stores.IStore, passwordProcessor passwordProcessors.IPasswordProcessor) *BasicAuthenticationBackend {
	return &BasicAuthenticationBackend{
		store:             store,
		passwordProcessor: passwordProcessor,
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
		Backend:           InitBasicAuthenticationBackend(store, passwordProcessor),
		Store:             store,
		PasswordProcessor: passwordProcessor,
	}
}

type CredentialsPair struct {
	CredentialType   enums.CredentialType `json:"credentialType"`
	FirstCredential  string               `json:"first"`
	SecondCredential string               `json:"second"`
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

func (backend *BasicAuthenticationBackend) GetUserID(ctx context.Context, credentials CredentialsPair) (int, *uuid.UUID) {
	switch credentials.CredentialType {
	case enums.EmailAndPassword:
		return backend.getUserFromEmailAndPassword(ctx, credentials.FirstCredential, credentials.SecondCredential)
	case enums.EmailAndCode:
		return backend.getUserFromEmailAndCode(ctx, credentials.FirstCredential, credentials.SecondCredential)
	case enums.PhoneAndPassword:
		return backend.getUserFromPhoneAndPassword(ctx, credentials.FirstCredential, credentials.SecondCredential)
	case enums.PhoneAndCode:
		return backend.getUserFromPhoneAndCode(ctx, credentials.FirstCredential, credentials.SecondCredential)
	default:
		return enums.CredentialsTypeNotFound, nil
	}
}

func (backend *BasicAuthenticationBackend) GetUser(ctx context.Context, headerToken string, _ string) (int, models.IAuthenticationBackendUser) {
	decodedTokenInBytes, err := base64.StdEncoding.DecodeString(headerToken)
	if err != nil {
		sentry.CaptureException(err)
		return enums.IncorrectToken, nil
	}

	var credentials CredentialsPair
	err = json.Unmarshal(decodedTokenInBytes, &credentials)
	if err != nil {
		return enums.IncorrectToken, nil
	}

	status, userID := backend.GetUserID(ctx, credentials)
	if status != enums.Ok {
		return status, nil
	}

	user := backend.store.GetUserView(ctx, *userID)
	if user == nil {
		return enums.UserNotFound, nil
	}

	return enums.Ok, &BasicAuthenticationBackendUser{
		IsAdmin: functools.Contains(config.GetEnvironment().AdminRole, user.Roles),
		Roles:   user.Roles,
		UserID:  user.Id,
	}
}
