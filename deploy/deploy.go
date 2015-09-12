package deploy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/couriers"
	"github.com/ChrisMcKenzie/dropship/couriers/github"
	"github.com/ChrisMcKenzie/dropship/servers"
	"github.com/ChrisMcKenzie/dropship/ssh"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

func HandleDeploy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	handler := p.ByName("provider")

	var deployment *couriers.Deployment
	var err error
	var c couriers.Courier
	if handler == "github.com" {
		log.Info("handling deploy with github courier...")
		c = github.NewGithubCourier()
		deployment, err = c.Handle(r)
		if err != nil {
			log.Errorf("[GITHUB] %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if deployment == nil {
			// encountered a ping event just pass ok
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	s := []servers.Server{}
	if server, ok := deployment.Servers[deployment.Environment]; ok {
		if server.Provider == "consul" {
			s, err = servers.GetServersFromConsul(server.Options)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if server.Provider == "list" {
			s, err = servers.GetServersFromPayload(server.Options)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		err = errors.New(fmt.Sprintf("servers for environment %s are not defined", deployment.Environment))
		log.Error(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	run(c, *deployment, s)

	fmt.Fprintf(w, "deploying %s", p.ByName("repo_name"))
}

func run(c couriers.Courier, d couriers.Deployment, servers []servers.Server) {
	log.Infof("Deploying to %s", d.Environment)
	for _, server := range servers {
		go func() {
			log.Infof("Excecuting Command on %s", server.Address)
			res := ssh.Execute(
				d.Commands,
				server.Address,
				ssh.NewClientConfig(server.User, server.Password),
			)

			err := c.UpdateStatus(d, "success", fmt.Sprintf("Deployed to %s", server.Address))
			if err != nil {
				log.Error(err)
			}
			log.Info(res)
		}()
	}
}
