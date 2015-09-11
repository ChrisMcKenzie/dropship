package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ChrisMcKenzie/dropship/dropship/database"
	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

var log = logging.GetLogger()

func storeAccessToken() {

}

func AddHook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, _ := r.Cookie("github")
	log.Debug(c.Value)
	if c.Value == "" {
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Value},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	name := "web"

	url := fmt.Sprintf("%s/deploy/github.com/%s/%s", os.Getenv("APP_URL"), p.ByName("repo_owner"), p.ByName("repo_name"))

	hook, _, err := client.Repositories.CreateHook(
		p.ByName("repo_owner"),
		p.ByName("repo_name"),
		&github.Hook{
			Name:   &name,
			Events: []string{"deployment"},
			Config: map[string]interface{}{
				"url":          url,
				"content_type": "json",
			},
		},
	)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = database.StoreTokenFor(
		fmt.Sprintf("%s/%s", p.ByName("repo_owner"), p.ByName("repo_name")),
		c.Value,
	)

	if err != nil {
		log.Errorf("unable to store access token for repo %s", err)
	}

	log.Debugf("Hook Created: %v", hook)
	json, _ := json.Marshal(hook)
	w.Write(json)
}

func GetRepos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, _ := r.Cookie("github")
	log.Debug(c.Value)
	if c.Value == "" {
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Value},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List("", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		return
	}

	re, err := json.Marshal(repos)
	if err != nil {
		return
	}

	w.Write(re)
}
