package controllers

import (
	"auth/enums"
	"auth/inout"
	"auth/mocks"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	_, store, esb := mocks.InitMockApp(ctrl)
	ctx := context.Background()

	store.
		EXPECT().
		GetPhoneConfirmationCode(ctx, "+71234567890").
		Return("123456")

	store.
		EXPECT().
		GetPhone(ctx, "+71234567890").
		Return(enums.Ok, &models.Phone{
			Id:      1,
			Created: 0,
			UserId:  2,
			Value:   "+71234567890",
		})

	store.
		EXPECT().
		GetEmailConfirmationCode("email@email.com").
		Return("654321")

	store.
		EXPECT().
		GetEmail(ctx, "email@email.com").
		Return(enums.Ok, &models.Email{
			Id:      3,
			Created: 0,
			UserId:  4,
			Value:   "email@email.com",
		})

	body := &inout.CreateUserRequestV1{
		Password:  "123",
		Phone:     "+71234567890",
		Email:     "email@email.com",
		PhoneCode: "123456",
		EmailCode: "654321",
	}

	store.
		EXPECT().
		CreateUser(ctx, body).
		Return(enums.Ok, &models.User{
			Id:      5,
			Created: 0,
		})

	esb.
		EXPECT().
		OnUserChanged([]int64{2, 4, 5}).
		Times(1)

	body.Phone = "71234567890"

	status, user := CreateUser(store, esb, ctx, body)

	require.Equal(t, &models.User{
		Id:      5,
		Created: 0,
	}, user)

	require.Equal(t, enums.Ok, status)

}
