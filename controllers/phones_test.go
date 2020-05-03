package controllers

import (
	"auth/enums"
	"auth/inout"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreatePhoneConfirmation(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	controller, store, dispatcher, _ := InitControllerWithMockedInternals(ctrl)
	ctx := context.Background()

	store.
		EXPECT().
		GetRandomCodeForPhoneConfirmation().
		Return("123456").
		Times(1)

	store.
		EXPECT().
		CreatePhoneConfirmationCode(ctx, "+71234567890", "123456", time.Minute*15).
		Return(&models.PhoneConfirmation{
			Created: 0,
			Expire:  0,
			Phone:   "+71234567890",
			Code:    "123456",
		})

	dispatcher.
		EXPECT().
		Send("phoneConfirmation", 1, &inout.CreatePhoneConfirmationEventV1{
			Phone: "+71234567890",
			Code:  "123456",
		})

	phone := "71234567890"
	status, confirmation := controller.CreatePhoneConfirmation(ctx, phone)
	require.Equal(t, "+71234567890", confirmation.Phone)
	require.NotNil(t, confirmation)
	require.Equal(t, enums.Ok, status)
}

func TestCreatePhoneConfirmationWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	controller, _, _, _ := InitControllerWithMockedInternals(ctrl)
	ctx := context.Background()
	phone := "qwerty"
	status, confirmation := controller.CreatePhoneConfirmation(ctx, phone)
	require.Nil(t, confirmation)
	require.Equal(t, enums.IncorrectPhone, status)
}
