package controllers

import (
	"auth/enums"
	"auth/mocks"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app, store, esb, passwordProcessor := mocks.InitMockApp(ctrl)
	ctx := context.Background()

	passwordProcessor.
		EXPECT().
		EncodePassword(ctx, "hello").
		Return("olleh")

	store.
		EXPECT().
		CreatePassword(ctx, models.UserID(10), gomock.Not("hello")).
		Return(enums.Ok, &models.Password{
			Id:      1,
			Created: 0,
			UserId:  10,
			Value:   "",
		})

	esb.
		EXPECT().
		OnPasswordChanged(models.UserID(10)).
		Times(1)

	status, password := CreatePassword(store, esb, app.GetLoginController(), ctx, 10, "hello")
	require.NotEqual(t, "hello", password.Value)
	require.NotNil(t, password)
	require.Equal(t, enums.Ok, status)
}
