package modelsFunctools

import (
	"auth/functools"
	"auth/models"
)

func UserIDListToInt64List(id []models.UserID) []int64 {
	identifiers := make([]int64, len(id))

	for i, v := range id {
		identifiers[i] = int64(v)
	}

	return identifiers
}

func UserIDListToString(id []models.UserID, delimiter string) string {
	return functools.Int64SliceToString(UserIDListToInt64List(id), delimiter)
}

func UserIDListToPGArray(id []models.UserID) string {
	return functools.Int64ListToPGArray(UserIDListToInt64List(id))
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
