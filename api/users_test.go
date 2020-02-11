package api

import (
	"auth/enums"
	"auth/functools"
	"auth/inout"
	"auth/mocks"
	"auth/models"
	"auth/repositories"
	"bytes"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserEmptyBody(t *testing.T) {
	t.Parallel()
	body := "{}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, _, _ := mocks.InitMockApp(ctrl)

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))

	store.
		EXPECT().
		CreateUser(r.Context(), gomock.Any()).
		Times(0)

	r.Header.Add("Content-Type", "application/json")
	status, message := createUserV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	v, ok := message.(*inout.CreateUserBadRequestV1)
	require.True(t, ok)
	require.Len(t, v.Password, 1)
}

func TestCreateUserWithOnlyPassword(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, _, _ := mocks.InitMockApp(ctrl)

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))

	store.
		EXPECT().
		CreateUser(r.Context(), "olleh").
		Times(0)

	r.Header.Add("Content-Type", "application/json")
	status, message := createUserV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	v, ok := message.(*inout.CreateUserBadRequestV1)
	require.True(t, ok)
	require.Len(t, v.Errors, 1)
}

func TestCreateUserWithOnlyEmail(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))

	app, store, _, passwordProcessor := mocks.InitMockApp(ctrl)

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hello").
		Return("olleh")

	store.
		EXPECT().
		CreateUser(r.Context(), "olleh").
		Times(0)

	store.
		EXPECT().
		GetEmailConfirmationCode(r.Context(), "mail@mail.com").
		Times(1).
		Return("")

	r.Header.Add("Content-Type", "application/json")
	status, message := createUserV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	v, ok := message.(*inout.CreateUserBadRequestV1)
	require.True(t, ok)
	require.Len(t, v.Email, 1)
}

func TestCreateUserWithEmailAndEmailCode(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\", \"emailCode\": \"123456\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))

	app, store, _, passwordProcessor := mocks.InitMockApp(ctrl)

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hello").
		Return("olleh")

	store.
		EXPECT().
		CreateUser(r.Context(), gomock.Any()).
		Times(0)

	store.
		EXPECT().
		GetEmailConfirmationCode(r.Context(), "mail@mail.com").
		Times(1).
		Return("")

	r.Header.Add("Content-Type", "application/json")
	status, message := createUserV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	v, ok := message.(*inout.CreateUserBadRequestV1)
	require.True(t, ok)
	require.Len(t, v.Email, 1)
}

func TestCreateUserWithEmailAndEmailCodeAfterIncorrectEmailConfirmationCodeReceived(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\", \"emailCode\": \"123456\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))

	app, store, _, passwordProcessor := mocks.InitMockApp(ctrl)

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hello").
		Return("olleh")

	store.
		EXPECT().
		CreateUser(r.Context(), gomock.Any()).
		Times(0)

	store.
		EXPECT().
		GetEmailConfirmationCode(r.Context(), "mail@mail.com").
		Times(1).
		Return("654321")

	r.Header.Add("Content-Type", "application/json")
	status, message := createUserV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	v, ok := message.(*inout.CreateUserBadRequestV1)
	require.True(t, ok)
	require.Len(t, v.EmailCode, 1)
}

func TestSuccessfulCreateUserWithEmail(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\", \"emailCode\": \"123456\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))

	app, store, esb, passwordProcessor := mocks.InitMockApp(ctrl)

	userID := uuid.NewV4()

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hello").
		Return("olleh")

	store.
		EXPECT().
		CreateUser(r.Context(), gomock.Any()).
		Times(0)

	store.
		EXPECT().
		GetEmailConfirmationCode(r.Context(), "mail@mail.com").
		Times(1).
		Return("123456")

	store.
		EXPECT().
		GetEmail(gomock.Any(), "mail@mail.com").
		Times(1).
		Return(0, nil)

	store.
		EXPECT().
		CreateUser(gomock.Any(), &inout.CreateUserRequestV1{
			Password:  "olleh",
			Email:     "mail@mail.com",
			EmailCode: "123456",
		}).
		Times(1).
		Return(enums.Ok, &models.User{
			Id:      userID,
			Created: 2,
		})

	esb.
		EXPECT().
		OnUserChanged([]uuid.UUID{userID}).
		Times(1)

	r.Header.Add("Content-Type", "application/json")
	status, message := createUserV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusCreated)
	v, ok := message.(*inout.GetUserResponseV1)
	require.True(t, ok)
	require.NotNil(t, v)
}

func TestGetUsersV1QueryForAdminUser(t *testing.T) {
	t.Parallel()

	adminUserID := uuid.NewV4()
	requestedIdentifiers := []string{uuid.NewV4().String(), adminUserID.String()}
	require.Equal(t, repositories.GetUsersQuery{
		Limit: 50,
		Page:  1,
		Id:    functools.StringsSliceToUUIDSlice(requestedIdentifiers),
	}, GetUsersV1Query(map[string][]string{
		"id": requestedIdentifiers,
	}, &models.AccessTokenPayload{
		IsAdmin: true,
		UserID:  adminUserID,
	}))
}

func TestGetUsersV1QueryForRegularUser(t *testing.T) {
	t.Parallel()

	userID := uuid.NewV4()
	requestedIdentifiers := []string{uuid.NewV4().String(), userID.String()}
	require.Equal(t, repositories.GetUsersQuery{
		Limit: 50,
		Page:  1,
		Id:    []uuid.UUID{userID},
	}, GetUsersV1Query(map[string][]string{
		"id": requestedIdentifiers,
	}, &models.AccessTokenPayload{
		IsAdmin: false,
		UserID:  userID,
	}))
}
