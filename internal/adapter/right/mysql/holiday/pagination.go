package mysqlholidayadapter

func defaultPagination(limit, page, max int) (int, int) {
	if limit <= 0 || limit > max {
		limit = max
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return limit, offset
}
