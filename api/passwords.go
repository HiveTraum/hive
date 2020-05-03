package api

import (
	"auth/enums"
	"auth/inout"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (api *API) CreatePasswordV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreatePasswordResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, password := api.Controller.CreatePassword(r.Context(), uuid.FromBytesOrNil(body.UserID), body.Value)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreatePasswordResponseV1{
			Data: &inout.CreatePasswordResponseV1_Ok{
				Ok: &inout.Password{
					Id:      password.Id.Bytes(),
					Created: password.Created,
					UserID:  password.UserId.Bytes(),
				}}})
	case enums.UserNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePasswordResponseV1{
			Data: &inout.CreatePasswordResponseV1_ValidationError_{
				ValidationError: &inout.CreatePasswordResponseV1_ValidationError{
					UserID: []string{"Пользователь не найден"},
				}}})
	case enums.IncorrectPassword:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePasswordResponseV1{
			Data: &inout.CreatePasswordResponseV1_ValidationError_{
				ValidationError: &inout.CreatePasswordResponseV1_ValidationError{
					Value: []string{"Не удалось обработать полученный пароль, попробуйте другой"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}
