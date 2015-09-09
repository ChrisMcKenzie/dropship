package couriers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/google/go-github/github"
	"github.com/libgit2/git2go"
	"gopkg.in/yaml.v2"
)

var log = logging.GetLogger()
var pathTemplate = "/tmp/dropship/%s/%s"

type (
	Payload struct {
		Deployment github.Deployment `json:"deployment"`
		Repository github.Repository `json:"repository"`
		Sender     github.User       `json:"sender"`
	}

	GitHubCourier struct {
		ApiKey string
	}
)

func NewGithubCourier() *GitHubCourier {
	return &GitHubCourier{}
}

func (c *GitHubCourier) Handle(r *http.Request) (Deployment, error) {
	var d Deployment
	headers := r.Header

	if headers.Get("X-GitHub-Event") != "deployment" {
		return d, errors.New("Unable to handle event " + headers.Get("X-GitHub-Event"))
	}

	payload := Payload{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		return d, err
	}

	if err := json.Unmarshal(payload.Deployment.Payload, &d); err != nil {
		return d, err
	}

	// Clone Repo
	_, err := cloneRepo(payload)
	if err != nil {
		return d, err
	}

	// Read dropship.yml
	bytes, err := ioutil.ReadFile(
		fmt.Sprintf(
			pathTemplate,
			*payload.Repository.Owner.Name,
			*payload.Repository.Name,
		),
	)
	d, err = parseDeployment(bytes)
	if err != nil {
		return d, err
	}

	return d, nil
}

func parseDeployment(file []byte) (Deployment, error) {
	var d Deployment
	err := yaml.Unmarshal(file, &d)
	if err != nil {
		return d, err
	}

	return d, nil
}

func cloneRepo(payload Payload) (*git.Repository, error) {
	repo, err := git.Clone(
		fmt.Sprintf(
			"git://git@github.com/%s/%s.git",
			*payload.Repository.Owner.Name,
			*payload.Repository.Name,
		),
		fmt.Sprintf(
			"/tmp/dropship/%s/%s",
			*payload.Repository.Owner.Name,
			*payload.Repository.Name,
		),
		&git.CloneOptions{},
	)

	if err != nil {
		return nil, err
	}

	return repo, nil
}
