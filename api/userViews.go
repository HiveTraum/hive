package api

import (
	"auth/extractors"
	"auth/functools"
	"auth/inout"
	"auth/repositories"
	"net/http"
)

func getUsersViewV1Query(r *http.Request) repositories.GetUsersViewStoreQuery {
	query := r.URL.Query()
	pagination := functools.GetPagination(query)
	return repositories.GetUsersViewStoreQuery{
		Limit:  pagination.Limit,
		Page:   pagination.Page,
		Id:     functools.StringsSliceToUUIDSlice(query["id"]),
		Roles:  functools.StringsSliceToUUIDSlice(query["roles"]),
		Phones: query["phones"],
		Emails: query["emails"],
	}
}

func (api *API) GetUsersViewV1(w http.ResponseWriter, r *http.Request) {
	query := getUsersViewV1Query(r)
	users, pagination := api.Controller.GetUserViews(r.Context(), query)

	userViews := make([]*inout.UserView, len(users))

	for i, u := range users {
		userViews[i] = &inout.UserView{
			Id:      u.Id.Bytes(),
			Created: u.Created,
			Roles:   u.Roles,
			Phones:  u.Phones,
			Emails:  u.Emails,
		}
	}

	api.Renderer.Render(w, r, http.StatusOK, &inout.ListUserViewResponseV1{Data: userViews, Pagination: &inout.Pagination{
		HasPrevious: pagination.HasPrevious,
		HasNext:     pagination.HasNext,
		Count:       pagination.Count,
	}})
}

func (api *API) GetUserViewV1(w http.ResponseWriter, r *http.Request) {

	id, _ := extractors.GetUUID(r)
	userView := api.Controller.GetUserView(r.Context(), id)

	if userView == nil {
		api.Renderer.Render(w, r, http.StatusNotFound, nil)
	} else {
		api.Renderer.Render(w, r, http.StatusOK, &inout.GetUserViewResponseV1{
			Data: &inout.UserView{
				Id:      userView.Id.Bytes(),
				Created: userView.Created,
				Roles:   userView.Roles,
				Phones:  userView.Phones,
				Emails:  userView.Emails,
			}})
	}
}
