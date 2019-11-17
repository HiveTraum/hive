package views

import (
	"auth/api"
	"auth/controllers"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/middlewares"
	"auth/repositories"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func getUsersViewV1Query(r *http.Request) repositories.GetUsersViewQuery {
	usersQuery := api.GetUsersV1Query(r)

	rolesQuery := r.URL.Query()["role_id"]
	var roles []int64

	if len(rolesQuery) > 0 {
		identifiersQueryInt := make([]int64, len(rolesQuery))
		for i, q := range rolesQuery {
			idQueryInt, _ := strconv.Atoi(q)
			identifiersQueryInt[i] = int64(idQueryInt)
		}

		roles = identifiersQueryInt
	}

	return repositories.GetUsersViewQuery{
		GetUsersQuery: usersQuery,
		Roles:         roles,
	}
}

func getUsersViewV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserViewResponseV1) {

	query := getUsersViewV1Query(r.Request)
	users := app.GetStore().GetUsersView(r.Context(), query)
	return http.StatusOK, &inout.ListUserViewResponseV1{Data: users}
}

func getUserViewV1(r *functools.Request, app infrastructure.AppInterface, id int64) (int, *inout.GetUserViewResponseV1) {

	userView := controllers.GetUserView(app.GetStore(), app.GetESB(), r.Context(), id)

	if userView == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, userView
}

func UsersViewV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return getUsersViewV1(r, app)
	}
}

func UserViewV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(request *functools.Request) (i int, message proto.Message) {
		vars := mux.Vars(request.Request)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		return getUserViewV1(request, app, id)
	}
}
