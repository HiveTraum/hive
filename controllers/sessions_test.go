package controllers

import (
	"auth/enums"
	"auth/mocks"
	"auth/models"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateSessionFromTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	userID := uuid.NewV4()
	fingerprint := "123"
	refreshToken := "321"
	accessToken := "123321"

	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}

	loginController.
		EXPECT().
		Login(ctx, accessToken).
		Times(1).
		Return(enums.Ok, &models.AccessTokenPayload{
			StandardClaims: jwt.StandardClaims{},
			IsAdmin:        false,
			Roles:          nil,
			UserID:         userID,
			SecretID:       secret.Id,
		})

	store.
		EXPECT().
		GetSession(ctx, fingerprint, refreshToken, userID).
		Times(1).
		Return(&models.Session{
			RefreshToken: refreshToken,
			Fingerprint:  fingerprint,
			UserID:       userID,
			SecretID:     secret.Id,
			Created:      1,
			UserAgent:    "chrome",
			AccessToken:  accessToken,
		})

	status, loggedUserID := getUserFromTokens(store, loginController, ctx, accessToken, fingerprint, refreshToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, *loggedUserID)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromTokensWithIncorrectAuth(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	fingerprint := "123"
	refreshToken := "321"
	accessToken := "123321"

	loginController.
		EXPECT().
		Login(ctx, accessToken).
		Times(1).
		Return(enums.IncorrectToken, nil)

	status, loggedUserID := getUserFromTokens(store, loginController, ctx, accessToken, fingerprint, refreshToken)
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromTokensWithoutSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	fingerprint := "123"
	refreshToken := "321"
	accessToken := "123321"

	loginController.
		EXPECT().
		Login(ctx, accessToken).
		Times(1).
		Return(enums.SecretNotFound, nil)

	status, loggedUserID := getUserFromTokens(store, loginController, ctx, accessToken, fingerprint, refreshToken)
	require.Equal(t, enums.SecretNotFound, status)
	require.Nil(t, loggedUserID)
}

// Email and password

func TestCreateSessionFromEmailAndPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	encodedPassword := "321"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	loginController.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(true)

	status, loggedUserID := getUserFromEmailAndPassword(store, loginController, ctx, email, password)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	email := "mail@mail.com"

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return("")

	status, loggedUserID := getUserFromEmailAndPassword(store, loginController, ctx, email, password)
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	email := "mail@mail.com"

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := getUserFromEmailAndPassword(store, loginController, ctx, email, password)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := getUserFromEmailAndPassword(store, loginController, ctx, email, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	encodedPassword := "321"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	loginController.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(false)

	status, loggedUserID := getUserFromEmailAndPassword(store, loginController, ctx, email, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, loggedUserID)
}

// Email and code

func TestCreateSessionFromEmailAndCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(code)

	status, loggedUserID := getUserFromEmailAndCode(store, loginController, ctx, email, code)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	email := "mail@mail.com"

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return("")

	status, loggedUserID := getUserFromEmailAndCode(store, loginController, ctx, email, code)
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithoutCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	email := "mail@mail.com"

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("")

	status, loggedUserID := getUserFromEmailAndCode(store, loginController, ctx, email, code)
	require.Equal(t, enums.EmailConfirmationCodeNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithIncorrectCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	email := "mail@mail.com"

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("321")

	status, loggedUserID := getUserFromEmailAndCode(store, loginController, ctx, email, code)
	require.Equal(t, enums.IncorrectEmailCode, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	email := "mail@mail.com"

	loginController.
		EXPECT().
		NormalizeEmail(ctx, email).
		Times(1).
		Return(email)

	store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(code)

	store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := getUserFromEmailAndCode(store, loginController, ctx, email, code)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, loggedUserID)
}

// Phone and password

func TestCreateSessionFromPhoneAndPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	encodedPassword := "321"
	phone := "+71234567890"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	loginController.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(true)

	status, loggedUserID := getUserFromPhoneAndPassword(store, loginController, ctx, phone, password)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	phone := "phone"

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return("")

	status, loggedUserID := getUserFromPhoneAndPassword(store, loginController, ctx, phone, password)
	require.Equal(t, enums.IncorrectPhone, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	phone := "+71234567890"

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := getUserFromPhoneAndPassword(store, loginController, ctx, phone, password)
	require.Equal(t, enums.PhoneNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	phone := "+71234567890"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := getUserFromPhoneAndPassword(store, loginController, ctx, phone, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	password := "123"
	encodedPassword := "321"
	phone := "+71234567890"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	loginController.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(false)

	status, loggedUserID := getUserFromPhoneAndPassword(store, loginController, ctx, phone, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, loggedUserID)
}

// Phone and code

func TestCreateSessionFromPhoneAndCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	phone := "+71234567890"
	userID := uuid.NewV4()

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhone(ctx, phone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return(code)

	status, loggedUserID := getUserFromPhoneAndCode(store, loginController, ctx, phone, code)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	phone := "+71234567890"

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return("")

	status, loggedUserID := getUserFromPhoneAndCode(store, loginController, ctx, phone, code)
	require.Equal(t, enums.IncorrectPhone, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithoutCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	phone := "+71234567890"

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return("")

	status, loggedUserID := getUserFromPhoneAndCode(store, loginController, ctx, phone, code)
	require.Equal(t, enums.PhoneConfirmationCodeNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithIncorrectCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	phone := "+71234567890"

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, phone).
		Times(1).
		Return("321")

	status, loggedUserID := getUserFromPhoneAndCode(store, loginController, ctx, phone, code)
	require.Equal(t, enums.IncorrectPhoneCode, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithoutPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, _, loginController := mocks.InitMockApp(ctrl)

	code := "123"
	phone := "+71234567890"

	loginController.
		EXPECT().
		NormalizePhone(ctx, phone).
		Times(1).
		Return(phone)

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

	status, loggedUserID := getUserFromPhoneAndCode(store, loginController, ctx, phone, code)
	require.Equal(t, enums.PhoneNotFound, status)
	require.Nil(t, loggedUserID)
}