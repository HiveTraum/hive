package api

import (
	"auth/auth/backends"
	"auth/enums"
	"auth/functools"
	"auth/inout"
	"auth/repositories"
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserEmptyBody(t *testing.T) {
	t.Parallel()
	body := "{}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	ctx := request.Context()

	controller.
		EXPECT().
		CreateUser(ctx, &inout.CreateUserResponseV1_Request{}).
		Return(enums.PasswordRequired, nil)

	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.Password, 1)
}

func TestCreateUserWithOnlyPassword(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	ctx := request.Context()

	controller.
		EXPECT().
		CreateUser(ctx, &inout.CreateUserResponseV1_Request{
			Password: "hello",
		}).
		Return(enums.MinimumOneFieldRequired, nil)

	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.Errors, 1)
}

func TestCreateUserWithOnlyEmail(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	ctx := request.Context()

	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	controller.
		EXPECT().
		CreateUser(ctx, &inout.CreateUserResponseV1_Request{
			Password: "hello",
			Email:    "mail@mail.com",
		}).
		Return(enums.EmailConfirmationCodeNotFound, nil)

	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.EmailCode, 1)
}

func TestCreateUserWithEmailAndEmailCode(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\", \"emailCode\": \"123456\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	ctx := request.Context()

	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	controller.
		EXPECT().
		CreateUser(ctx, &inout.CreateUserResponseV1_Request{
			Password:  "hello",
			Email:     "mail@mail.com",
			EmailCode: "123456",
		}).
		Return(enums.EmailConfirmationCodeNotFound, nil)

	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.Email, 1)
}

func TestCreateUserWithEmailAndEmailCodeAfterIncorrectEmailConfirmationCodeReceived(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\", \"emailCode\": \"123456\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	ctx := request.Context()

	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	controller.
		EXPECT().
		CreateUser(ctx, &inout.CreateUserResponseV1_Request{
			Password: "hello",
			Email:    "mail@mail.com",
		}).
		Return(enums.EmailNotFound, nil)

	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.EmailCode, 1)
}

func TestSuccessfulCreateUserWithEmail(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"hello\", \"email\": \"mail@mail.com\", \"emailCode\": \"123456\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	ctx := request.Context()

	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	controller.
		EXPECT().
		CreateUser(ctx, &inout.CreateUserResponseV1_Request{
			Password:  "hello",
			Email:     "mail@mail.com",
			EmailCode: "123456",
		}).
		Return(enums.Ok, nil)

	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusCreated)
	user := message.GetOk()
	require.NotNil(t, user)
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
	}, &backends.BasicAuthenticationBackendUser{
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
	}, &backends.BasicAuthenticationBackendUser{
		IsAdmin: false,
		UserID:  userID,
	}))
}
