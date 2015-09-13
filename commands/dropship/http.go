package dropship

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	_ "github.com/chrismckenzie/dropship/plugin/courier/github"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router *gin.Engine
	logger *logrus.Logger
}

func NewHTTPServer(addr string) {
	server := &HTTPServer{
		router: gin.Default(),
		logger: logrus.New(),
	}

	server.registerHandlers()

	http.ListenAndServe(addr, server.router)
}

func (s *HTTPServer) registerHandlers() {
	s.router.StaticFile("/", "./ui/index.html")
	s.router.POST("/deploy/:courier", s.DeployHook)

	s.router.GET("/auth/:courier", s.Auth)

	api := s.router.Group("/api", s.AuthMiddleware())
	{
		api.POST("/repos", s.PostRepo)
		api.GET("/repos", s.GetRepos)
	}
}
