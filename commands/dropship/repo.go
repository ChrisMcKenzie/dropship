package dropship

import (
	"net/http"

	"github.com/ChrisMcKenzie/dropship/database"
	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) PostRepo(c *gin.Context) {
	var repo model.Repo

	if c.Bind(&repo) == nil {
		err := database.CreateRepo(&repo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(201, repo)
}

func (s *HTTPServer) GetRepos(c *gin.Context) {
	c.JSON(200, database.GetRepos())
}
