package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct{}

// NewContext creates a new context with timeout
func (b *BaseHandler) NewContext(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), timeout)
}

// ResponseJSON sends a JSON response with the given status code and data
func (b *BaseHandler) ResponseJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// ResponseError sends an error response with the given status code and error message
func (b *BaseHandler) ResponseError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// ResponseUnauthorized sends an unauthorized error response
func (b *BaseHandler) ResponseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
}
