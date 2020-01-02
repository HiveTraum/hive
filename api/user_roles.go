package api

import (
	"auth/controllers"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"auth/models"
	"auth/repositories"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetUserRolesV1Query(r *functools.Request) repositories.GetUserRoleQuery {

	return repositories.GetUserRoleQuery{
		Limit:  r.GetLimit(),
		UserId: functools.StringsSliceToInt64String(r.URL.Query()["users"]),
		RoleId: functools.StringsSliceToInt64String(r.URL.Query()["roles"]),
	}
}

func getUserRolesV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserRolesResponseV1) {

	query := GetUserRolesV1Query(r)
	userRoles := app.GetStore().GetUserRoles(r.Context(), query)
	usersData := make([]*inout.GetUserRoleResponseV1, len(userRoles))

	for i, userRole := range userRoles {
		usersData[i] = &inout.GetUserRoleResponseV1{
			Id:      int64(userRole.Id),
			Created: userRole.Created,
			UserId:  int64(userRole.UserId),
			RoleId:  int64(userRole.RoleId),
		}
	}

	return http.StatusOK, &inout.ListUserRolesResponseV1{Data: usersData}
}

func createUserRoleV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := inout.CreateUserRoleRequestV1{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, userRole := controllers.CreateUserRole(app.GetStore(), app.GetESB(), r.Context(), models.UserID(body.UserId), models.RoleID(body.RoleId))

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.GetUserRoleResponseV1{
			Id:      int64(userRole.Id),
			Created: userRole.Created,
			UserId:  int64(userRole.UserId),
			RoleId:  int64(userRole.RoleId),
		}
	case enums.RoleNotFound:
		return http.StatusBadRequest, &inout.CreateUserRoleBadRequestV1{
			RoleId: []string{"Такой роли не существует"},
		}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreateUserRoleBadRequestV1{
			UserId: []string{"Такого пользователя не существует"},
		}
	case enums.UserRoleAlreadyExist:
		return http.StatusBadRequest, &inout.CreateUserRoleBadRequestV1{
			Errors: []string{"Данная роль уже есть у пользователя"},
		}
	default:
		return http.StatusCreated, &inout.GetUserRoleResponseV1{
			Id:      int64(userRole.Id),
			Created: userRole.Created,
			UserId:  int64(userRole.UserId),
			RoleId:  int64(userRole.RoleId),
		}
	}
}

func deleteUserRoleV1(r *functools.Request, app infrastructure.AppInterface, id models.UserRoleID) (int, proto.Message) {

	status, _ := controllers.DeleteUserRole(app.GetStore(), app.GetESB(), r.Context(), id)

	switch status {
	case enums.Ok:
		return http.StatusNoContent, nil
	case enums.UserRoleNotFound:
		return http.StatusNoContent, nil
	default:
		return http.StatusNoContent, nil
	}
}

func UserRolesAPIV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		switch r.Method {
		case http.MethodPost:
			return createUserRoleV1(r, app)
		default:
			return getUserRolesV1(r, app)
		}
	}
}

func UserRoleAPIV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		vars := mux.Vars(r.Request)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		return deleteUserRoleV1(r, app, models.UserRoleID(id))
	}
}
