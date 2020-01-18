package modelsFunctools

import (
	"auth/functools"
	uuid "github.com/satori/go.uuid"
)

func SecretIDListToStringList(id []uuid.UUID) []string {
	identifiers := make([]string, len(id))

	for i, v := range id {
		identifiers[i] = v.String()
	}

	return identifiers
}

func SecretIDListToPGArray(id []uuid.UUID) string {
	return functools.StringsToPGArray(SecretIDListToStringList(id))
}
