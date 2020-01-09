package api

import (
	"auth/app"
	"auth/controllers"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"auth/models"
	"github.com/golang/protobuf/proto"
	"net/http"
)

func createEmailV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {
	body := inout.CreateEmailRequestV1{}
	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, email := controllers.CreateEmail(app.GetStore(), app.GetESB(), r.Context(), body.Email, body.Code, models.UserID(body.UserId))

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreateEmailResponseV1{
			Id:      int64(email.Id),
			Created: email.Created,
			UserId:  int64(email.UserId),
			Email:   email.Value,
		}
	case enums.IncorrectEmailCode:
		return http.StatusBadRequest, &inout.CreateEmailBadRequestV1{
			Code: []string{"Неверный код"},
		}
	case enums.EmailNotFound:
		return http.StatusBadRequest, &inout.CreateEmailBadRequestV1{
			Email: []string{"Не удалось найти код для данного email."},
		}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreateEmailBadRequestV1{
			UserId: []string{"Такого пользователя не существует"},
		}
	case enums.IncorrectEmail:
		return http.StatusBadRequest, &inout.CreateEmailBadRequestV1{
			Email: []string{"Некорректный email"},
		}
	default:
		return http.StatusCreated, &inout.CreateEmailResponseV1{
			Id:      int64(email.Id),
			Created: email.Created,
			UserId:  int64(email.UserId),
			Email:   email.Value,
		}
	}
}

func EmailsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createEmailV1(r, app)
	}
}

func createEmailConfirmationV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := inout.CreateEmailConfirmationRequestV1{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, emailConfirmation := controllers.CreateEmailConfirmation(r.Context(), app.GetStore(), app.GetESB(), body.Email)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreateEmailConfirmationResponseV1{
			Created: emailConfirmation.Created,
			Expire:  emailConfirmation.Expire,
			Email:   emailConfirmation.Email,
		}
	case enums.IncorrectPhone:
		return http.StatusBadRequest, &inout.CreateEmailConfirmationBadRequestV1{
			Email: []string{"Некорректный email"},
		}
	default:
		return http.StatusCreated, &inout.CreateEmailConfirmationResponseV1{
			Created: emailConfirmation.Created,
			Expire:  emailConfirmation.Expire,
			Email:   emailConfirmation.Email,
		}
	}
}

func EmailConfirmationsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createEmailConfirmationV1(r, app)
	}
}
