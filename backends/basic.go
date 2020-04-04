package backends

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"context"
	"encoding/base64"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type BasicAuthenticationBackendUser struct {
	IsAdmin bool
	Roles   []string
	UserID  uuid.UUID
}

type BasicAuthenticationBackend struct {
	Store infrastructure.StoreInterface
}

func (user BasicAuthenticationBackendUser) GetIsAdmin() bool {
	return user.IsAdmin
}

func (user BasicAuthenticationBackendUser) GetRoles() []string {
	return user.Roles
}

func (user BasicAuthenticationBackendUser) GetUserID() uuid.UUID {
	return user.UserID
}

func (backend BasicAuthenticationBackend) GetUserID(ctx context.Context, credential string) (int, *uuid.UUID) {
	email := functools.NormalizeEmail(credential)
	if email != "" {
		status, emailObject := backend.Store.GetEmail(ctx, email)
		return status, &emailObject.UserId
	}

	phoneCredential := strings.Split(credential, "/")
	if len(phoneCredential) < 2 {
		phone := functools.NormalizePhone(phoneCredential[1], phoneCredential[0])
		status, phoneObject := backend.Store.GetPhone(ctx, phone)
		return status, &phoneObject.UserId
	}

	return enums.UserNotFound, nil
}

func (backend BasicAuthenticationBackend) GetUser(ctx context.Context, token string) (int, infrastructure.AuthenticationBackendUser) {
	decodedTokenInBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		sentry.CaptureException(err)
		return enums.IncorrectToken, nil
	}

	decodedToken := string(decodedTokenInBytes)
	credentials := strings.Split(decodedToken, ":")
	if len(credentials) < 2 {
		return enums.IncorrectToken, nil
	}

	credential, password := credentials[0], credentials[1]
	if credential == "" || password == "" {
		return enums.IncorrectToken, nil
	}

	status, userID := backend.GetUserID(ctx, credential)
	if status != enums.Ok {
		return status, nil
	}

	user := backend.Store.GetUserView(ctx, *userID)
	if user == nil {
		return enums.UserNotFound, nil
	}

	return enums.Ok, BasicAuthenticationBackendUser{
		IsAdmin: functools.Contains(config.GetEnvironment().AdminRole, user.Roles),
		Roles:   user.Roles,
		UserID:  user.Id,
	}
}
