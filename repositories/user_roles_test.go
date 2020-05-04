package repositories

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUserRole(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	status, userRole := CreateUserRole(pool, ctx, user.Id, role.Id)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, userRole)
}

func TestCreateUserRoleWithoutRole(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	status, userRole := CreateUserRole(pool, ctx, user.Id, uuid.NewV4())
	require.Equal(t, enums.RoleNotFound, status)
	require.Nil(t, userRole)
}

func TestCreateUserRoleWithoutUser(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	status, userRole := CreateUserRole(pool, ctx, uuid.NewV4(), role.Id)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, userRole)
}

func TestCreateUserRoleWithoutUserAndRole(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	status, userRole := CreateUserRole(pool, ctx, uuid.NewV4(), uuid.NewV4())
	require.True(t, functools.In([]int{enums.RoleNotFound, enums.UserNotFound}, status))
	require.Nil(t, userRole)
}

func TestCreateUserRoleThatAlreadyExist(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	CreateUserRole(pool, ctx, user.Id, role.Id)
	status, userRole := CreateUserRole(pool, ctx, user.Id, role.Id)
	require.Equal(t, enums.UserRoleAlreadyExist, status)
	require.Nil(t, userRole)
}

func TestGetUserRolesWithTwoRoles(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, adminRole := CreateRole(pool, ctx, "admin")
	CreateUserRole(pool, ctx, user.Id, role.Id)
	CreateUserRole(pool, ctx, user.Id, adminRole.Id)
	userRoles, _ := GetUserRoles(pool, ctx, GetUserRoleQuery{
		UserId:     []uuid.UUID{user.Id},
		RoleId:     nil,
		Pagination: &models.PaginationRequest{Limit: 10},
	})
	require.Len(t, userRoles, 2)
}

func TestDeleteUserRole(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, userRole := CreateUserRole(pool, ctx, user.Id, role.Id)
	DeleteUserRole(pool, ctx, userRole.Id)
	userRoles, _ := GetUserRoles(pool, ctx, GetUserRoleQuery{
		UserId:     []uuid.UUID{user.Id},
		RoleId:     nil,
		Pagination: &models.PaginationRequest{Limit: 10},
	})

	require.Len(t, userRoles, 0)
}

func TestGetUserRolesWithLimitedPagination(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, lore := CreateRole(pool, ctx, "lore")
	_, userRole := CreateUserRole(pool, ctx, user.Id, role.Id)
	CreateUserRole(pool, ctx, user.Id, lore.Id)
	userRoles, _ := GetUserRoles(pool, ctx, GetUserRoleQuery{
		Pagination: &models.PaginationRequest{Limit: 1},
	})

	require.Len(t, userRoles, 1)
	require.Equal(t, userRole.Id, userRoles[0].Id)
}

func TestGetUserRolesWithPagination(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, lore := CreateRole(pool, ctx, "lore")
	CreateUserRole(pool, ctx, user.Id, role.Id)
	CreateUserRole(pool, ctx, user.Id, lore.Id)
	userRoles, _ := GetUserRoles(pool, ctx, GetUserRoleQuery{
		Pagination: &models.PaginationRequest{Limit: 10},
	})

	require.Len(t, userRoles, 2)
}

func TestGetUserRolesWithLimitedPaginationAndSecondPage(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, lore := CreateRole(pool, ctx, "lore")
	CreateUserRole(pool, ctx, user.Id, role.Id)
	_, userRole := CreateUserRole(pool, ctx, user.Id, lore.Id)
	userRoles, _ := GetUserRoles(pool, ctx, GetUserRoleQuery{
		Pagination: &models.PaginationRequest{Limit: 1, Page: 2},
	})

	require.Len(t, userRoles, 1)
	require.Equal(t, userRole.Id, userRoles[0].Id)
}
