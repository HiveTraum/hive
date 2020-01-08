package modelsFunctools

import (
	"auth/functools"
	"auth/models"
)

func SecretIDListToInt64List(id []models.SecretID) []int64 {
	identifiers := make([]int64, len(id))

	for i, v := range id {
		identifiers[i] = int64(v)
	}

	return identifiers
}

func SecretIDListToPGArray(id []models.SecretID) string {
	return functools.Int64ListToPGArray(SecretIDListToInt64List(id))
}
