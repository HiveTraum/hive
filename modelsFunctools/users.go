package modelsFunctools

import (
	"auth/functools"
	uuid "github.com/satori/go.uuid"
	"strings"
)

func UserIDListToStringList(id []uuid.UUID) []string {
	identifiers := make([]string, len(id))

	for i, v := range id {
		identifiers[i] = v.String()
	}

	return identifiers
}

func UserIDListToString(id []uuid.UUID, delimiter string) string {
	return strings.Join(UserIDListToStringList(id), delimiter)
}

func UserIDListToPGArray(id []uuid.UUID) string {
	return functools.StringsToPGArray(UserIDListToStringList(id))
}

func StringsSliceToUserIDSlice(id []string) []uuid.UUID {
	identifiers := make([]uuid.UUID, len(id))

	for i, v := range id {
		identifiers[i] = uuid.FromStringOrNil(v)
	}

	return identifiers
}
