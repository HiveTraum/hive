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
)

func GetUserRolesV1Query(r *functools.Request) repositories.GetUserRoleQuery {

	query := r.URL.Query()
	return repositories.GetUserRoleQuery{
		Pagination: functools.GetPagination(query),
		UserId:     functools.StringsSliceToUUIDSlice(query["users"]),
		RoleId:     functools.StringsSliceToUUIDSlice(query["roles"]),
	}
}

func getUserRolesV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserRolesResponseV1) {

	query := GetUserRolesV1Query(r)
	userRoles, pagination := app.GetStore().GetUserRoles(r.Context(), query)
	usersData := make([]*inout.UserRole, len(userRoles))

	for i, userRole := range userRoles {
		usersData[i] = &inout.UserRole{
			Id:      userRole.Id.Bytes(),
			Created: userRole.Created,
			UserID:  userRole.UserId.Bytes(),
			RoleID:  userRole.RoleId.Bytes(),
		}
	}

	return http.StatusOK, &inout.ListUserRolesResponseV1{Data: usersData, Pagination: &inout.Pagination{
		HasPrevious: pagination.HasPrevious,
		HasNext:     pagination.HasNext,
		Count:       pagination.Count,
	}}
}

func createUserRoleV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.CreateUserRoleResponseV1) {

	body := inout.CreateUserRoleResponseV1_Request{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, userRole := controllers.CreateUserRole(app.GetStore(), app.GetESB(), r.Context(), uuid.FromBytesOrNil(body.UserID), uuid.FromBytesOrNil(body.RoleID))

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_Ok{Ok: &inout.UserRole{
				Id:      userRole.Id.Bytes(),
				Created: userRole.Created,
				UserID:  userRole.UserId.Bytes(),
				RoleID:  userRole.RoleId.Bytes(),
			}},
		}
	case enums.RoleNotFound:
		return http.StatusBadRequest, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserRoleResponseV1_ValidationError{
					RoleID: []string{"Такой роли не существует"},
				}},
		}
	case enums.UserNotFound:
		return http.StatusBadRequest, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserRoleResponseV1_ValidationError{
					UserID: []string{"Такого пользователя не существует"},
				}}}
	case enums.UserRoleAlreadyExist:
		return http.StatusBadRequest, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserRoleResponseV1_ValidationError{
					Errors: []string{"Данная роль уже есть у пользователя"},
				}}}
	default:
		return unhandledStatus(r, status), nil
	}
}

func deleteUserRoleV1(r *functools.Request, app infrastructure.AppInterface, id uuid.UUID) (int, proto.Message) {

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
		id, err := uuid.FromString(vars["id"])

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		return deleteUserRoleV1(r, app, id)
	}
}
