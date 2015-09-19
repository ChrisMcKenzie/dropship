package dropship

import (
	"net/http"
	"strings"

	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/ChrisMcKenzie/dropship/session"
	"github.com/gin-gonic/gin"
)

const BearerPrefix string = "Bearer "

func (s *HTTPServer) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.User
		auth := c.Request.Header.Get("Authorization")
		// fmt.Println(auth[len(BearerPrefix):])
		if strings.HasPrefix(auth, BearerPrefix) {
			user = session.GetUserFromJWT(auth[len(BearerPrefix):])
			if user != nil {
				c.Set("User", user)
			}
		}

		// fmt.Println(user)
		if user == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
