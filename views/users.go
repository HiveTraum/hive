package views

import (
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

func getUsersViewV1Query(r *functools.Request) repositories.GetUsersViewStoreQuery {
	query := r.URL.Query()
	return repositories.GetUsersViewStoreQuery{
		Limit:  r.GetLimit(),
		Id:     modelsFunctools.StringsSliceToUserIDSlice(query["id"]),
		Roles:  modelsFunctools.StringsSliceToRoleIDSlice(query["roles"]),
		Phones: query["phones"],
		Emails: query["emails"],
	}
}

func getUsersViewV1(r *functools.Request, app infrastructure.AppInterface) (int, *inout.ListUserViewResponseV1) {
	query := getUsersViewV1Query(r)
	users := app.GetStore().GetUsersView(r.Context(), query)

	userViews := make([]*inout.GetUserViewResponseV1, len(users))

	for i, u := range users {
		userViews[i] = &inout.GetUserViewResponseV1{
			Id:      u.Id.Bytes(),
			Created: u.Created,
			Roles:   u.Roles,
			Phones:  u.Phones,
			Emails:  u.Emails,
		}
	}

	return http.StatusOK, &inout.ListUserViewResponseV1{Data: userViews}
}

func getUserViewV1(r *functools.Request, app infrastructure.AppInterface, id uuid.UUID) (int, *inout.GetUserViewResponseV1) {

	userView := app.GetStore().GetUserView(r.Context(), id)

	if userView == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, &inout.GetUserViewResponseV1{
		Id:      userView.Id.Bytes(),
		Created: userView.Created,
		Roles:   userView.Roles,
		Phones:  userView.Phones,
		Emails:  userView.Emails,
	}
}

func UsersViewV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(r *functools.Request) (int, proto.Message) {
		return getUsersViewV1(r, app)
	}
}

func UserViewV1(app infrastructure.AppInterface) middlewares.ResponseControllerHandler {
	return func(request *functools.Request) (i int, message proto.Message) {
		vars := mux.Vars(request.Request)
		id, err := uuid.FromString(vars["id"])

		if err != nil {
			// Сознательно отправляем отчет об ошибке, т.к. в vars["id"] не должны попасть не числовые значения.
			// Если такое произошло - что то пошло не так
			sentry.CaptureException(err)
			return http.StatusBadRequest, nil
		}

		return getUserViewV1(request, app, id)
	}
}
