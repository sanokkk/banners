package http

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

// baseAuthMiddleware set up middleware for checking if user has auth token
func (s *Server) baseAuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	if !slices.Contains([]string{"user_token", "admin_token"}, token) {
		_ = c.AbortWithError(http.StatusUnauthorized, ErrNotAuthorized)

		return
	}

	c.Next()
}

// baseAuthMiddleware set up middleware for checking if user has admin auth token
func (s *Server) adminAuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")

	switch token {
	case "user_token":
		_ = c.AbortWithError(http.StatusForbidden, ErrForbidden)
	case "admin_token":
		c.Next()
	default:
		_ = c.AbortWithError(http.StatusUnauthorized, ErrNotAuthorized)
	}
}

// timeoutMiddleware sets the middleware to control response timeout
func (s *Server) timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(s.timeout),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	_ = c.AbortWithError(http.StatusInternalServerError, ErrTimeout)

	return
}
