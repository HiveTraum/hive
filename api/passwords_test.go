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

	body, _ := json.Marshal(&inout.CreatePasswordRequestV1{
		UserID: userID.Bytes(),
		Value:  "hello",
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, _, loginController := mocks.InitMockApp(ctrl)

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	loginController.
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
	v, ok := message.(*inout.CreatePasswordBadRequestResponseV1)
	require.True(t, ok)
	require.Len(t, v.UserID, 1)
	require.Len(t, v.Value, 0)
}

func TestCreatePasswordWithUserV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb, passwordProcessor := mocks.InitMockApp(ctrl)
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

	body, _ := json.Marshal(&inout.CreatePasswordRequestV1{
		UserID: userID.Bytes(),
		Value:  "hello",
	})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusCreated)
	v, ok := message.(*inout.CreatePasswordResponseV1)
	require.True(t, ok)
	require.Equal(t, userID.Bytes(), v.UserID)
}

func TestCreatePasswordWithTooLongValueV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb, passwordProcessor := mocks.InitMockApp(ctrl)
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

	body, _ := json.Marshal(&inout.CreatePasswordRequestV1{
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
	_, ok := message.(*inout.CreatePasswordResponseV1)
	require.True(t, ok)
}
