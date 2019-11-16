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
	"net/http"
)

func createPasswordV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := inout.CreatePasswordRequestV1{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, password := controllers.CreatePassword(app.GetStore(), app.GetESB(), r.Context(), body.UserId, body.Value)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreatePasswordResponseV1{
			Id:      password.Id,
			Created: password.Created,
			UserId:  password.UserId,
		}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreatePasswordBadRequestResponseV1{
			UserId: []string{"Пользователь не найден"},
		}
	case enums.IncorrectPassword:
		return http.StatusBadRequest, &inout.CreatePasswordBadRequestResponseV1{
			Value: []string{"Не удалось обработать полученный пароль, попробуйте другой"},
		}
	default:
		return http.StatusCreated, &inout.CreatePasswordResponseV1{
			Id:      password.Id,
			Created: password.Created,
			UserId:  password.UserId,
		}
	}
}

func PasswordsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createPasswordV1(r, app)
	}
}
