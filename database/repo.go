package database

import (
	"errors"

	"github.com/ChrisMcKenzie/dropship/model"
)

func CreateRepo(repo *model.Repo) error {
	if repo.Name == "" {
		return errors.New("Name is a required field")
	}

	if repo.Owner == "" {
		return errors.New("Owner is a required field")
	}

	if repo.Courier == "" {
		return errors.New("Courier is a required field")
	}

	db.Create(repo)

	return nil
}

func GetRepo(owner, name string) *model.Repo {
	repo := &model.Repo{
		Owner: owner,
		Name:  name,
	}
	db.Where(repo).First(repo)
	return repo
}

func GetRepos() (repos []model.Repo) {
	db.Find(&repos)
	return repos
}
