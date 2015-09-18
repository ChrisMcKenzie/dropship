package github

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/ChrisMcKenzie/dropship/util"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	ghauth "golang.org/x/oauth2/github"
)

const (
	DefaultAPI = "https://api.github.com/"
	DefaultURL = "https://github.com"
	Kind       = "github.com"
)

type GitHub struct {
	Client string
	Secret string
	URL    string
	API    string
}

func New(client, secret string) *GitHub {
	return &GitHub{client, secret, DefaultURL, DefaultAPI}
}

func (g *GitHub) GetKind() string {
	return Kind
}

func (g *GitHub) Authorize(c *gin.Context) (*model.Authentication, error) {
	fmt.Println(fmt.Sprintf("http://%s/auth/%s", c.Request.Host, g.GetKind()))
	config := &oauth2.Config{
		ClientID:     g.Client,
		ClientSecret: g.Secret,
		Scopes:       []string{"user", "repo", "admin:repo_hook"},
		Endpoint:     ghauth.Endpoint,
		RedirectURL:  fmt.Sprintf("http://%s/auth/%s", c.Request.Host, g.GetKind()),
	}

	code := c.Query("code")
	state := c.Query("state")
	if len(code) == 0 {
		random := util.GenerateRandom()
		util.CreateCookie(c.Writer, "state", random)
		c.Redirect(
			http.StatusTemporaryRedirect,
			config.AuthCodeURL(random, oauth2.AccessTypeOnline),
		)
		return nil, nil
	}

	cookieState, err := util.GetCookieValue(c.Request, "state")
	util.DeleteCookie(c.Writer, "state")
	if cookieState != state || err != nil {
		return nil, errors.New("Invalid State Token")
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	client := GetClient(token.AccessToken)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, err
	}

	login := new(model.Authentication)
	login.Login = *user.Login
	login.Email = *user.Email
	login.Token = token.AccessToken
	login.Expiry = token.Expiry

	return login, nil
}

func (g *GitHub) GetRepos(user *model.User) ([]*model.Repo, error) {
	var repos []*model.Repo
	client := GetClient(user.Token)

	list, _, err := client.Repositories.List("", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		return repos, err
	}

	for _, repo := range list {
		repo := model.Repo{
			UserID:   user.Id,
			Name:     *repo.Name,
			Owner:    *repo.Owner.Login,
			Courier:  "github.com",
			URL:      *repo.HTMLURL,
			CloneURL: *repo.GitURL,
			Private:  *repo.Private,
		}

		repos = append(repos, &repo)
	}

	return repos, nil
}

func (g *GitHub) GetScript(user *model.User, repo *model.Repo, hook *model.Hook) ([]byte, error) {
	client := GetClient(user.Token)
	var opts = new(github.RepositoryContentGetOptions)
	opts.Ref = hook.Sha
	content, _, _, err := client.Repositories.GetContents(repo.Owner, repo.Name, "dropship.yml", opts)
	if err != nil {
		return nil, err
	}
	return content.Decode()
}

func (g *GitHub) ParseHook(r *http.Request) (*model.Deployment, error) {

	if r.Header.Get("X-Github-Event") == "ping" {
		return nil, nil
	}

	payload := GetPayload(r)
	data, err := ParseHook(payload)
	if err != nil {
		return nil, err
	}

	deploy := new(model.Deployment)
	deploy.Id = *data.Deployment.ID
	deploy.Owner = *data.Repository.Owner.Login
	deploy.Repo = *data.Repository.Name
	deploy.Sha = *data.Deployment.SHA
	deploy.Environment = *data.Deployment.Environment
	deploy.Author = *data.Sender.Name
	deploy.Gravatar = *data.Sender.GravatarID
	deploy.Timestamp = time.Now().UTC()
	deploy.Message = *data.Deployment.Description

	if len(deploy.Owner) == 0 {
		deploy.Owner = *data.Repository.Owner.Name
	}

	return deploy, nil
}
