package dto

// Common response structures

// TokensResponse represents authentication tokens
type TokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty" swaggertype:"object"`
}

// SuccessResponse creates a success API response
func SuccessResponse(data interface{}, message ...string) *APIResponse {
	resp := &APIResponse{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}

// ErrorResponseFromError creates an error API response
func ErrorResponseFromError(err error) *APIResponse {
	return &APIResponse{
		Success: false,
		Error:   err.Error(),
	}
}

// ErrorResponseWithDetails creates an error API response with details
func ErrorResponseWithDetails(code int, message string, details interface{}) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
