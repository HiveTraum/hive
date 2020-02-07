package api

import (
	"auth/controllers"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"auth/modelsFunctools"
	"auth/repositories"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func createUserV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := &inout.CreateUserRequestV1{}
	err := r.ParseBody(body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, user := controllers.CreateUser(app.GetStore(), app.GetESB(), app.GetLoginController(), r.Context(), body)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.GetUserResponseV1{
			Id:      user.Id.Bytes(),
			Created: user.Created,
		}
	case enums.MinimumOneFieldRequired:
		return http.StatusBadRequest, &inout.CreateUserBadRequestV1{
			Errors: []string{"Необходимо указать телефон или почту",},
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
		return http.StatusCreated, &inout.GetUserResponseV1{
			Id:      user.Id.Bytes(),
			Created: user.Created,
		}
	}
}

func GetUsersV1Query(r *functools.Request) repositories.GetUsersQuery {
	query := r.URL.Query()
	pagination := functools.GetPagination(query)

	return repositories.GetUsersQuery{
		Limit: pagination.Limit,
		Page:  pagination.Page,
		Id:    modelsFunctools.StringsSliceToUserIDSlice(r.URL.Query()["id"]),
	}
}

func getUsersV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserResponseV1) {

	query := GetUsersV1Query(r)
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

func deleteUserV1(r *functools.Request, app infrastructure.AppInterface, id uuid.UUID) (int, *inout.GetUserResponseV1) {

	ctx := r.Context()

	status, payload := app.GetLoginController().Login(ctx, r.GetAccessToken())
	if status != enums.Ok {
		return http.StatusUnauthorized, nil
	}
	if id != payload.UserID || payload.IsAdmin {
		return http.StatusForbidden, nil
	}

	_, deletedUser := app.GetStore().DeleteUser(ctx, id)
	if deletedUser == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusNoContent, &inout.GetUserResponseV1{
		Id:      deletedUser.Id.Bytes(),
		Created: deletedUser.Created,
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

		if request.Method == "DELETE" {
			return deleteUserV1(request, app, id)
		} else {
			return getUserV1(request, app, id)
		}
	}
}
