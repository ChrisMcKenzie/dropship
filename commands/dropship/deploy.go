package dropship

import (
	"log"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/database"
	"github.com/ChrisMcKenzie/dropship/plugin/courier"
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) DeployHook(c *gin.Context) {
	user := ToUser(c)
	courierParam := c.Param("courier")
	courier := courier.Lookup(courierParam)
	if courier == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	deployment, err := courier.ParseHook(c.Request)
	if err != nil {
		log.Printf("Unable to parse hook. %s\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	repo := database.GetRepo(deployment.Owner, deployment.Repo)
	if repo == nil {
		log.Printf("Unable to retrieve Repo.")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if repo.Active == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	script, err := courier.GetScript(user, repo, deployment)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	_, err = parseScript(script)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	// add task to work pool
}
