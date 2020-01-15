package modelsFunctools

import (
	"auth/functools"
	"auth/models"
	"strings"
)

func UserIDListToStringList(id []models.UserID) []string {
	identifiers := make([]string, len(id))

	for i, v := range id {
		identifiers[i] = string(v)
	}

	return identifiers
}

func UserIDListToString(id []models.UserID, delimiter string) string {
	return strings.Join(UserIDListToStringList(id), delimiter)
}

func UserIDListToPGArray(id []models.UserID) string {
	return functools.StringsToPGArray(UserIDListToStringList(id))
}

func Int64SliceToUserIDSlice(id []int64) []models.UserID {
	identifiers := make([]models.UserID, len(id))

	for i, v := range id {
		identifiers[i] = models.UserID(v)
	}

	return identifiers
}

func StringsSliceToUserIDSlice(str []string) []models.UserID {
	return Int64SliceToUserIDSlice(functools.StringsSliceToInt64String(str))
}
