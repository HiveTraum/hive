package api

import (
	"auth/enums"
	"auth/inout"
	"auth/models"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePasswordParsingJSON(t *testing.T) {
	t.Parallel()
	userID := uuid.NewV4()
	encodedUserID := base64.StdEncoding.EncodeToString(userID.Bytes())
	body := fmt.Sprintf("{\"userID\": \"%s\", \"value\": \"password\"}", encodedUserID)

	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	request.Header.Add("Content-Type", "application/json")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreatePassword(ctx, userID, "password").
		Return(enums.Ok, &models.Password{})

	api.API.CreatePasswordV1(httptest.NewRecorder(), request)

	ctrl.Finish()
}

func TestCreatePasswordParsingJSONPB(t *testing.T) {
	t.Parallel()
	userID := uuid.NewV4()
	body, _ := protojson.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value:  "password",
	})
	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreatePassword(ctx, userID, "password").
		Return(enums.Ok, &models.Password{})

	api.API.CreatePasswordV1(httptest.NewRecorder(), request)

	ctrl.Finish()
}

func TestCreatePasswordParsingBinary(t *testing.T) {
	t.Parallel()
	userID := uuid.NewV4()
	body, _ := proto.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value:  "password",
	})
	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/octet-stream")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreatePassword(ctx, userID, "password").
		Return(enums.Ok, &models.Password{})

	api.API.CreatePasswordV1(httptest.NewRecorder(), request)

	ctrl.Finish()
}

func TestCreatePasswordRenderingJSON(t *testing.T) {
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
		CreatePassword(ctx, gomock.Any(), gomock.Any()).
		Return(enums.Ok, &models.Password{Id: userID, Created: 1})

	recorder := httptest.NewRecorder()
	api.API.CreatePasswordV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message inout.CreatePasswordResponseV1
	_ = protojson.Unmarshal(responseBody, &message)
	user := message.GetOk()
	require.NotNil(t, user)
	require.Equal(t, userID.Bytes(), user.Id)
	require.Equal(t, int64(1), user.Created)
	ctrl.Finish()
}

func TestCreatePasswordRenderingBinary(t *testing.T) {
	t.Parallel()
	body, _ := proto.Marshal(&inout.CreatePasswordResponseV1_Request{})
	userID := uuid.NewV4()

	ctrl := gomock.NewController(t)
	api := InitAPIWithMockedInternals(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/octet-stream")
	ctx := request.Context()

	api.
		Controller.
		EXPECT().
		CreatePassword(ctx, gomock.Any(), gomock.Any()).
		Return(enums.Ok, &models.Password{Id: userID, Created: 1})

	recorder := httptest.NewRecorder()
	api.API.CreatePasswordV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message inout.CreatePasswordResponseV1
	_ = proto.Unmarshal(responseBody, &message)
	user := message.GetOk()
	require.Equal(t, http.StatusCreated, response.StatusCode)
	require.NotNil(t, user)
	require.Equal(t, userID.Bytes(), user.Id)
	require.Equal(t, int64(1), user.Created)
	ctrl.Finish()
}
