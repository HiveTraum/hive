package controllers

import (
	"auth/enums"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	controller, store, _, passwordProcessor := InitControllerWithMockedInternals(ctrl)
	ctx := context.Background()

	userID := uuid.NewV4()

	passwordProcessor.
		EXPECT().
		EncodePassword(ctx, "hello").
		Return("olleh")

	store.
		EXPECT().
		CreatePassword(ctx, userID, gomock.Not("hello")).
		Return(enums.Ok, &models.Password{
			Id:      uuid.NewV4(),
			Created: 0,
			UserId:  userID,
			Value:   "",
		})

	status, password := controller.CreatePassword(ctx, userID, "hello")
	require.NotEqual(t, "hello", password.Value)
	require.NotNil(t, password)
	require.Equal(t, enums.Ok, status)
}
