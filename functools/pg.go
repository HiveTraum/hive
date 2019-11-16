package functools

import (
	"strconv"
	"strings"
)

func StringsToPGArray(list []string) string {
	return "{" + strings.Join(list[:], ",") + "}"
}

func Int64ListToPGArray(list []int64) string {
	result := make([]string, len(list))
	for i, id := range list {
		result[i] = strconv.Itoa(int(id))
	}

	return StringsToPGArray(result)
}
