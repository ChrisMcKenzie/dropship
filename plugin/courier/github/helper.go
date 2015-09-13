package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type DeploymentHook struct {
	Deployment github.Deployment `json:"deployment"`
	Repository github.Repository `json:"repository"`
	Sender     github.User       `json:"sender"`
}

func GetPayload(req *http.Request) []byte {
	var payload = req.FormValue("payload")
	if len(payload) == 0 {
		raw, _ := ioutil.ReadAll(req.Body)
		return raw
	}
	return []byte(payload)
}

func ParseHook(payload []byte) (hook DeploymentHook, err error) {
	err = json.Unmarshal(payload, &hook)
	return
}

func GetClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}
