package functools

func LimitPageToLimitOffset(limit, page int) (int, int) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 1
	}

	offset := page*limit - limit

	return limit, offset
}

func HasNext(totalCount int64, limit, page int) bool {
	limit, offset := LimitPageToLimitOffset(limit, page)
	return int64(limit)+int64(offset) < totalCount
}

func HasPrevious(page int) bool {
	return page > 1
}
