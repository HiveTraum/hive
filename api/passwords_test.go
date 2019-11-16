package api

import (
	"auth/enums"
	"auth/functools"
	"auth/inout"
	"auth/mocks"
	"auth/models"
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePasswordWithoutUserV1(t *testing.T) {
	t.Parallel()
	body := "{\"user_id\": 1, \"value\": \"hello\"}"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, _ := mocks.InitMockApp(ctrl)

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	store.
		EXPECT().
		CreatePassword(gomock.Any(), int64(1), gomock.Any()).
		Times(1).
		Return(enums.UserNotFound, nil)
	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusBadRequest)
	v, ok := message.(*inout.CreatePasswordBadRequestResponseV1)
	require.True(t, ok)
	require.Len(t, v.UserId, 1)
	require.Len(t, v.Value, 0)
}

func TestCreatePasswordWithUserV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb := mocks.InitMockApp(ctrl)
	ctx := context.Background()
	store.
		EXPECT().
		CreatePassword(ctx, int64(2), gomock.Any()).
		Return(enums.Ok, &models.Password{
			Id:      1,
			Created: 0,
			UserId:  2,
			Value:   "some value",
		}).
		Times(1)

	esb.
		EXPECT().
		OnPasswordChanged(int64(2)).
		Times(1)

	body := fmt.Sprintf("{\"user_id\": %d, \"value\": \"hello\"}", 2)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusCreated)
	v, ok := message.(*inout.CreatePasswordResponseV1)
	require.True(t, ok)
	require.Equal(t, v.UserId, int64(2))
}

func TestCreatePasswordWithTooLongValueV1(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb := mocks.InitMockApp(ctrl)
	ctx := context.Background()
	store.
		EXPECT().
		CreatePassword(ctx, int64(3), gomock.Any()).
		Return(enums.Ok, &models.Password{
			Id:      1,
			Created: 0,
			UserId:  3,
			Value:   "some value",
		}).
		Times(1)

	esb.
		EXPECT().
		OnPasswordChanged(int64(3)).
		Times(1)

	body := fmt.Sprintf("{\"user_id\": %d, \"value\": \"hellohellohellohellohellohellohellohellohellohellohell"+
		"ohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohe"+
		"llohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello"+
		"hellohellohellohellohellohello\"}", 3)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	r.Header.Add("Content-Type", "application/json")
	status, message := createPasswordV1(&functools.Request{Request: r}, app)
	require.Equal(t, status, http.StatusCreated)
	_, ok := message.(*inout.CreatePasswordResponseV1)
	require.True(t, ok)
}
