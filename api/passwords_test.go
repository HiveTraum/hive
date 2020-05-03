package api

import (
	"auth/enums"
	"auth/inout"
	"auth/models"
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePasswordWithoutUserV1(t *testing.T) {
	t.Parallel()
	userID := uuid.NewV4()

	body, _ := json.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value:  "hello",
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, _, _, passwordProcessor := mocks.InitMockApp(ctrl)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	ctx := request.Context()

	controller.
		EXPECT().
		CreatePassword(ctx, userID, "hello").
		Return(enums.UserNotFound, nil).
		Times(1)

	request.Header.Add("Content-Type", string(enums.JSONContentType))
	recorder := httptest.NewRecorder()
	api.CreatePasswordV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message inout.CreatePasswordResponseV1
	_ = json.Unmarshal(responseBody, &message)
	require.Equal(t, response.StatusCode, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.UserID, 1)
	require.Len(t, validationError.Value, 0)
}

func TestCreatePasswordWithUserV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api, controller, _ := InitAPIWithMockedInternals(ctrl)

	userID := uuid.NewV4()

	body, _ := json.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value:  "hello",
	})

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	ctx := request.Context()

	controller.
		EXPECT().
		CreatePassword(ctx, userID, "hello").
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value:   "hello",
		})

	recorder := httptest.NewRecorder()
	api.CreatePasswordV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreatePasswordResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.Status, http.StatusCreated)
	password := message.GetOk()
	require.NotNil(t, password)
	require.Equal(t, userID.Bytes(), password.UserID)
}

func TestCreatePasswordWithTooLongValueV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api, controller, _ := InitAPIWithMockedInternals(ctrl)
	ctx := context.Background()

	userID := uuid.NewV4()

	controller.
		EXPECT().
		CreatePassword(ctx, userID, "hellohellohellohellohellohellohellohellohellohellohell"+
			"ohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohe"+
			"llohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello"+
			"hellohellohellohellohellohello").
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 1,
			UserId:  userID,
			Value: "hellohellohellohellohellohellohellohellohellohellohell" +
				"ohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohe" +
				"llohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello" +
				"hellohellohellohellohellohello",
		})

	body, _ := json.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value: "hellohellohellohellohellohellohellohellohellohellohell" +
			"ohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohe" +
			"llohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello" +
			"hellohellohellohellohellohello",
	})

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.CreatePasswordV1(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var message *inout.CreatePasswordResponseV1
	_ = proto.Unmarshal(responseBody, message)
	require.Equal(t, response.StatusCode, http.StatusCreated)
	require.NotNil(t, message)
	password := message.GetOk()
	require.NotNil(t, password)
}
