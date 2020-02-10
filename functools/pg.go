package functools

import (
	"strings"
)

func StringsToPGArray(list []string) string {
	return "{" + strings.Join(list[:], ",") + "}"
}
