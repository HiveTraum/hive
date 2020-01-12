package controllers

import (
	"auth/enums"
	"auth/mocks"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLoginController_EncodePassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	require.NotEmpty(t, encodedPassword)
	require.NotEqual(t, password, encodedPassword)
	require.Len(t, encodedPassword, 60)
}

func TestLoginController_VerifyPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	isVerified := controller.VerifyPassword(ctx, password, encodedPassword)
	require.True(t, isVerified)
}

func TestLoginController_NormalizeEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	email := "mail@mail.com"
	normalizedEmail := controller.NormalizeEmail(ctx, email)
	require.NotEmpty(t, normalizedEmail)
	require.Equal(t, email, normalizedEmail)
}

func TestLoginController_NormalizeIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	incorrectEmail := "mail"
	normalizedEmail := controller.NormalizeEmail(ctx, incorrectEmail)
	require.Empty(t, normalizedEmail)
	require.NotEqual(t, incorrectEmail, normalizedEmail)
	require.Equal(t, "", normalizedEmail)
}

func TestLoginController_NormalizePhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	phone := "+71234567890"
	normalizedPhone := controller.NormalizePhone(ctx, phone)
	require.NotEmpty(t, normalizedPhone)
	require.Equal(t, phone, normalizedPhone)
}

func TestLoginController_NormalizeIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	phone := "+123"
	normalizedPhone := controller.NormalizePhone(ctx, phone)
	require.Empty(t, normalizedPhone)
	require.NotEqual(t, phone, normalizedPhone)
	require.Equal(t, "", normalizedPhone)
}

func TestLoginController_EncodeAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	token := controller.EncodeAccessToken(ctx, 1, []string{"admin"}, "123", time.Now().Add(time.Millisecond))
	require.NotEmpty(t, token)
}

func TestLoginController_DecodeAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	encodedToken := controller.EncodeAccessToken(ctx, 1, []string{"admin"}, "123", time.Now().Add(time.Millisecond))
	status, decodedToken := controller.DecodeAccessToken(ctx, encodedToken, "123")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, decodedToken)
	require.Equal(t, models.UserID(1), decodedToken.UserID)
	require.Contains(t, decodedToken.Roles, "admin")
}

func TestLoginController_DecodeIncorrectAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	status, decodedToken := controller.DecodeAccessToken(ctx, "123", "123")
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, decodedToken)
}

func TestLoginController_DecodeExpiredAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	encodedToken := controller.EncodeAccessToken(ctx, 1, []string{"admin"}, "123", time.Now().Add(time.Second))
	time.Sleep(time.Second * 2)
	status, decodedToken := controller.DecodeAccessToken(ctx, encodedToken, "123")
	require.Equal(t, enums.InvalidToken, status)
	require.Nil(t, decodedToken)
}

func TestLoginController_DecodeAccessTokenWithoutValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	encodedToken := controller.EncodeAccessToken(ctx, 1, []string{"admin"}, "123", time.Now().Add(time.Millisecond))
	status, decodedToken := controller.DecodeAccessTokenWithoutValidation(ctx, encodedToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, decodedToken)
	require.Equal(t, models.UserID(1), decodedToken.UserID)
	require.Contains(t, decodedToken.Roles, "admin")
}

func TestLoginController_DecodeIncorrectAccessTokenWithoutValidation(t *testing.T) {

	// Декодирование без валидации должно проводить базовую верификацию токена на корректность, например структура токена

	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	status, decodedToken := controller.DecodeAccessTokenWithoutValidation(ctx, "123")
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, decodedToken)
}

func TestLoginController_DecodeExpiredAccessTokenWithoutValidation(t *testing.T) {

	// Декодирование без валидации должно корректно декодировать токены с истекшим сроком

	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}
	encodedToken := controller.EncodeAccessToken(ctx, 1, []string{"admin"}, "123", time.Now().Add(time.Second))
	time.Sleep(time.Second * 2)
	status, decodedToken := controller.DecodeAccessTokenWithoutValidation(ctx, encodedToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, decodedToken)
	require.Equal(t, models.UserID(1), decodedToken.UserID)
	require.Contains(t, decodedToken.Roles, "admin")
}

func TestLoginController_LoginByTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	refreshToken := "refreshToken"
	fingerprint := "fingerprint"
	userID := models.UserID(1)
	secretID := models.SecretID(1)
	secretValue := "123"

	store.
		EXPECT().
		GetSession(ctx, fingerprint, refreshToken, userID).
		Times(1).
		Return(&models.Session{
			RefreshToken: refreshToken,
			Fingerprint:  fingerprint,
			UserID:       userID,
			SecretID:     secretID,
			Created:      1,
			UserAgent:    "chrome",
		})

	store.
		EXPECT().
		GetSecret(ctx, models.SecretID(1)).
		Times(1).
		Return(&models.Secret{
			Id:      secretID,
			Created: 1,
			Value:   secretValue,
		})

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(&models.User{
			Id:      userID,
			Created: 1,
		})

	encodedToken := controller.EncodeAccessToken(ctx, userID, []string{"admin"}, secretValue, time.Now().Add(time.Second))
	status, user := controller.LoginByTokens(ctx, refreshToken, encodedToken, fingerprint)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, user)
}

func TestLoginController_LoginByIncorrectTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}

	refreshToken := "refreshToken"
	fingerprint := "fingerprint"

	status, user := controller.LoginByTokens(ctx, refreshToken, "123", fingerprint)
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByTokensWithoutSession(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	refreshToken := "refreshToken"
	fingerprint := "fingerprint"
	userID := models.UserID(1)
	secretValue := "123"

	store.
		EXPECT().
		GetSession(ctx, fingerprint, refreshToken, userID).
		Times(1).
		Return(nil)

	encodedToken := controller.EncodeAccessToken(ctx, userID, []string{"admin"}, secretValue, time.Now().Add(time.Second))
	status, user := controller.LoginByTokens(ctx, refreshToken, encodedToken, fingerprint)
	require.Equal(t, enums.SessionNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByTokensWithoutSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	refreshToken := "refreshToken"
	fingerprint := "fingerprint"
	userID := models.UserID(1)
	secretID := models.SecretID(1)
	secretValue := "123"

	store.
		EXPECT().
		GetSession(ctx, fingerprint, refreshToken, userID).
		Times(1).
		Return(&models.Session{
			RefreshToken: refreshToken,
			Fingerprint:  fingerprint,
			UserID:       userID,
			SecretID:     secretID,
			Created:      1,
			UserAgent:    "chrome",
		})

	store.
		EXPECT().
		GetSecret(ctx, models.SecretID(1)).
		Times(1).
		Return(nil)

	encodedToken := controller.EncodeAccessToken(ctx, userID, []string{"admin"}, secretValue, time.Now().Add(time.Second))
	status, user := controller.LoginByTokens(ctx, refreshToken, encodedToken, fingerprint)
	require.Equal(t, enums.SecretNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByTokensWithoutUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	refreshToken := "refreshToken"
	fingerprint := "fingerprint"
	userID := models.UserID(1)
	secretID := models.SecretID(1)
	secretValue := "123"

	store.
		EXPECT().
		GetSession(ctx, fingerprint, refreshToken, userID).
		Times(1).
		Return(&models.Session{
			RefreshToken: refreshToken,
			Fingerprint:  fingerprint,
			UserID:       userID,
			SecretID:     secretID,
			Created:      1,
			UserAgent:    "chrome",
		})

	store.
		EXPECT().
		GetSecret(ctx, models.SecretID(1)).
		Times(1).
		Return(&models.Secret{
			Id:      secretID,
			Created: 1,
			Value:   secretValue,
		})

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(nil)

	encodedToken := controller.EncodeAccessToken(ctx, userID, []string{"admin"}, secretValue, time.Now().Add(time.Second))
	status, user := controller.LoginByTokens(ctx, refreshToken, encodedToken, fingerprint)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	email := "mail@mail.com"
	emailCode := "1234"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(emailCode)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      models.EmailID(1),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      models.PasswordID(1),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(&models.User{
			Id:      userID,
			Created: 1,
		})

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, user)
}

func TestLoginController_LoginByEmailWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}

	password := "123"
	email := "mail"
	emailCode := "1234"

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailWithoutEmailConfirmationCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	email := "mail@mail.com"
	emailCode := "1234"

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("")

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.EmailConfirmationCodeNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailWithIncorrectEmailConfirmationCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	email := "mail@mail.com"

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("4321")

	status, user := controller.LoginByEmail(ctx, email, "1234", password)
	require.Equal(t, enums.IncorrectEmailCode, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	email := "mail@mail.com"
	emailCode := "1234"

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(emailCode)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, nil)

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	email := "mail@mail.com"
	emailCode := "1234"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(emailCode)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      models.EmailID(1),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, nil)

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, "321")
	email := "mail@mail.com"
	emailCode := "1234"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(emailCode)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      models.EmailID(1),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      models.PasswordID(1),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailWithoutUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	email := "mail@mail.com"
	emailCode := "1234"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(emailCode)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      models.EmailID(1),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      models.PasswordID(1),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(nil)

	status, user := controller.LoginByEmail(ctx, email, emailCode, password)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, user)
}
