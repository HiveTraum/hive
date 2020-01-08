package api

import (
	"auth/app"
	"auth/controllers"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"github.com/golang/protobuf/proto"
	"net/http"
)

func createSessionV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := inout.CreateSessionRequestV1{}
	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	loginController := app.GetLoginController()
	ctx := r.Context()

	_, session, token := controllers.CreateSession(app.GetStore(), loginController, ctx, body)

	return http.StatusCreated, &inout.CreateSessionResponseV1{
		RefreshToken: session.RefreshToken,
		AccessToken:  token,
	}
}

func SessionsAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return createSessionV1(r, app)
	}
}
