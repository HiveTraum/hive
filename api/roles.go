package api

import (
	"auth/enums"
	"auth/extractors"
	"auth/functools"
	"auth/inout"
	"auth/repositories"
	"net/http"
)

func GetRolesV1Query(r *http.Request) repositories.GetRolesQuery {
	query := r.URL.Query()
	return repositories.GetRolesQuery{
		Pagination:  functools.GetPagination(query),
		Identifiers: query["id"],
	}
}

func (api *API) GetRoleV1(w http.ResponseWriter, r *http.Request) {
	id, _ := extractors.GetUUID(r)
	status, role := api.Controller.GetRole(r.Context(), id)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusOK, &inout.GetRoleResponseV1{
			Data: &inout.Role{
				Id:      role.Id.Bytes(),
				Created: role.Created,
				Title:   role.Title,
			}})
	case enums.RoleNotFound:
		api.Renderer.Render(w, r, http.StatusNotFound, &inout.GetRoleResponseV1{})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}

func (api *API) GetRolesV1(w http.ResponseWriter, r *http.Request) {

	query := GetRolesV1Query(r)
	roles, pagination := api.Controller.GetRoles(r.Context(), query)
	rolesData := make([]*inout.Role, len(roles))

	for i, role := range roles {
		rolesData[i] = &inout.Role{
			Id:      role.Id.Bytes(),
			Created: role.Created,
			Title:   role.Title,
		}
	}

	api.Renderer.Render(w, r, http.StatusOK, &inout.ListRoleResponseV1{Data: rolesData, Pagination: &inout.Pagination{
		HasPrevious: pagination.HasPrevious,
		HasNext:     pagination.HasNext,
		Count:       pagination.Count,
	}})
}

func (api *API) CreateRoleV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreateRoleResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	ctx := r.Context()

	status, role := api.Controller.CreateRole(ctx, body.Title)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreateRoleResponseV1{
			Data: &inout.CreateRoleResponseV1_Ok{
				Ok: &inout.Role{
					Id:      role.Id.Bytes(),
					Created: role.Created,
					Title:   role.Title,
				}}})
	case enums.RoleAlreadyExist:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateRoleResponseV1{
			Data: &inout.CreateRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateRoleResponseV1_ValidationError{
					Title: []string{"Роль с таким названием уже существует"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}
