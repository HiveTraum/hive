package modelsFunctools

import (
	"auth/functools"
	uuid "github.com/satori/go.uuid"
)

func RoleIDListToStringList(id []uuid.UUID) []string {
	identifiers := make([]string, len(id))

	for i, v := range id {
		identifiers[i] = v.String()
	}

	return identifiers
}

func RoleIDListToPGArray(id []uuid.UUID) string {
	return functools.StringListToPGArray(RoleIDListToStringList(id))
}

func StringsSliceToRoleIDSlice(id []string) []uuid.UUID {
	identifiers := make([]uuid.UUID, len(id))

	for i, v := range id {
		identifiers[i] = uuid.FromStringOrNil(v)
	}

	return identifiers
}

