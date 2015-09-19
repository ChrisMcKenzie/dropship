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

	db.Create(&repo)

	return nil
}

func GetRepo(owner, name string) *model.Repo {
	repo := &model.Repo{
		Owner: owner,
		Name:  name,
	}
	db.Preload("User").Where(repo).First(repo)
	return repo
}

func GetRepos() (repos []model.Repo) {
	db.Find(&repos)
	for i, _ := range repos {
		db.Model(repos[i]).Related(&repos[i].User)
	}
	return repos
}
