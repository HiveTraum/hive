package repositories

import (
	"auth/config"
	"auth/enums"
	"auth/models"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUserRole(t *testing.T) {
	pool := config.InitPool(nil)
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
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	status, userRole := CreateUserRole(pool, ctx, user.Id, 1)
	require.Equal(t, enums.RoleNotFound, status)
	require.Nil(t, userRole)
}

func TestCreateUserRoleWithoutUser(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	status, userRole := CreateUserRole(pool, ctx, 1, role.Id)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, userRole)
}

func TestCreateUserRoleWithoutUserAndRole(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	status, userRole := CreateUserRole(pool, ctx, 1, 1)
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, userRole)
}

func TestCreateUserRoleThatAlreadyExist(t *testing.T) {
	pool := config.InitPool(nil)
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
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, adminRole := CreateRole(pool, ctx, "admin")
	CreateUserRole(pool, ctx, user.Id, role.Id)
	CreateUserRole(pool, ctx, user.Id, adminRole.Id)
	userRoles := GetUserRoles(pool, ctx, GetUserRoleQuery{
		UserId: []models.UserID{user.Id},
		RoleId: nil,
		Limit:  10,
	})
	require.Len(t, userRoles, 2)
}

func TestDeleteUserRole(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUserRoles(pool, ctx)
	PurgeRoles(pool, ctx)
	PurgeUsers(pool, ctx)
	user := CreateUser(pool, ctx)
	_, role := CreateRole(pool, ctx, "role")
	_, userRole := CreateUserRole(pool, ctx, user.Id, role.Id)
	DeleteUserRole(pool, ctx, userRole.Id)
	userRoles := GetUserRoles(pool, ctx, GetUserRoleQuery{
		UserId: []models.UserID{user.Id},
		RoleId: nil,
		Limit:  10,
	})

	require.Len(t, userRoles, 0)
}
