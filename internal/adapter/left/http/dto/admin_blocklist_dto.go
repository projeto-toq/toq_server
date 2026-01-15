package dto

// AdminBlocklistListRequest captures filters for GET /admin/blocklist
// Pagination is coarse because Redis SCAN is used; ordering is not guaranteed.
type AdminBlocklistListRequest struct {
	Page     int64 `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64 `form:"pageSize,default=100" binding:"omitempty,min=1,max=500"`
}

// AdminBlocklistItemResponse represents a single blocked JTI entry.
type AdminBlocklistItemResponse struct {
	JTI       string `json:"jti"`
	ExpiresAt int64  `json:"expiresAt"` // Unix seconds
}

// AdminBlocklistListResponse bundles items and pagination metadata.
type AdminBlocklistListResponse struct {
	Items      []AdminBlocklistItemResponse `json:"items"`
	Pagination PaginationResponse           `json:"pagination"`
}

// AdminBlocklistAddRequest payload to add a JTI to the blocklist.
type AdminBlocklistAddRequest struct {
	JTI string `json:"jti" binding:"required,uuid4"`
	TTL int64  `json:"ttl" binding:"required,min=1"` // seconds
}

// AdminBlocklistDeleteRequest payload to remove a JTI from the blocklist.
type AdminBlocklistDeleteRequest struct {
	JTI string `json:"jti" binding:"required,uuid4"`
}
