package api

import (
	"auth/controllers"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"auth/repositories"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
)

func createUserV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := &inout.CreateUserRequestV1{}
	err := r.ParseBody(body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, user := controllers.CreateUser(app.GetStore(), app.GetESB(), app.GetPasswordProcessor(), r.Context(), body)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.GetUserResponseV1{
			Id:      user.Id.Bytes(),
			Created: user.Created,
		}
	case enums.MinimumOneFieldRequired:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Errors: []string{"Необходимо указать телефон или почту"},
		}

	// Password validations

	case enums.PasswordRequired:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Password: []string{"Необходимо указать пароль"},
		}
	case enums.IncorrectPassword:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Password: []string{"Не удалось обработать полученный пароль, попробуйте другой"},
		}

	// Phone Validations

	case enums.IncorrectPhoneCode:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			PhoneCode: []string{"Неверный код"},
		}
	case enums.PhoneNotFound:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Phone: []string{"Не удалось найти код для данного телефона."},
		}
	case enums.IncorrectPhone:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Phone: []string{"Некорректный номер телефона"},
		}

	// Email Validations

	case enums.IncorrectEmailCode:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			EmailCode: []string{"Неверный код"},
		}
	case enums.EmailNotFound:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Email: []string{"Не удалось найти код для данного email."},
		}
	case enums.IncorrectEmail:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Email: []string{"Некорректный email"},
		}

	default:
		return unhandledStatus(r, status)
	}
}

func GetUsersV1Query(query url.Values, payload infrastructure.AuthenticationBackendUser) repositories.GetUsersQuery {
	pagination := functools.GetPagination(query)

	var requestedUserIdentifiers []uuid.UUID
	if payload.GetIsAdmin() {
		requestedUserIdentifiers = functools.StringsSliceToUUIDSlice(query["id"])
	} else {
		requestedUserIdentifiers = []uuid.UUID{payload.GetUserID()}
	}

	return repositories.GetUsersQuery{
		Limit: pagination.Limit,
		Page:  pagination.Page,
		Id:    requestedUserIdentifiers,
	}
}

func getUsersV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserResponseV1) {

	status, payload := app.GetLoginController().Login(r.Context(), r.GetAuthorizationHeader())
	if status != enums.Ok || payload == nil {
		return http.StatusUnauthorized, nil
	}

	query := GetUsersV1Query(r.URL.Query(), payload)
	users := app.GetStore().GetUsers(r.Context(), query)
	usersData := make([]*inout.GetUserResponseV1, len(users))

	for i, user := range users {
		usersData[i] = &inout.GetUserResponseV1{
			Id:      user.Id.Bytes(),
			Created: user.Created,
		}
	}

	return http.StatusOK, &inout.ListUserResponseV1{Data: usersData}
}

func getUserV1(r *functools.Request, app infrastructure.AppInterface, id uuid.UUID) (int, *inout.GetUserResponseV1) {
	user := app.GetStore().GetUser(r.Context(), id)

	if user == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, &inout.GetUserResponseV1{
		Id:      user.Id.Bytes(),
		Created: user.Created,
	}
}

func deleteUserV1(r *functools.Request, app infrastructure.AppInterface, id uuid.UUID) (int, proto.Message) {

	ctx := r.Context()

	status, deletedUser := controllers.DeleteUser(app.GetStore(), app.GetESB(), ctx, id)

	switch status {
	case enums.Ok:
		return http.StatusOK, &inout.GetUserResponseV1{
			Id:      deletedUser.Id.Bytes(),
			Created: deletedUser.Created,
		}
	case enums.UserNotFound:
		return http.StatusNotFound, nil
	default:
		return unhandledStatus(r, status)
	}
}

func UsersAPIV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		if r.Method == "GET" {
			return getUsersV1(r, app)
		} else {
			return createUserV1(r, app)
		}
	}
}

func UserAPIV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(request *functools.Request) (i int, message proto.Message) {

		vars := mux.Vars(request.Request)
		id, err := uuid.FromString(vars["id"])

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		status, payload := app.GetLoginController().Login(request.Context(), request.GetAuthorizationHeader())
		if status != enums.Ok {
			return http.StatusUnauthorized, nil
		}
		if id != payload.GetUserID() || payload.GetIsAdmin() {
			return http.StatusForbidden, nil
		}

		if request.Method == "DELETE" {
			return deleteUserV1(request, app, id)
		} else {
			return getUserV1(request, app, id)
		}
	}
}
