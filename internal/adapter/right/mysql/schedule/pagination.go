package mysqlscheduleadapter

// defaultPagination clamps limit to (0, max] and normalizes page to >=1 returning limit and offset.
// Used by listing queries to avoid unbounded scans and to enforce adapter-side safety limits.
// Parameters:
//   - limit: requested page size; values <=0 or above max are replaced by max.
//   - page: requested page number (1-indexed); values <=0 default to 1.
//   - max: hard upper bound for page size enforced at adapter level.
//
// Returns: normalized limit and offset suitable for LIMIT/OFFSET clauses.
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
