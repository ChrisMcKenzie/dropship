package main

import (
	"encoding/json"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/dropship"
	"github.com/ChrisMcKenzie/dropship/dropship/auth"
	"github.com/ChrisMcKenzie/dropship/dropship/github"
	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/thoas/stats"
)

var log = logging.GetLogger()
var port = "3000"

func Logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infof("[%s] %s", r.Method, r.URL)
		h(w, r, p)
	}
}

func main() {
	router := httprouter.New()
	s := stats.New()

	ga := auth.NewGithubAuth()

	router.GET("/",
		Logger(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			http.FileServer(http.Dir("./ui")).ServeHTTP(w, r)
		}))

	router.GET("/api/github/repos", Logger(github.GetRepos))
	router.POST("/api/github/repos/:repo_owner/:repo_name/hook", Logger(github.AddHook))
	router.GET("/auth/github", Logger(ga.AuthHandle))
	router.GET("/auth/github/callback", Logger(ga.CallbackHandle))

	router.GET("/_service/stats",
		Logger(func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			s, err := json.Marshal(s.Data())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Write(s)
		}))

	router.POST("/deploy/:provider/:repo_owner/:repo_name",
		Logger(dropship.HandleDeploy))

	log.Info("Dropship listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Handler(router)))
}
