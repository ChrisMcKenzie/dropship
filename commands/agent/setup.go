package agent

import (
	"errors"
	"time"

	"github.com/ChrisMcKenzie/dropship/repo"
	"github.com/ChrisMcKenzie/dropship/structs"
)

func setup(service structs.Service, shutdownCh <-chan struct{}) (*updater, error) {
	if service.CheckInterval == "" {
		service.CheckInterval = "10s"
	}

	tickerDur, err := time.ParseDuration(service.CheckInterval)
	if err != nil {
		return nil, err
	}

	repo := repo.GetRepo(service.Artifact.Repo)
	if repo == nil {
		return nil, errors.New("Unable to find repo " + service.Artifact.Repo)
	}

	updater := &updater{
		time.NewTicker(tickerDur),
		shutdownCh,
		service,
		repo,
	}

	return updater, nil
}
