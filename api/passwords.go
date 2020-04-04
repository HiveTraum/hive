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

func createPasswordV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.CreatePasswordResponseV1) {

	body := inout.CreatePasswordResponseV1_Request{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, password := controllers.CreatePassword(app.GetStore(), app.GetESB(), app.GetPasswordProcessor(), r.Context(), uuid.FromBytesOrNil(body.UserID), body.Value)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreatePasswordResponseV1{
			Data: &inout.CreatePasswordResponseV1_Ok{
				Ok: &inout.Password{
					Id:      password.Id.Bytes(),
					Created: password.Created,
					UserID:  password.UserId.Bytes(),
				}}}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreatePasswordResponseV1{
			Data: &inout.CreatePasswordResponseV1_ValidationError_{
				ValidationError: &inout.CreatePasswordResponseV1_ValidationError{
					UserID: []string{"Пользователь не найден"},
				}}}
	case enums.IncorrectPassword:
		return http.StatusBadRequest, &inout.CreatePasswordResponseV1{
			Data: &inout.CreatePasswordResponseV1_ValidationError_{
				ValidationError: &inout.CreatePasswordResponseV1_ValidationError{
					Value: []string{"Не удалось обработать полученный пароль, попробуйте другой"},
				}}}
	default:
		return unhandledStatus(r, status), nil
	}
}

func PasswordsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createPasswordV1(r, app)
	}
}
