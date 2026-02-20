package handler

import (
	"github.com/gin-gonic/gin"
)

// APIError represents a standardized error response.
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// RespondWithError sends a standardized JSON error response.
func RespondWithError(c *gin.Context, code int, errorCode string, message string) {
	resp := APIError{
		Error:   errorCode,
		Message: message,
	}
	c.AbortWithStatusJSON(code, resp)
}
