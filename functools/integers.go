package functools

import (
	"fmt"
	"strings"
)

func All(value int, list []int) bool {
	for _, e := range list {
		if e != value {
			return false
		}
	}

	return true
}

func Max(list []int) int {
	max := list[0]
	for _, e := range list {
		if e > max {
			max = e
		}
	}

	return max
}

func Int64SliceToString(slice []int64, delimiter string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", delimiter, -1), "[]")
}
