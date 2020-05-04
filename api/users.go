package api

import (
	"auth/enums"
	"auth/extractors"
	"auth/functools"
	"auth/inout"
	"auth/models"
	"auth/repositories"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
)

func (api *API) CreateUserV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreateUserResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, user := api.Controller.CreateUser(r.Context(), body.Password, body.Email, body.EmailCode, body.Phone, body.PhoneCode)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_Ok{
				Ok: &inout.User{
					Id:      user.Id.Bytes(),
					Created: user.Created,
				}}})
	case enums.MinimumOneFieldRequired:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Errors: []string{"Необходимо указать телефон или почту"},
				}}})

	// Password validations

	case enums.PasswordRequired:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Password: []string{"Необходимо указать пароль"},
				}}})
	case enums.IncorrectPassword:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Password: []string{"Не удалось обработать полученный пароль, попробуйте другой"},
				}}})

	// Phone Validations

	case enums.IncorrectPhoneCode:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					PhoneCode: []string{"Неверный код"},
				}}})
	case enums.PhoneNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Phone: []string{"Не удалось найти код для данного телефона."},
				}}})
	case enums.IncorrectPhone:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Phone: []string{"Некорректный номер телефона"},
				}}})

	// Email Validations

	case enums.IncorrectEmailCode:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					EmailCode: []string{"Неверный код"},
				}}})
	case enums.EmailNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Email: []string{"Не удалось найти код для данного email."},
				}}})
	case enums.IncorrectEmail:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserResponseV1{
			Data: &inout.CreateUserResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserResponseV1_ValidationError{
					Email: []string{"Некорректный email"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}

func (api *API) GetUsersV1Query(query url.Values, user models.IAuthenticationBackendUser) repositories.GetUsersQuery {
	pagination := functools.GetPagination(query, api.environment)

	var requestedUserIdentifiers []uuid.UUID
	if user.GetIsAdmin() {
		requestedUserIdentifiers = functools.StringsSliceToUUIDSlice(query["id"])
	} else {
		requestedUserIdentifiers = []uuid.UUID{user.GetUserID()}
	}

	return repositories.GetUsersQuery{
		Limit: pagination.Limit,
		Page:  pagination.Page,
		Id:    requestedUserIdentifiers,
	}
}

func (api *API) GetUsersV1(w http.ResponseWriter, r *http.Request) {

	user := repositories.GetUserFromContext(r.Context())
	query := api.GetUsersV1Query(r.URL.Query(), user)
	users := api.Controller.GetUsers(r.Context(), query)
	usersData := make([]*inout.User, len(users))

	for i, user := range users {
		usersData[i] = &inout.User{
			Id:      user.Id.Bytes(),
			Created: user.Created,
		}
	}

	api.Renderer.Render(w, r, http.StatusOK, &inout.ListUserResponseV1{Data: usersData})
}

func (api *API) GetUserV1(w http.ResponseWriter, r *http.Request) {

	id, err := extractors.GetUUID(r)
	if err != nil {
		sentry.CaptureException(err)
		api.Renderer.Render(w, r, http.StatusBadRequest, nil)
		return
	}

	user := api.Controller.GetUser(r.Context(), id)

	if user == nil {
		api.Renderer.Render(w, r, http.StatusNotFound, nil)
		return
	}

	api.Renderer.Render(w, r, http.StatusOK, &inout.GetUserResponseV1{
		Data: &inout.User{
			Id:      user.Id.Bytes(),
			Created: user.Created,
		}},
	)
}

func (api *API) DeleteUserV1(w http.ResponseWriter, r *http.Request) {

	id, err := extractors.GetUUID(r)
	if err != nil {
		sentry.CaptureException(err)
		api.Renderer.Render(w, r, http.StatusBadRequest, nil)
		return
	}

	status, deletedUser := api.Controller.DeleteUser(r.Context(), id)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusOK, &inout.GetUserResponseV1{
			Data: &inout.User{
				Id:      deletedUser.Id.Bytes(),
				Created: deletedUser.Created,
			},
		})
	case enums.UserNotFound:
		api.Renderer.Render(w, r, http.StatusNotFound, nil)
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}
