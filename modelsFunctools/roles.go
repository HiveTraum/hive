package modelsFunctools

import (
	"auth/functools"
	"auth/models"
)

func RoleIDListToInt64List(id []models.RoleID) []int64 {
	identifiers := make([]int64, len(id))

	for i, v := range id {
		identifiers[i] = int64(v)
	}

	return identifiers
}

func RoleIDListToString(id []models.RoleID, delimiter string) string {
	return functools.Int64SliceToString(RoleIDListToInt64List(id), delimiter)
}

func RoleIDListToPGArray(id []models.RoleID) string {
	return functools.Int64ListToPGArray(RoleIDListToInt64List(id))
}

func Int64SliceToRoleIDSlice(id []int64) []models.RoleID {
	identifiers := make([]models.RoleID, len(id))

	for i, v := range id {
		identifiers[i] = models.RoleID(v)
	}

	return identifiers
}

func StringsSliceToRoleIDSlice(str []string) []models.RoleID {
	return Int64SliceToRoleIDSlice(functools.StringsSliceToInt64String(str))
}
