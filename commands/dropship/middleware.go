package dropship

import "github.com/gin-gonic/gin"

func (s *HTTPServer) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if authenticated
		c.Next()
	}
}
