package dropship

import (
	"fmt"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/couriers"
	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/ChrisMcKenzie/dropship/servers"
	"github.com/ChrisMcKenzie/dropship/ssh"
	"github.com/julienschmidt/httprouter"
)

var log = logging.GetLogger()

func HandleDeploy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	handler := p.ByName("provider")

	var deployment couriers.Deployment
	var err error
	var c couriers.Courier
	if handler == "github.com" {
		c = couriers.NewGithubCourier()
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
	} else if deployment.Servers.Provider == "list" {
		s, err = servers.GetServersFromPayload(deployment.Servers.Options)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	run(c, deployment, s)

	fmt.Fprintf(w, "deploying %s", p.ByName("repo_name"))
}

func run(c couriers.Courier, d couriers.Deployment, servers []servers.Server) {
	log.Infof("Deploying to %v", servers)
	for _, server := range servers {
		go func() {
			log.Debugf("Excecuting Command on %s", server.Address)
			res := ssh.Execute(
				d.Commands,
				server.Address,
				ssh.NewClientConfig(server.User, server.Password),
			)

			err := c.UpdateStatus(d, "success", fmt.Sprintf("Deployed to %s", server.Address))
			if err != nil {
				log.Error(err)
			}
			log.Debug(res)
		}()
	}
}
