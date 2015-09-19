package dropship

import (
	"github.com/ChrisMcKenzie/dropship/database"
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) GetUsers(c *gin.Context) {
	c.JSON(200, database.GetUsers())
}
