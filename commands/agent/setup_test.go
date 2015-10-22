package agent

import (
	"bytes"
	"io"
	"testing"

	"github.com/ChrisMcKenzie/dropship/repo"
	"github.com/ChrisMcKenzie/dropship/structs"
)

var (
	mockCh      = make(chan struct{})
	mockService = structs.Service{
		CheckInterval: "",
		Artifact: structs.ArtifactConfig{
			Repo: "mock",
		},
	}
)

type MockRepo struct{}

func (r *MockRepo) GetName() string {
	return "mock"
}

func (r *MockRepo) IsUpdated(s structs.Service) (bool, error) {
	return true, nil
}

func (r *MockRepo) Download(s structs.Service) (io.Reader, repo.MetaData, error) {
	return new(bytes.Buffer), repo.MetaData{}, nil
}

func TestSetup(t *testing.T) {
	repo.Register(&MockRepo{})

	updater, err := setup(mockService, mockCh)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if updater.service.CheckInterval != "10s" {
		t.Error("CheckInterval was not properly defaulted.")
		t.Fail()
	}

	mockService.CheckInterval = "invalid time"
	_, intErr := setup(mockService, mockCh)
	if intErr == nil {
		t.Error("expected setup to fail due to bad interval")
		t.Fail()
	}

	mockService.CheckInterval = "1s"
	mockService.Artifact.Repo = "Unknown Repo"
	_, repoErr := setup(mockService, mockCh)
	if repoErr == nil {
		t.Error("expected setup to fail due to bad repo")
		t.Fail()
	}
}
