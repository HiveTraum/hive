package api

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"hive/auth/backends"
	"hive/enums"
	"hive/functools"
	"hive/inout"
	"hive/models"
	"hive/repositories"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserParsingJSON(t *testing.T) {
	t.Parallel()
	body := "{\"password\": \"password\", \"email\": \"email\", \"emailCode\": \"emailCode\", \"phone\": \"phone\", \"phoneCode\": \"phoneCode\"}"
	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	request.Header.Add("Content-Type", "application/json")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreateUser(ctx, "password", "email", "emailCode", "phone", "phoneCode").
		Return(enums.Ok, &models.User{})

	api.API.CreateUserV1(httptest.NewRecorder(), request)

	ctrl.Finish()
}

func TestCreateUserParsingJSONPB(t *testing.T) {
	t.Parallel()
	body, _ := protojson.Marshal(&inout.CreateUserResponseV1_Request{
		Password:  "password",
		Phone:     "phone",
		Email:     "email",
		PhoneCode: "phoneCode",
		EmailCode: "emailCode",
	})
	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreateUser(ctx, "password", "email", "emailCode", "phone", "phoneCode").
		Return(enums.Ok, &models.User{})

	api.API.CreateUserV1(httptest.NewRecorder(), request)

	ctrl.Finish()
}

func TestCreateUserParsingBinary(t *testing.T) {
	t.Parallel()
	body, _ := proto.Marshal(&inout.CreateUserResponseV1_Request{
		Password:  "password",
		Phone:     "phone",
		Email:     "email",
		PhoneCode: "phoneCode",
		EmailCode: "emailCode",
	})
	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/octet-stream")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreateUser(ctx, "password", "email", "emailCode", "phone", "phoneCode").
		Return(enums.Ok, &models.User{})

	api.API.CreateUserV1(httptest.NewRecorder(), request)

	ctrl.Finish()
}

func TestCreateUserRenderingJSON(t *testing.T) {
	t.Parallel()
	body := "{}"
	userID := uuid.NewV4()

	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	request.Header.Add("Content-Type", "application/json")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreateUser(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(enums.Ok, &models.User{Id: userID, Created: 1})

	recorder := httptest.NewRecorder()
	api.API.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message inout.CreateUserResponseV1
	_ = protojson.Unmarshal(responseBody, &message)
	user := message.GetOk()
	require.NotNil(t, user)
	require.Equal(t, userID.Bytes(), user.Id)
	require.Equal(t, int64(1), user.Created)
	ctrl.Finish()
}

func TestCreateUserRenderingBinary(t *testing.T) {
	t.Parallel()
	body, _ := proto.Marshal(&inout.CreateUserResponseV1_Request{})
	userID := uuid.NewV4()

	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/octet-stream")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreateUser(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(enums.Ok, &models.User{Id: userID, Created: 1})

	recorder := httptest.NewRecorder()
	api.API.CreateUserV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message inout.CreateUserResponseV1
	_ = proto.Unmarshal(responseBody, &message)
	user := message.GetOk()
	require.NotNil(t, user)
	require.Equal(t, userID.Bytes(), user.Id)
	require.Equal(t, int64(1), user.Created)
	ctrl.Finish()
}

func TestGetUsersV1QueryForAdminUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)
	adminUserID := uuid.NewV4()
	requestedIdentifiers := []string{uuid.NewV4().String(), adminUserID.String()}
	require.Equal(t, repositories.GetUsersQuery{
		Limit: 50,
		Page:  1,
		Id:    functools.StringsSliceToUUIDSlice(requestedIdentifiers),
	}, api.API.GetUsersV1Query(map[string][]string{
		"id": requestedIdentifiers,
	}, &backends.BasicAuthenticationBackendUser{
		IsAdmin: true,
		UserID:  adminUserID,
	}))
	ctrl.Finish()
}

func TestGetUsersV1QueryForRegularUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)
	userID := uuid.NewV4()
	requestedIdentifiers := []string{uuid.NewV4().String(), userID.String()}
	require.Equal(t, repositories.GetUsersQuery{
		Limit: 50,
		Page:  1,
		Id:    []uuid.UUID{userID},
	}, api.API.GetUsersV1Query(map[string][]string{
		"id": requestedIdentifiers,
	}, &backends.BasicAuthenticationBackendUser{
		IsAdmin: false,
		UserID:  userID,
	}))
	ctrl.Finish()
}
