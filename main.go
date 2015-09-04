package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/ssh"
	"github.com/google/go-github/github"
	"github.com/hashicorp/consul/api"
	"github.com/julienschmidt/httprouter"
	"github.com/the-control-group/data-service-api/logger"
	"github.com/thoas/stats"
)

type (
	Payload struct {
		Deployment github.Deployment `json:"deployment"`
		Repository github.Repository `json:"repository"`
		Sender     github.User       `json:"sender"`
	}

	Server struct {
		Address  string `json:"address"`
		User     string `json:"username"`
		Password string `json:"password"`
	}

	Deployment struct {
		Command string `json:"command"`
		Servers struct {
			Provider string          `json:"provider"`
			Options  json.RawMessage `json:"options"`
		} `json:"servers"`
	}
)

var log = logger.NewLogger()
var port = "3000"

func Logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infof("[%s] %s", r.Method, r.URL)
		h(w, r, p)
	}
}

func handleGithubDeploy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	headers := r.Header

	if headers.Get("X-GitHub-Event") != "deploy" {
		log.Infof("Unable to handle event type %s", headers.Get("X-GitHub-Event"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	payload := Payload{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deployment := Deployment{}
	if err := json.Unmarshal(payload.Deployment.Payload, &deployment); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	servers := []Server{}
	if deployment.Servers.Provider == "consul" {
		getServersFromConsul(deployment.Servers.Options)
	} else if deployment.Servers.Provider == "payload" {
		servers, err = getServersFromPayload(deployment.Servers.Options)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	run(deployment.Command, servers)

	fmt.Fprintf(w, "deploying %s", payload.Repository.Name)
}

func run(c string, servers []Server) {
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

func getServersFromConsul(opt json.RawMessage) (servers []Server, err error) {
	options := struct {
		serviceName string `json:"service_name"`
	}{}
	err = json.Unmarshal(opt, &options)
	if err != nil {
		return
	}

	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return
	}

	catalog := client.Catalog()
	services, _, err := catalog.Service(options.serviceName, "", nil)
	if err != nil {
		return
	}

	for _, service := range services {
		servers = append(
			servers,
			Server{
				Address: service.Address,
			},
		)
	}

	return
}

func getServersFromPayload(opt json.RawMessage) (servers []Server, err error) {
	list := struct {
		List []Server `json:"list"`
	}{}
	if err = json.Unmarshal(opt, &list); err != nil {
		return
	}

	servers = list.List

	return
}

func main() {
	router := httprouter.New()
	s := stats.New()

	router.POST("/deploy/github.com/:repo_owner/:repo_name",
		Logger(handleGithubDeploy))

	log.Info("Dropship listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Handler(router)))
}
