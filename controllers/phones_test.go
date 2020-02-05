package controllers

import (
	"auth/enums"
	"auth/mocks"
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
	_, store, esb, _ := mocks.InitMockApp(ctrl)
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

	esb.
		EXPECT().
		OnPhoneCodeConfirmationCreated("+71234567890", "123456").
		Times(1)

	phone := "71234567890"
	status, confirmation := CreatePhoneConfirmation(store, esb, ctx, phone)
	require.Equal(t, "+71234567890", confirmation.Phone)
	require.NotNil(t, confirmation)
	require.Equal(t, enums.Ok, status)
}

func TestCreatePhoneConfirmationWithIncorrectPhone(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, store, esb, _ := mocks.InitMockApp(ctrl)
	ctx := context.Background()
	phone := "qwerty"
	status, confirmation := CreatePhoneConfirmation(store, esb, ctx, phone)
	require.Nil(t, confirmation)
	require.Equal(t, enums.IncorrectPhone, status)
}
