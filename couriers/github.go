package couriers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/go-github/github"
)

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

	return d, nil
}
