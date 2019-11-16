package functools

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
