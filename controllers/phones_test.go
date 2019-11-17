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
		CreatePhoneConfirmationCode(ctx, "+71234567890", gomock.Any(), time.Minute*15).
		Return(&models.PhoneConfirmation{
			Created: 0,
			Expire:  0,
			Phone:   "+71234567890",
			Code:    "123456",
		})

	esb.
		EXPECT().
		OnPhoneCodeConfirmationCreated("+71234567890", gomock.Any()).
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

//
//func TestCreatePhoneWithoutUser(t *testing.T) {
//	app := config.InitApp()
//	config.PurgeUsers(app.Db, context.Background())
//	config.PurgePhones(app.Db, context.Background())
//	phone := "71234567890"
//	status, confirmation := CreatePhoneConfirmation(app, phone)
//	require.Equal(t, "+71234567890", confirmation.Phone)
//	require.NotNil(t, confirmation)
//	require.Equal(t, enums.Ok, status)
//	status, phoneObj := CreatePhone(app, context.Background(), phone, confirmation.Code, 1)
//	require.Equal(t, enums.UserNotFound, status)
//	require.Nil(t, phoneObj)
//}
//
//func TestCreatePhoneWithIncorrectPhone(t *testing.T) {
//	app := config.InitApp()
//	config.PurgeUsers(app.Db, context.Background())
//	config.PurgePhones(app.Db, context.Background())
//	phone := "qwerty"
//	status, confirmation := CreatePhoneConfirmation(app, phone)
//	require.Equal(t, enums.IncorrectPhone, status)
//	require.Nil(t, confirmation)
//}
//
//func TestCreatePhoneWithUser(t *testing.T) {
//	app := config.InitApp()
//	ctx := context.Background()
//	config.PurgeUsers(app.Db, ctx)
//	config.PurgeUserViews(app.Db, ctx)
//	config.PurgePhones(app.Db, ctx)
//	user := CreateUser(app, ctx)
//	require.NotNil(t, user)
//	phone := "71234567890"
//	status, confirmation := CreatePhoneConfirmation(app, phone)
//	require.Equal(t, "+71234567890", confirmation.Phone)
//	require.NotNil(t, confirmation)
//	require.Equal(t, enums.Ok, status)
//	status, phoneObj := CreatePhone(app, ctx, phone, confirmation.Code, user.Id)
//	require.Equal(t, enums.Ok, status)
//	require.NotNil(t, phoneObj)
//	require.Equal(t, "+71234567890", phoneObj.Value)
//}
//
//func TestCreatePhoneWithAlreadyAssignedPhone(t *testing.T) {
//	app := config.InitApp()
//	ctx := context.Background()
//	config.PurgeUsers(app.Db, ctx)
//	config.PurgePhones(app.Db, ctx)
//	config.PurgeUserViews(app.Db, ctx)
//	user := CreateUser(app, ctx)
//	require.NotNil(t, user)
//	phone := "71234567890"
//	status, confirmation := CreatePhoneConfirmation(app, phone)
//	require.Equal(t, "+71234567890", confirmation.Phone)
//	require.NotNil(t, confirmation)
//	require.Equal(t, enums.Ok, status)
//	status, phoneObj := CreatePhone(app, ctx, phone, confirmation.Code, user.Id)
//	require.Equal(t, enums.Ok, status)
//	require.NotNil(t, phoneObj)
//	require.Equal(t, "+71234567890", phoneObj.Value)
//	_, phoneFromDB := repositories.GetPhone(app.Db, ctx, "+71234567890")
//	require.Equal(t, phoneFromDB.Id, phoneObj.Id)
//	require.Equal(t, phoneFromDB.UserId, user.Id)
//	newUser := CreateUser(app, ctx)
//	status, phoneObj = CreatePhone(app, ctx, phone, confirmation.Code, newUser.Id)
//	require.Equal(t, enums.Ok, status)
//	require.NotNil(t, phoneObj)
//	require.Equal(t, "+71234567890", phoneObj.Value)
//	_, phoneFromDB = repositories.GetPhone(app.Db, ctx, "+71234567890")
//	require.NotEqual(t, phoneFromDB.UserId, user.Id)
//}
