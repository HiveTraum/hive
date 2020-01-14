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

// Email and password

func TestLoginController_LoginByEmailAndPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	email := "mail@mail.com"
	userID := models.UserID(1)

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

	status, user := controller.LoginByEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, user)
}

func TestLoginController_LoginByEmailAndPasswordWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}

	status, user := controller.LoginByEmailAndPassword(ctx, "mail", "password")
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndPasswordWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	email := "mail@mail.com"
	password := "1234"

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, nil)

	status, user := controller.LoginByEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndPasswordWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	email := "mail@mail.com"
	userID := models.UserID(1)

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

	status, user := controller.LoginByEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndPasswordWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, "321")
	email := "mail@mail.com"
	userID := models.UserID(1)

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

	status, user := controller.LoginByEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndPasswordWithoutUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	email := "mail@mail.com"
	userID := models.UserID(1)

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

	status, user := controller.LoginByEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, user)
}

// Email and code

func TestLoginController_LoginByEmailAndCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	code := "123"
	email := "mail@mail.com"
	userID := models.UserID(1)

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
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(code)

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(&models.User{
			Id:      userID,
			Created: 1,
		})

	status, user := controller.LoginByEmailAndCode(ctx, email, code)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, user)
}

func TestLoginController_LoginByEmailAndCodeWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}

	email := "mail"
	emailCode := "1234"

	status, user := controller.LoginByEmailAndCode(ctx, email, emailCode)
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndCodeWithoutEmailConfirmationCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	email := "mail@mail.com"
	emailCode := "1234"

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("")

	status, user := controller.LoginByEmailAndCode(ctx, email, emailCode)
	require.Equal(t, enums.EmailConfirmationCodeNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndCodeWithIncorrectEmailConfirmationCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	email := "mail@mail.com"

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("4321")

	status, user := controller.LoginByEmailAndCode(ctx, email, "1234")
	require.Equal(t, enums.IncorrectEmailCode, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndCodeWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

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

	status, user := controller.LoginByEmailAndCode(ctx, email, emailCode)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByEmailAndCodeWithoutUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	code := "123"
	email := "mail@mail.com"
	userID := models.UserID(1)

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
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(code)

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(nil)

	status, user := controller.LoginByEmailAndCode(ctx, email, code)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, user)
}

// Phone and password

func TestLoginController_LoginByPhoneAndPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	phone := "+71234567890"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      models.PhoneID(1),
			Created: 1,
			UserId:  userID,
			Value:   phone,
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

	status, user := controller.LoginByPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, user)
}

func TestLoginController_LoginByPhoneAndPasswordWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}

	status, user := controller.LoginByPhoneAndPassword(ctx, "123", "password")
	require.Equal(t, enums.IncorrectPhone, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndPasswordWithoutPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	phone := "+71234567890"
	password := "1234"

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, nil)

	status, user := controller.LoginByPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.PhoneNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndPasswordWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	phone := "+71234567890"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      models.PhoneID(1),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, nil)

	status, user := controller.LoginByPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndPasswordWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, "321")
	phone := "+71234567890"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      models.PhoneID(1),
			Created: 1,
			UserId:  userID,
			Value:   phone,
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

	status, user := controller.LoginByPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndPasswordWithoutUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	password := "123"
	encodedPassword := controller.EncodePassword(ctx, password)
	phone := "+71234567890"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      models.PhoneID(1),
			Created: 1,
			UserId:  userID,
			Value:   phone,
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

	status, user := controller.LoginByPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, user)
}

// Phone and code

func TestLoginController_LoginByPhoneAndCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	code := "123"
	phone := "+71234567890"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      models.PhoneID(1),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return(code)

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(&models.User{
			Id:      userID,
			Created: 1,
		})

	status, user := controller.LoginByPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, user)
}

func TestLoginController_LoginByPhoneAndCodeWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	controller := LoginController{Store: nil}

	phone := "1234"
	code := "1234"

	status, user := controller.LoginByPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.IncorrectPhone, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndCodeWithoutPhoneConfirmationCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	phone := "+71234567890"
	code := "1234"

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return("")

	status, user := controller.LoginByPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.PhoneConfirmationCodeNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndCodeWithIncorrectPhoneConfirmationCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	phone := "+71234567890"

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return("4321")

	status, user := controller.LoginByPhoneAndCode(ctx, phone, "1234")
	require.Equal(t, enums.IncorrectPhoneCode, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndCodeWithoutPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	phone := "+71234567890"
	code := "1234"

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return(code)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, nil)

	status, user := controller.LoginByPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.PhoneNotFound, status)
	require.Nil(t, user)
}

func TestLoginController_LoginByPhoneAndCodeWithoutUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, _ := mocks.InitMockApp(ctrl)
	controller := LoginController{Store: store}

	code := "123"
	phone := "+71234567890"
	userID := models.UserID(1)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      models.PhoneID(1),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return(code)

	store.
		EXPECT().
		GetUser(ctx, userID).
		Times(1).
		Return(nil)

	status, user := controller.LoginByPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, user)
}