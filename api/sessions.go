package api

import (
	"hive/enums"
	"hive/inout"
	"hive/models"
	"hive/repositories"
	"net/http"
	"time"
)

func (api *API) CreateSessionV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreateSessionResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	ctx := r.Context()
	user := repositories.GetUserFromContext(ctx)
	refreshToken := repositories.GetRefreshTokenCookie(r, api.environment)

	var status int
	var session *models.Session

	if user != nil {
		session = api.Controller.CreateSession(ctx, user.GetUserID(), body.Fingerprint, body.UserAgent)
		status = enums.Ok
	} else if refreshToken != nil {
		status, session = api.Controller.UpdateSession(ctx, *refreshToken, body.Fingerprint, body.UserAgent)
	} else {
		api.Renderer.Render(w, r, http.StatusUnauthorized, nil)
		return
	}

	switch status {
	case enums.Ok:

		http.SetCookie(w, &http.Cookie{
			Name:     enums.RefreshToken,
			Value:    session.RefreshToken.String(),
			Domain:   r.Referer(),
			Expires:  time.Now().Add(time.Hour * 24 * time.Duration(api.environment.RefreshTokenLifetime)),
			Secure:   true,
			HttpOnly: true,
			Path:     "/api/v1/sessions",
		})

		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_Ok{
				Ok: &inout.Session{
					RefreshToken: session.RefreshToken.Bytes(),
					AccessToken:  session.AccessToken,
					Created:      session.Created,
					Expired:      session.Expires,
				}}})
	case
		enums.SessionNotFound,
		enums.UserNotFound,
		enums.IncorrectToken,
		enums.InvalidToken,
		enums.SecretNotFound:
		api.Renderer.Render(w, r, http.StatusUnauthorized, nil)
	case enums.IncorrectPassword:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					Password: []string{"Некорректный пароль"},
				}}})
	case enums.PasswordNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					Password: []string{"Не установлен пароль"},
				}}})
	case enums.EmailNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					Email: []string{"Email не найден"},
				}}})
	case enums.IncorrectEmail:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					Email: []string{"Некорректный email"},
				}}})
	case enums.IncorrectEmailCode:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					EmailCode: []string{"Некорректный код подтверждения"},
				}}})
	case enums.EmailConfirmationCodeNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					EmailCode: []string{"Не найден код подтверждения для данного email"},
				}}})
	case enums.PhoneNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					Phone: []string{"Телефон не найден"},
				}}})
	case enums.IncorrectPhone:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					Phone: []string{"Некорректный телефон"},
				}}})
	case enums.IncorrectPhoneCode:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					PhoneCode: []string{"Некорректный код подтверждения"},
				}}})
	case enums.PhoneConfirmationCodeNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateSessionResponseV1{
			Data: &inout.CreateSessionResponseV1_ValidationError_{
				ValidationError: &inout.CreateSessionResponseV1_ValidationError{
					PhoneCode: []string{"Не найден код подтверждения для данного телефона"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}
