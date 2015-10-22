package repo

import (
	"io"

	"github.com/ChrisMcKenzie/dropship/structs"
)

type (
	Repo interface {
		GetName() string
		IsUpdated(structs.Service) (bool, error)
		Download(structs.Service) (io.Reader, MetaData, error)
	}

	MetaData struct {
		Hash        string
		ContentType string
	}
)

var repos []Repo

func Register(r Repo) {
	repos = append(repos, r)
}

func GetRepo(name string) Repo {
	for _, repo := range repos {
		if repo.GetName() == name {
			return repo
		}
	}

	return nil
}
