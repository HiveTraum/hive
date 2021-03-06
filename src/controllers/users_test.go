package controllers

//func TestCreateUser(t *testing.T) {
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	controller, esb, store, _, passwordProcessor := InitControllerWithMockedInternals(ctrl)
//	ctx := context.Background()
//
//	firstUserID, secondUserID, thirdUserID := uuid.NewV4(), uuid.NewV4(), uuid.NewV4()
//
//	passwordProcessor.
//		EXPECT().
//		EncodePassword(ctx, "hello").
//		Return("olleh")
//
//	store.
//		EXPECT().
//		GetPhoneConfirmationCode(ctx, "+71234567890").
//		Return("123456")
//
//	store.
//		EXPECT().
//		GetPhone(ctx, "+71234567890").
//		Return(enums.Ok, &models.Phone{
//			Id:      uuid.NewV4(),
//			Created: 0,
//			UserId:  firstUserID,
//			Value:   "+71234567890",
//		})
//
//	store.
//		EXPECT().
//		GetEmailConfirmationCode(ctx, "email@email.com").
//		Return("654321")
//
//	store.
//		EXPECT().
//		GetEmail(ctx, "email@email.com").
//		Return(enums.Ok, &models.Email{
//			Id:      uuid.NewV4(),
//			Created: 0,
//			UserId:  secondUserID,
//			Value:   "email@email.com",
//		})
//
//	body := inout.CreateUserResponseV1_Request{
//		Password:         "hello",
//		Phone:            "71234567890",
//		Email:            "email@email.com",
//		PhoneCode:        "123456",
//		EmailCode:        "654321",
//		PhoneCountryCode: "RU",
//	}
//
//	store.
//		EXPECT().
//		CreateUser(ctx, &inout.CreateUserResponseV1_Request{
//			Password:         "olleh",
//			Phone:            "+71234567890",
//			Email:            body.Email,
//			PhoneCode:        body.PhoneCode,
//			EmailCode:        body.EmailCode,
//			PhoneCountryCode: "RU",
//		}).
//		Return(enums.Ok, &models.User{
//			Id:      thirdUserID,
//			Created: 0,
//		})
//
//	esb.
//		EXPECT().
//		OnUserChanged([]uuid.UUID{firstUserID, secondUserID, thirdUserID}).
//		Times(1)
//
//	body.Phone = "71234567890"
//
//	status, user := controller.CreateUser(ctx, &body)
//
//	require.Equal(t, &models.User{
//		Id:      thirdUserID,
//		Created: 0,
//	}, user)
//
//	require.Equal(t, enums.Ok, status)
//}
