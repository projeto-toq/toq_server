package mysqlvisitadapter

// defaultPagination normalizes limit and page parameters, returning limit and offset.
func defaultPagination(limit, page, max int) (int, int) {
	if limit <= 0 || limit > max {
		limit = max
	}
	if page <= 0 {
		page = 1
	}
	return limit, (page - 1) * limit
}
