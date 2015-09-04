package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/couriers"
	"github.com/ChrisMcKenzie/dropship/servers"
	"github.com/ChrisMcKenzie/dropship/ssh"
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/thoas/stats"
)

var log = logrus.New()
var port = "3000"

func Logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infof("[%s] %s", r.Method, r.URL)
		h(w, r, p)
	}
}

func run(c string, servers []servers.Server) {
	for _, server := range servers {
		go func() {
			log.Debugf("Excecuting Command on %s", server.Address)
			res := ssh.Execute(
				c,
				server.Address,
				ssh.NewClientConfig(server.User, server.Password),
			)
			log.Debug(res)
		}()
	}
}

func main() {
	router := httprouter.New()
	s := stats.New()

	router.GET("/_service/stats", Logger(func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s, err := json.Marshal(s.Data())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(s)
	}))

	router.POST("/deploy/:provider/:repo_owner/:repo_name",
		Logger(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			handler := p.ByName("provider")

			var deployment couriers.Deployment
			var err error
			if handler == "github.com" {
				c := couriers.NewGithubCourier()
				deployment, err = c.Handle(r)
				if err != nil {
					log.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			s := []servers.Server{}
			if deployment.Servers.Provider == "consul" {
				s, err = servers.GetServersFromConsul(deployment.Servers.Options)
				if err != nil {
					log.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else if deployment.Servers.Provider == "payload" {
				s, err = servers.GetServersFromPayload(deployment.Servers.Options)
				if err != nil {
					log.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			run(deployment.Command, s)

			fmt.Fprintf(w, "deploying %s", p.ByName("repo_name"))
		}))

	log.Info("Dropship listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Handler(router)))
}
