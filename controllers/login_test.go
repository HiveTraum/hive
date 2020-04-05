package controllers

import (
	"auth/backends"
	"auth/enums"
	"auth/infrastructure"
	"auth/mocks"
	"auth/models"
	"auth/processors"
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLoginController_EncodePassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	processor := processors.PasswordProcessor{}
	password := "123"
	encodedPassword := processor.EncodePassword(ctx, password)
	require.NotEmpty(t, encodedPassword)
	require.NotEqual(t, password, encodedPassword)
	require.Len(t, encodedPassword, 60)
}

func TestLoginController_VerifyPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	processor := processors.PasswordProcessor{}
	password := "123"
	encodedPassword := processor.EncodePassword(ctx, password)
	isVerified := processor.VerifyPassword(ctx, password, encodedPassword)
	require.True(t, isVerified)
}

func TestLoginController_EncodeAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	token := backend.EncodeAccessToken(ctx, uuid.NewV4(), []string{"admin"}, secret, time.Now().Add(time.Millisecond))
	require.NotEmpty(t, token)
}

func TestLoginController_DecodeAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	userID := uuid.NewV4()
	encodedToken := backend.EncodeAccessToken(ctx, userID, []string{"admin"}, secret, time.Now().Add(time.Millisecond))
	status, decodedToken := backend.DecodeAccessToken(ctx, encodedToken, secret.Value)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, decodedToken)
	require.Equal(t, userID, decodedToken.UserID)
	require.Contains(t, decodedToken.Roles, "admin")
	require.Equal(t, secret.Id, decodedToken.SecretID)
}

func TestLoginController_DecodeIncorrectAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	status, decodedToken := backend.DecodeAccessToken(ctx, "123", uuid.NewV4())
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, decodedToken)
}

func TestLoginController_DecodeExpiredAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	encodedToken := backend.EncodeAccessToken(ctx, uuid.NewV4(), []string{"admin"}, secret, time.Now().Add(time.Second))
	time.Sleep(time.Second * 2)
	status, decodedToken := backend.DecodeAccessToken(ctx, encodedToken, secret.Value)
	require.Equal(t, enums.InvalidToken, status)
	require.Nil(t, decodedToken)
}

func TestLoginController_DecodeAccessTokenWithoutValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	userID := uuid.NewV4()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}
	encodedToken := backend.EncodeAccessToken(ctx, userID, []string{"admin"}, secret, time.Now().Add(time.Millisecond))
	status, decodedToken := backend.DecodeAccessTokenWithoutValidation(ctx, encodedToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, decodedToken)
	require.Equal(t, userID, decodedToken.UserID)
	require.Contains(t, decodedToken.Roles, "admin")
}

func TestLoginController_DecodeIncorrectAccessTokenWithoutValidation(t *testing.T) {

	// Декодирование без валидации должно проводить базовую верификацию токена на корректность, например структура токена

	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	status, decodedToken := backend.DecodeAccessTokenWithoutValidation(ctx, "123")
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, decodedToken)
}

func TestLoginController_DecodeExpiredAccessTokenWithoutValidation(t *testing.T) {

	// Декодирование без валидации должно корректно декодировать токены с истекшим сроком

	t.Parallel()
	ctx := context.Background()
	backend := backends.JWTAuthenticationBackend{}
	userID := uuid.NewV4()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}
	encodedToken := backend.EncodeAccessToken(ctx, userID, []string{"admin"}, secret, time.Now().Add(time.Second))
	time.Sleep(time.Second * 2)
	status, decodedToken := backend.DecodeAccessTokenWithoutValidation(ctx, encodedToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, decodedToken)
	require.Equal(t, userID, decodedToken.UserID)
	require.Contains(t, decodedToken.Roles, "admin")
}

func TestLoginController_LoginByTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _, _ := mocks.InitMockApp(ctrl)
	backend := backends.JWTAuthenticationBackend{Store: store}
	controller := LoginController{Backends: map[string]infrastructure.AuthenticationBackend{"Bearer": backend}}

	userID := uuid.NewV4()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}

	store.
		EXPECT().
		GetSecret(ctx, secret.Id).
		Times(1).
		Return(&models.Secret{
			Id:      secret.Id,
			Created: 1,
			Value:   secret.Value,
		})

	encodedToken := "Bearer " + backend.EncodeAccessToken(ctx, userID, []string{"admin"}, secret, time.Now().Add(time.Second))
	status, payload := controller.Login(ctx, encodedToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, payload)
	require.True(t, payload.GetIsAdmin())
	require.Equal(t, userID, payload.GetUserID())
	require.Equal(t, []string{"admin"}, payload.GetRoles())
	// Todo
	//require.Equal(t, secret.Id, payload.SecretID)
}

func TestLoginController_LoginByIncorrectTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Backends: map[string]infrastructure.AuthenticationBackend{"Bearer": backends.JWTAuthenticationBackend{}}}

	status, payload := controller.Login(ctx, "Bearer 123")
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, payload)
}

func TestLoginController_LoginByTokensWithoutSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _, _ := mocks.InitMockApp(ctrl)
	backend := backends.JWTAuthenticationBackend{Store: store}
	controller := LoginController{Backends: map[string]infrastructure.AuthenticationBackend{"Bearer": backend}}

	userID := uuid.NewV4()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}

	store.
		EXPECT().
		GetSecret(ctx, secret.Id).
		Times(1).
		Return(nil)

	encodedToken := backend.EncodeAccessToken(ctx, userID, []string{"admin"}, secret, time.Now().Add(time.Second))
	status, payload := controller.Login(ctx, "Bearer "+encodedToken)
	require.Equal(t, enums.SecretNotFound, status)
	require.Nil(t, payload)
}
