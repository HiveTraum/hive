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

func createEmailV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.CreateEmailResponseV1) {
	body := inout.CreateEmailResponseV1_Request{}
	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, email := controllers.CreateEmail(app.GetStore(), app.GetESB(), r.Context(), body.Email, body.Code, uuid.FromBytesOrNil(body.UserID))

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_Ok{
				Ok: &inout.Email{
					Id:      email.Id.Bytes(),
					Created: email.Created,
					UserID:  email.UserId.Bytes(),
					Email:   email.Value,
				}}}
	case enums.IncorrectEmailCode:
		return http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					Code: []string{"Неверный код"},
				}}}
	case enums.EmailNotFound:
		return http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					Email: []string{"Не удалось найти код для данного email."},
				}}}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					UserId: []string{"Такого пользователя не существует"},
				}}}
	case enums.IncorrectEmail:
		return http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					Email: []string{"Некорректный email"},
				}}}
	default:
		return unhandledStatus(r, status), nil
	}
}

func EmailsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createEmailV1(r, app)
	}
}

func createEmailConfirmationV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.CreateEmailConfirmationResponseV1) {

	body := inout.CreateEmailConfirmationResponseV1_Request{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, emailConfirmation := controllers.CreateEmailConfirmation(r.Context(), app.GetStore(), app.GetESB(), body.Email)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreateEmailConfirmationResponseV1{
			Data: &inout.CreateEmailConfirmationResponseV1_Ok{
				Ok: &inout.EmailConfirmation{
					Created: emailConfirmation.Created,
					Expire:  emailConfirmation.Expire,
					Email:   emailConfirmation.Email,
				}}}
	case enums.IncorrectPhone:
		return http.StatusBadRequest, &inout.CreateEmailConfirmationResponseV1{
			Data: &inout.CreateEmailConfirmationResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailConfirmationResponseV1_ValidationError{
					Email: []string{"Некорректный email"}}}}
	default:
		return unhandledStatus(r, status), nil
	}
}

func EmailConfirmationsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createEmailConfirmationV1(r, app)
	}
}
