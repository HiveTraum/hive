package api

import (
	"hive/enums"
	"hive/inout"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (api *API) CreatePhoneV1(w http.ResponseWriter, r *http.Request) {
	body := &inout.CreatePhoneResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, phone := api.Controller.CreatePhone(r.Context(), body.Phone, body.Code, uuid.FromBytesOrNil(body.UserID))

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_Ok{
				Ok: &inout.Phone{
					Id:      phone.Id.Bytes(),
					Created: phone.Created,
					UserID:  phone.UserId.Bytes(),
					Phone:   phone.Value,
				}}})
	case enums.IncorrectPhoneCode:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					Code: []string{"Неверный код"},
				}}})
	case enums.PhoneNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					Phone: []string{"Не удалось найти код для данного телефона."},
				}}})
	case enums.UserNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					UserID: []string{"Такого пользователя не существует"},
				}}})
	case enums.IncorrectPhone:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePhoneResponseV1{
			Data: &inout.CreatePhoneResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneResponseV1_ValidationError{
					Phone: []string{"Некорректный номер телефона"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}

func (api *API) CreatePhoneConfirmationV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreatePhoneConfirmationResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, phoneConfirmation := api.Controller.CreatePhoneConfirmation(r.Context(), body.Phone)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreatePhoneConfirmationResponseV1{
			Data: &inout.CreatePhoneConfirmationResponseV1_Ok{
				Ok: &inout.PhoneConfirmation{
					Created: phoneConfirmation.Created,
					Expire:  phoneConfirmation.Expire,
					Phone:   phoneConfirmation.Phone,
				}}})
	case enums.IncorrectPhone:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreatePhoneConfirmationResponseV1{
			Data: &inout.CreatePhoneConfirmationResponseV1_ValidationError_{
				ValidationError: &inout.CreatePhoneConfirmationResponseV1_ValidationError{
					Phone: []string{"Некорректный номер телефона"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}
