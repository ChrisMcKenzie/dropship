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

func ToUser(c *gin.Context) *model.User {
	contextValue, _ := c.Get("User")
	user, ok := contextValue.(*model.User)
	if !ok {
		return nil
	}

	return user
}

func (s *HTTPServer) PostRepo(c *gin.Context) {
	var repo model.Repo
	user := ToUser(c)
	if user == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println(user)
	// Set repo fields
	repo.Name = c.Param("name")
	repo.Owner = c.Param("owner")
	repo.Courier = c.Param("courier")
	repo.Active = true
	repo.UserId = user.Id
	repo.User = *user

	// TODO(ChrisMcKenzie): Add Privatekey and secret to repo
	repo.Token = session.GenerateRandom()
	// Activate(Addhook) Repo on courier
	courier := courier.Lookup(c.Param("courier"))
	if courier == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	url := httputil.GetBaseURL(c.Request) + "/" + repo.Courier
	if err := courier.Activate(&repo, url); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add Repo to Datastore
	if err := database.CreateRepo(&repo); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, repo)
}

func (s *HTTPServer) GetCourierRepos(c *gin.Context) {
	user := ToUser(c)
	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	co := courier.Lookup(c.Param("courier"))
	if co == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	repos, err := co.GetRepos(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, repos)
}

func (s *HTTPServer) GetRepos(c *gin.Context) {
	c.JSON(200, database.GetRepos())
}
