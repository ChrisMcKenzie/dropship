package dropship

import (
	"log"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/database"
	"github.com/ChrisMcKenzie/dropship/httputil"
	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/ChrisMcKenzie/dropship/plugin/courier"
	"github.com/ChrisMcKenzie/dropship/session"
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
		c.String(http.StatusBadRequest, err.Error())
		return
	} else if login == nil {
		return
	}

	// get or create user
	user := &model.User{
		Login:   login.Login,
		Email:   login.Email,
		Name:    login.Name,
		Courier: courier.GetKind(),
	}
	database.FindOrCreateUser(user)
	if user.Id == 0 {
		log.Println("Unable to create user account")
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	user.Token = login.Token
	user.TokenExpiry = login.Expiry

	if err := database.UpdateUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Print(user)
	// generate user session
	token, err := session.CreateToken(httputil.GetBaseURL(c.Request), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// redirect user to home
	c.Redirect(http.StatusSeeOther, "/#access_token="+token)
}
