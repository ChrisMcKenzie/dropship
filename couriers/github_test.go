package couriers

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/github"
)

func TestGithubCloneRepo(t *testing.T) {
	name := "dropship"
	owner := "ChrisMcKenzie"
	_, err := cloneRepo(Payload{
		Repository: github.Repository{
			Name:  &name,
			Owner: &github.User{Name: &owner},
		},
	})

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = os.Stat(fmt.Sprintf("/tmp/dropship/%s/%s", owner, name))
	if os.IsNotExist(err) {
		t.Log(err)
		t.Fail()
	}

	os.RemoveAll("/tmp/dropship")
}

func TestGithubParseDeployment(t *testing.T) {
	yaml := []byte(`
servers:
  provider: list
  options:
          list:
            - 127.0.0.1
command: ls
`)

	deployment, err := parseDeployment(yaml)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Logf("%v", deployment)
}
