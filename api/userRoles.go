package api

import (
	"auth/enums"
	"auth/extractors"
	"auth/functools"
	"auth/inout"
	"auth/repositories"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func GetUserRolesV1Query(r *http.Request) repositories.GetUserRoleQuery {

	query := r.URL.Query()
	return repositories.GetUserRoleQuery{
		Pagination: functools.GetPagination(query),
		UserId:     functools.StringsSliceToUUIDSlice(query["users"]),
		RoleId:     functools.StringsSliceToUUIDSlice(query["roles"]),
	}
}

func (api *API) GetUserRolesV1(w http.ResponseWriter, r *http.Request) {

	query := GetUserRolesV1Query(r)
	userRoles, pagination := api.Controller.GetUserRoles(r.Context(), query)
	usersData := make([]*inout.UserRole, len(userRoles))

	for i, userRole := range userRoles {
		usersData[i] = &inout.UserRole{
			Id:      userRole.Id.Bytes(),
			Created: userRole.Created,
			UserID:  userRole.UserId.Bytes(),
			RoleID:  userRole.RoleId.Bytes(),
		}
	}

	api.Renderer.Render(w, r, http.StatusOK, &inout.ListUserRolesResponseV1{Data: usersData, Pagination: &inout.Pagination{
		HasPrevious: pagination.HasPrevious,
		HasNext:     pagination.HasNext,
		Count:       pagination.Count,
	}})
}

func (api *API) CreateUserRoleV1(w http.ResponseWriter, r *http.Request) {

	body := &inout.CreateUserRoleResponseV1_Request{}
	err := api.Parser.Parse(r, w, body)
	if err != nil {
		return
	}

	status, userRole := api.Controller.CreateUserRole(r.Context(), uuid.FromBytesOrNil(body.UserID), uuid.FromBytesOrNil(body.RoleID))

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusCreated, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_Ok{Ok: &inout.UserRole{
				Id:      userRole.Id.Bytes(),
				Created: userRole.Created,
				UserID:  userRole.UserId.Bytes(),
				RoleID:  userRole.RoleId.Bytes(),
			}},
		})
	case enums.RoleNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserRoleResponseV1_ValidationError{
					RoleID: []string{"Такой роли не существует"},
				}},
		})
	case enums.UserNotFound:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserRoleResponseV1_ValidationError{
					UserID: []string{"Такого пользователя не существует"},
				}}})
	case enums.UserRoleAlreadyExist:
		api.Renderer.Render(w, r, http.StatusBadRequest, &inout.CreateUserRoleResponseV1{
			Data: &inout.CreateUserRoleResponseV1_ValidationError_{
				ValidationError: &inout.CreateUserRoleResponseV1_ValidationError{
					Errors: []string{"Данная роль уже есть у пользователя"},
				}}})
	default:
		api.Renderer.Render(w, r, unhandledStatus(r, status), nil)
	}
}

func (api *API) DeleteUserRoleV1(w http.ResponseWriter, r *http.Request) {

	id, _ := extractors.GetUUID(r)
	status, _ := api.Controller.DeleteUserRole(r.Context(), id)

	switch status {
	case enums.Ok:
		api.Renderer.Render(w, r, http.StatusNoContent, nil)
	case enums.UserRoleNotFound:
		api.Renderer.Render(w, r, http.StatusNoContent, nil)
	default:
		api.Renderer.Render(w, r, http.StatusNoContent, nil)
	}
}
