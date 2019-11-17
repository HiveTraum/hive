package api

import (
	"auth/controllers"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"auth/repositories"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strconv"
)

func GetUserRolesV1Query(r *http.Request) repositories.GetUserRoleQuery {
	limitQuery := r.URL.Query().Get("limit")
	if limitQuery == "" {
		limitQuery = string(DefaultLimit)
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		limit = DefaultLimit
	}

	return repositories.GetUserRoleQuery{
		Limit:  limit,
		UserId: functools.StringsSliceToInt64String(r.URL.Query()["users"]),
		RoleId: functools.StringsSliceToInt64String(r.URL.Query()["roles"]),
	}
}

func getUserRolesV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserRolesResponseV1) {

	query := GetUserRolesV1Query(r.Request)
	userRoles := app.GetStore().GetUserRoles(r.Context(), query)
	usersData := make([]*inout.GetUserRoleResponseV1, len(userRoles))

	for i, userRole := range userRoles {
		usersData[i] = &inout.GetUserRoleResponseV1{
			Created: userRole.Created,
			UserId:  userRole.UserId,
			RoleId:  userRole.RoleId,
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

	status, userRole := controllers.CreateUserRole(app.GetStore(), app.GetESB(), r.Context(), body.UserId, body.RoleId)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.GetUserRoleResponseV1{
			Created: userRole.Created,
			UserId:  userRole.UserId,
			RoleId:  userRole.RoleId,
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
			Created: userRole.Created,
			UserId:  userRole.UserId,
			RoleId:  userRole.RoleId,
		}
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
