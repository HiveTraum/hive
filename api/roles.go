package api

import (
	"auth/app"
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

func GetRolesV1Query(r *functools.Request) repositories.GetRolesQuery {
	query := r.URL.Query()
	return repositories.GetRolesQuery{
		Pagination:  functools.GetPagination(query),
		Identifiers: query["id"],
	}
}

func getRoleV1(r *functools.Request, app *app.App, id uuid.UUID) (int, proto.Message) {
	status, role := app.Store.GetRole(r.Context(), id)

	switch status {
	case enums.Ok:
		return http.StatusOK, &inout.GetRoleResponseV1{
			Id:      role.Id.Bytes(),
			Created: role.Created,
			Title:   role.Title,
		}
	case enums.RoleNotFound:
		return http.StatusNotFound, nil
	default:
		return unhandledStatus(r, status)
	}
}

func getRolesV1(r *functools.Request, app *app.App) (int, *inout.ListRoleResponseV1) {

	query := GetRolesV1Query(r)
	roles := app.Store.GetRoles(r.Context(), query)
	rolesData := make([]*inout.GetRoleResponseV1, len(roles))

	for i, role := range roles {
		rolesData[i] = &inout.GetRoleResponseV1{
			Id:      role.Id.Bytes(),
			Created: role.Created,
			Title:   role.Title,
		}
	}

	return http.StatusOK, &inout.ListRoleResponseV1{Data: rolesData}
}

func createRoleV1(r *functools.Request, app infrastructure.AppInterface) (int, proto.Message) {

	body := inout.CreateRoleRequestV1{}

	err := r.ParseBody(&body)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	status, role := controllers.CreateRole(app.GetStore(), app.GetESB(), r.Context(), body.Title)

	switch status {
	case enums.Ok:
		return http.StatusCreated, &inout.GetRoleResponseV1{
			Id:      role.Id.Bytes(),
			Created: role.Created,
			Title:   role.Title,
		}
	case enums.RoleAlreadyExist:
		return http.StatusBadRequest, &inout.CreateRoleBadRequestV1{
			Title: []string{"Роль с таким названием уже существует"},
		}
	default:
		return unhandledStatus(r, status)
	}
}

func RoleAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(request *functools.Request) (i int, message proto.Message) {

		vars := mux.Vars(request.Request)
		id, err := uuid.FromString(vars["id"])

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		return getRoleV1(request, app, id)
	}
}

func RolesAPIV1(app *app.App) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		if r.Method == "GET" {
			return getRolesV1(r, app)
		} else {
			return createRoleV1(r, app)
		}
	}
}
