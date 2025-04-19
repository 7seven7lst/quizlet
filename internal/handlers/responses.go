package handlers

// ErrorResponse represents an error response
// @model ErrorResponse
// @Description Error response from the API
type ErrorResponse struct {
	// The error message
	// @example "quiz attempt not found"
	Error string `json:"error" example:"quiz attempt not found"`
}

// SuccessResponse represents a success response
// @model SuccessResponse
// @Description Success response from the API
type SuccessResponse struct {
	// The success message
	// @example "Operation completed successfully"
	Message string `json:"message" example:"Operation completed successfully"`
} 