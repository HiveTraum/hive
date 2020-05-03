package api

import (
	"auth/enums"
	"auth/inout"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (api *API) CreateEmailV1(w http.ResponseWriter, r *http.Request) {
	body := &inout.CreateEmailResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, email := api.Controller.CreateEmail(r.Context(), body.Email, body.Code, uuid.FromBytesOrNil(body.UserID))

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_Ok{
				Ok: &inout.Email{
					Id:      email.Id.Bytes(),
					Created: email.Created,
					UserID:  email.UserId.Bytes(),
					Email:   email.Value,
				}}})
	case enums.IncorrectEmailCode:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					Code: []string{"Неверный код"},
				}}})
	case enums.EmailNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					Email: []string{"Не удалось найти код для данного email."},
				}}})
	case enums.UserNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					UserId: []string{"Такого пользователя не существует"},
				}}})
	case enums.IncorrectEmail:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateEmailResponseV1{
			Data: &inout.CreateEmailResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailResponseV1_ValidationError{
					Email: []string{"Некорректный email"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}

func (api *API) CreateEmailConfirmationV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreateEmailConfirmationResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, emailConfirmation := api.Controller.CreateEmailConfirmation(r.Context(), body.Email)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreateEmailConfirmationResponseV1{
			Data: &inout.CreateEmailConfirmationResponseV1_Ok{
				Ok: &inout.EmailConfirmation{
					Created: emailConfirmation.Created,
					Expire:  emailConfirmation.Expire,
					Email:   emailConfirmation.Email,
				}}})
	case enums.IncorrectPhone:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateEmailConfirmationResponseV1{
			Data: &inout.CreateEmailConfirmationResponseV1_ValidationError_{
				ValidationError: &inout.CreateEmailConfirmationResponseV1_ValidationError{
					Email: []string{"Некорректный email"}}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}
