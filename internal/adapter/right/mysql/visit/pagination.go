package mysqlvisitadapter

// defaultPagination normalizes and sanitizes pagination parameters
//
// This function ensures pagination parameters are within safe bounds to prevent:
//   - Memory exhaustion (excessive page size)
//   - Invalid page numbers (zero or negative)
//   - Denial of Service attacks (requesting millions of records)
//
// Normalization Rules:
//   - Limit: If ≤ 0 or > max, set to max (default protection)
//   - Page: If ≤ 0, set to 1 (1-indexed pagination)
//   - Offset: Calculated as (page - 1) * limit (zero-indexed database offset)
//
// Parameters:
//   - limit: Requested page size (0 means "use default")
//   - page: Requested page number (1-indexed, 0 means "first page")
//   - max: Maximum allowed page size (security/performance limit)
//
// Returns:
//   - limit: Normalized page size (guaranteed > 0 and ≤ max)
//   - offset: Database offset for LIMIT/OFFSET query (zero-indexed)
//
// Examples:
//
//	defaultPagination(0, 0, 50)    → (50, 0)    // First page, default size
//	defaultPagination(20, 1, 50)   → (20, 0)    // First page, custom size
//	defaultPagination(20, 3, 50)   → (20, 40)   // Third page: offset = (3-1)*20
//	defaultPagination(100, 2, 50)  → (50, 50)   // Limit capped to max=50
//	defaultPagination(-10, -5, 50) → (50, 0)    // Invalid inputs normalized
//
// Usage:
//
//	limit, offset := defaultPagination(filter.Limit, filter.Page, visitsMaxPageSize)
//	query := "SELECT * FROM visits LIMIT ? OFFSET ?"
func defaultPagination(limit, page, max int) (int, int) {
	// Normalize limit: use max if invalid or exceeds maximum
	if limit <= 0 || limit > max {
		limit = max
	}

	// Normalize page: default to first page if invalid
	if page <= 0 {
		page = 1
	}

	// Calculate zero-indexed offset for database query
	// Formula: offset = (page - 1) * limit
	// Example: page=1 → offset=0, page=2 → offset=limit, etc.
	return limit, (page - 1) * limit
}
