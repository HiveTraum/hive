package repositories

import (
	"auth/inout"
	"context"
)

func createOrUpdateUserRolesSQL() string {
	return `INSERT INTO user_roles(user_id, role_id) VALUES ($1, $2);`
}

func CreateOrUpdateUserRole(db DB, ctx context.Context, request *inout.CreateUserRoleRequestV1) {

}
