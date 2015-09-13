package dropship

import (
	"net/http"

	"github.com/ChrisMcKenzie/dropship/plugin/courier"
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) Auth(c *gin.Context) {
	courierParam := c.Param("courier")
	courier := courier.Lookup(courierParam)
	if courier == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	login, err := courier.Authorize(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else if login == nil {
		return
	}
}
