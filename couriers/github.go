package couriers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ChrisMcKenzie/dropship/dropship/database"
	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/google/go-github/github"
	"github.com/libgit2/git2go"
	"golang.org/x/oauth2"
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
	log.Info("Cloning repo...")
	_, err := cloneRepo(payload)
	if err != nil {
		return d, err
	}

	// Read dropship.yml
	bytes, err := ioutil.ReadFile(
		fmt.Sprintf(
			pathTemplate+"/dropship.yml",
			*payload.Repository.Owner.Login,
			*payload.Repository.Name,
		),
	)
	d, err = parseDeployment(bytes)
	if err != nil {
		return d, err
	}

	os.RemoveAll(fmt.Sprintf(pathTemplate, *payload.Repository.Owner.Login, *payload.Repository.Name))

	d.Id = *payload.Deployment.ID
	d.Owner = *payload.Repository.Owner.Login
	d.Repo = *payload.Repository.Name

	return d, nil
}

func (g *GitHubCourier) UpdateStatus(deployment Deployment, status string, desc string) error {
	token, err := database.GetTokenFor(
		fmt.Sprintf("%s/%s", deployment.Owner, deployment.Repo),
	)

	if err != nil {
		return err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	log.Debug(desc)
	_, _, err = client.Repositories.CreateDeploymentStatus(
		deployment.Owner,
		deployment.Repo,
		deployment.Id,
		&github.DeploymentStatusRequest{
			State:       &status,
			Description: &desc,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	log.Debugf("looking for keys in %s", os.Getenv("KEY_PATH")+"/.ssh/id_rsa.pub")
	ret, cred := git.NewCredSshKey(
		"git",
		os.Getenv("KEY_PATH")+"/.ssh/id_rsa.pub",
		os.Getenv("KEY_PATH")+"/.ssh/id_rsa",
		"",
	)
	return git.ErrorCode(ret), &cred
}

// Made this one just return 0 during troubleshooting...
func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

func cloneRepo(payload Payload) (*git.Repository, error) {
	cloneOptions := &git.CloneOptions{}
	// use FetchOptions instead of directly RemoteCallbacks
	// https://github.com/libgit2/git2go/commit/36e0a256fe79f87447bb730fda53e5cbc90eb47c
	cloneOptions.FetchOptions = &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	repo, err := git.Clone(
		fmt.Sprintf(
			"git@github.com:%s/%s.git",
			*payload.Repository.Owner.Login,
			*payload.Repository.Name,
		),
		fmt.Sprintf(
			"/tmp/dropship/%s/%s",
			*payload.Repository.Owner.Login,
			*payload.Repository.Name,
		),
		cloneOptions,
	)

	if err != nil {
		return nil, err
	}

	return repo, nil
}
