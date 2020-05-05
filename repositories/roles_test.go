package repositories

import (
	"hive/config"
	"hive/enums"
	"hive/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateRole(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeRoles(pool, ctx)
	status, role := CreateRole(pool, ctx, "admin")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, role)
	roles, _ := GetRoles(pool, ctx, GetRolesQuery{
		Pagination:  &models.PaginationRequest{Limit: 10},
		Identifiers: nil,
	})
	require.Len(t, roles, 1)
}

func TestCreateRoleThatAlreadyExist(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeRoles(pool, ctx)
	_, _ = CreateRole(pool, ctx, "admin")
	status, role := CreateRole(pool, ctx, "admin")
	require.Equal(t, enums.RoleAlreadyExist, status)
	require.Nil(t, role)
	roles, _ := GetRoles(pool, ctx, GetRolesQuery{
		Pagination:  &models.PaginationRequest{Limit: 10},
		Identifiers: nil,
	})
	require.Len(t, roles, 1)
}

func TestGetRolesWithEmptyTable(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeRoles(pool, ctx)
	roles, _ := GetRoles(pool, ctx, GetRolesQuery{
		Pagination:  &models.PaginationRequest{Limit: 10},
		Identifiers: nil,
	})
	require.Len(t, roles, 0)
}

func TestGetRole(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeRoles(pool, ctx)
	_, createdRole := CreateRole(pool, ctx, "admin")
	status, role := GetRole(pool, ctx, createdRole.Id)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, role)
	require.Equal(t, createdRole.Id, role.Id)
	require.Equal(t, "admin", role.Title)
}

func TestGetRoleWithEmptyTable(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeRoles(pool, ctx)
	status, role := GetRole(pool, ctx, uuid.NewV4())
	require.Equal(t, enums.RoleNotFound, status)
	require.Nil(t, role)
}
