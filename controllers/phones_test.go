package controllers

import (
	"hive/enums"
	"hive/functools"
	"hive/inout"
	"hive/models"
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
	controller := InitControllerWithMockedInternals(ctrl)
	ctx := context.Background()

	phone := "+79234567890"
	formattedPhone := functools.NormalizePhone(phone)

	controller.
		Store.
		EXPECT().
		GetRandomCodeForPhoneConfirmation().
		Return("123456").
		Times(1)

	controller.
		Store.
		EXPECT().
		CreatePhoneConfirmationCode(ctx, formattedPhone, "123456", time.Minute*15).
		Return(&models.PhoneConfirmation{
			Created: 0,
			Expire:  0,
			Phone:   formattedPhone,
			Code:    "123456",
		})

	controller.
		Dispatcher.
		EXPECT().
		Send("phoneConfirmation", int32(1), &inout.CreatePhoneConfirmationEventV1{
			Phone: formattedPhone,
			Code:  "123456",
		})

	status, confirmation := controller.Controller.CreatePhoneConfirmation(ctx, phone)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, confirmation)
	require.Equal(t, formattedPhone, confirmation.Phone)
}

func TestCreatePhoneConfirmationWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	controller := InitControllerWithMockedInternals(ctrl)
	ctx := context.Background()
	phone := "qwerty"
	status, confirmation := controller.Controller.CreatePhoneConfirmation(ctx, phone)
	require.Nil(t, confirmation)
	require.Equal(t, enums.IncorrectPhone, status)
}
