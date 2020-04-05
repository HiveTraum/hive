package api

import (
	"auth/enums"
	"auth/functools"
	"auth/inout"
	"auth/mocks"
	"auth/models"
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
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

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hello").
		Return("olleh")

	store.
		EXPECT().
		CreatePassword(gomock.Any(), userID, "olleh").
		Times(1).
		Return(enums.UserNotFound, nil)

	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	validationError := message.GetValidationError()
	require.NotNil(t, validationError)
	require.Len(t, validationError.UserID, 1)
	require.Len(t, validationError.Value, 0)
}

func TestCreatePasswordWithUserV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb, _, passwordProcessor := mocks.InitMockApp(ctrl)
	ctx := context.Background()

	userID := uuid.NewV4()

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hello").
		Return("olleh")

	store.
		EXPECT().
		CreatePassword(ctx, userID, "olleh").
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 0,
			UserId:  userID,
			Value:   "some value",
		}).
		Times(1)

	esb.
		EXPECT().
		OnPasswordChanged(userID).
		Times(1)

	body, _ := json.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value:  "hello",
	})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusCreated)
	password := message.GetOk()
	require.NotNil(t, password)
	require.Equal(t, userID.Bytes(), password.UserID)
}

func TestCreatePasswordWithTooLongValueV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb, _, passwordProcessor := mocks.InitMockApp(ctrl)
	ctx := context.Background()

	userID := uuid.NewV4()

	passwordProcessor.
		EXPECT().
		EncodePassword(gomock.Any(), "hellohellohellohellohellohellohellohellohellohellohell"+
			"ohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohe"+
			"llohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello"+
			"hellohellohellohellohellohello").
		Return("olleh")

	store.
		EXPECT().
		CreatePassword(ctx, userID, "olleh").
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 0,
			UserId:  userID,
			Value:   "some value",
		}).
		Times(1)

	esb.
		EXPECT().
		OnPasswordChanged(userID).
		Times(1)

	body, _ := json.Marshal(&inout.CreatePasswordResponseV1_Request{
		UserID: userID.Bytes(),
		Value: "hellohellohellohellohellohellohellohellohellohellohell" +
			"ohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohe" +
			"llohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello" +
			"hellohellohellohellohellohello",
	})

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusCreated)
	password := message.GetOk()
	require.NotNil(t, password)
}
