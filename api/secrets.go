package api

import (
	"auth/controllers"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func getSecretV1(r *functools.Request, app infrastructure.AppInterface, id uuid.UUID) (int, proto.Message) {

	secret := controllers.GetSecret(app.GetStore(), r.Context(), id)

	if secret == nil {
		return http.StatusNotFound, nil
	} else {
		return http.StatusOK, &inout.GetSecretResponseV1{
			Data: &inout.Secret{
				Id:      secret.Id.Bytes(),
				Created: secret.Created,
				Value:   secret.Value.Bytes(),
			}}
	}
}

func SecretAPIV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		vars := mux.Vars(r.Request)
		id, err := uuid.FromString(vars["id"])

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		return getSecretV1(r, app, id)
	}
}
