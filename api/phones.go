package api

import (
	"auth/app"
	"auth/controllers"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func createPhoneV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.CreatePhoneResponseV1) {
	b := inout.CreatePhoneResponseV1_Request{}
	err := r.ParseBody(&b)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, phone := controllers.CreatePhone(app.GetStore(), app.GetESB(), r.Context(), b.Phone, b.Code, uuid.FromBytesOrNil(b.UserID), b.PhoneCountryCode)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_Ok{
				Ok: &inout.Phone{
					Id:               phone.Id.Bytes(),
					Created:          phone.Created,
					UserID:           phone.UserId.Bytes(),
					Phone:            phone.Value,
					PhoneCountryCode: phone.CountryCode,
				}}}
	case enums.IncorrectPhoneCode:
		return http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					Code: []string{"Неверный код"},
				}}}
	case enums.PhoneNotFound:
		return http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					Phone: []string{"Не удалось найти код для данного телефона."},
				}}}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					UserID: []string{"Такого пользователя не существует"},
				}}}
	case enums.IncorrectPhone:
		return http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					Phone: []string{"Некорректный номер телефона"},
				}}}
	default:
		return unhandledStatus(r, status), nil
	}
}

func PhonesAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createPhoneV1(r, app)
	}
}

func createPhoneConfirmationV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.CreatePhoneConfirmationResponseV1) {

	body := inout.CreatePhoneConfirmationResponseV1_Request{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, phoneConfirmation := controllers.CreatePhoneConfirmation(app.GetStore(), app.GetESB(), r.Context(), body.Phone, body.PhoneCountryCode)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreatePhoneConfirmationResponseV1{
			Data: &inout.CreatePhoneConfirmationResponseV1_Ok{
				Ok: &inout.PhoneConfirmation{
					Created: phoneConfirmation.Created,
					Expire:  phoneConfirmation.Expire,
					Phone:   phoneConfirmation.Phone,
				}}}
	case enums.IncorrectPhone:
		return http.StatusBadRequest, &inout.CreatePhoneConfirmationResponseV1{
			Data: &inout.CreatePhoneConfirmationResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneConfirmationResponseV1_ValidationError{
					Phone: []string{"Некорректный номер телефона"},
				}}}
	default:
		return unhandledStatus(r, status), nil
	}
}

func PhoneConfirmationsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createPhoneConfirmationV1(r, app)
	}
}
