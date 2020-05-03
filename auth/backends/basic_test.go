package backends

import (
	"auth/enums"
	"auth/functools"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

// Email and password

func TestCreateSessionFromEmailAndPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	encodedPassword := "321"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	backend.
		Store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	backend.
		Store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	backend.
		PasswordProcessor.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(true)

	status, loggedUserID := backend.Backend.getUserFromEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	email := "mail"

	status, loggedUserID := backend.Backend.getUserFromEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	email := "mail@mail.com"

	backend.
		Store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := backend.Backend.getUserFromEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	backend.
		Store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	backend.Store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := backend.Backend.getUserFromEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndPasswordWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	encodedPassword := "321"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	backend.Store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	backend.Store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	backend.PasswordProcessor.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(false)

	status, loggedUserID := backend.Backend.getUserFromEmailAndPassword(ctx, email, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, loggedUserID)
}

// Email and code

func TestCreateSessionFromEmailAndCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	email := "mail@mail.com"
	userID := uuid.NewV4()

	backend.Store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, &models.Email{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   email,
		})

	backend.Store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(code)

	status, loggedUserID := backend.Backend.getUserFromEmailAndCode(ctx, email, code)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	email := "mail"

	status, loggedUserID := backend.Backend.getUserFromEmailAndCode(ctx, email, code)
	require.Equal(t, enums.IncorrectEmail, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithoutCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	email := "mail@mail.com"

	backend.Store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("")

	status, loggedUserID := backend.Backend.getUserFromEmailAndCode(ctx, email, code)
	require.Equal(t, enums.EmailConfirmationCodeNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithIncorrectCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	email := "mail@mail.com"
	backend.Store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return("321")

	status, loggedUserID := backend.Backend.getUserFromEmailAndCode(ctx, email, code)
	require.Equal(t, enums.IncorrectEmailCode, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromEmailAndCodeWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	email := "mail@mail.com"

	backend.Store.
		EXPECT().
		GetEmailConfirmationCode(ctx, email).
		Times(1).
		Return(code)

	backend.Store.
		EXPECT().
		GetEmail(ctx, email).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := backend.Backend.getUserFromEmailAndCode(ctx, email, code)
	require.Equal(t, enums.EmailNotFound, status)
	require.Nil(t, loggedUserID)
}

// Phone and password

func TestCreateSessionFromPhoneAndPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	encodedPassword := "321"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)
	userID := uuid.NewV4()

	backend.Store.
		EXPECT().
		GetPhone(ctx, formattedPhone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   formattedPhone,
		})

	backend.Store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	backend.PasswordProcessor.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(true)

	status, loggedUserID := backend.Backend.getUserFromPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithIncorrectEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	phone := "phone"

	status, loggedUserID := backend.Backend.getUserFromPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.IncorrectPhone, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithoutEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)

	backend.Store.
		EXPECT().
		GetPhone(ctx, formattedPhone).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := backend.Backend.getUserFromPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.PhoneNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithoutPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)
	userID := uuid.NewV4()

	backend.Store.
		EXPECT().
		GetPhone(ctx, formattedPhone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   phone,
		})

	backend.Store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := backend.Backend.getUserFromPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.PasswordNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndPasswordWithIncorrectPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	password := "123"
	encodedPassword := "321"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)
	userID := uuid.NewV4()

	backend.Store.
		EXPECT().
		GetPhone(ctx, formattedPhone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   formattedPhone,
		})

	backend.
		Store.
		EXPECT().
		GetLatestPassword(ctx, userID).
		Times(1).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   encodedPassword,
		})

	backend.
		PasswordProcessor.
		EXPECT().
		VerifyPassword(ctx, password, encodedPassword).
		Times(1).
		Return(false)

	status, loggedUserID := backend.Backend.getUserFromPhoneAndPassword(ctx, phone, password)
	require.Equal(t, enums.IncorrectPassword, status)
	require.Nil(t, loggedUserID)
}

// Phone and code

func TestCreateSessionFromPhoneAndCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)
	userID := uuid.NewV4()

	backend.
		Store.
		EXPECT().
		GetPhone(ctx, formattedPhone).
		Times(1).
		Return(enums.Ok, &models.Phone{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   formattedPhone,
		})

	backend.
		Store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, formattedPhone).
		Times(1).
		Return(code)

	status, loggedUserID := backend.Backend.getUserFromPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, userID, *loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	phone := "1"

	status, loggedUserID := backend.Backend.getUserFromPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.IncorrectPhone, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithoutCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)

	backend.Store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, formattedPhone).
		Times(1).
		Return("")

	status, loggedUserID := backend.Backend.getUserFromPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.PhoneConfirmationCodeNotFound, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithIncorrectCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)

	backend.Store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, formattedPhone).
		Times(1).
		Return("321")

	status, loggedUserID := backend.Backend.getUserFromPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.IncorrectPhoneCode, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromPhoneAndCodeWithoutPhone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitBasicAuthenticationWithMockedInternals(ctrl)

	code := "123"
	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)

	backend.Store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, formattedPhone).
		Times(1).
		Return(code)

	backend.Store.
		EXPECT().
		GetPhone(ctx, formattedPhone).
		Times(1).
		Return(enums.Ok, nil)

	status, loggedUserID := backend.Backend.getUserFromPhoneAndCode(ctx, phone, code)
	require.Equal(t, enums.PhoneNotFound, status)
	require.Nil(t, loggedUserID)
}
