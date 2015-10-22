package rackspace

import (
	"io"
	"log"

	"github.com/ChrisMcKenzie/dropship/repo"
	"github.com/ChrisMcKenzie/dropship/structs"
	"github.com/ncw/swift"
)

const (
	RepoName = "rackspace"
	AUTH_URL = "https://identity.api.rackspacecloud.com/v2.0"
)

type RackspaceRepo struct {
	connection *swift.Connection
}

func Setup(user, key, region string) {
	if user == "" || key == "" {
		return
	}

	rackConnection := &swift.Connection{
		// This should be your username
		UserName: user,
		// This should be your api key
		ApiKey: key,
		// This should be a v1 auth url, eg
		//  Rackspace US        https://auth.api.rackspacecloud.com/v1.0
		//  Rackspace UK        https://lon.auth.api.rackspacecloud.com/v1.0
		//  Memset Memstore UK  https://auth.storage.memset.com/v1.0
		AuthUrl: AUTH_URL,
		// Region to use - default is use first region if unset
		Region: region,
		// Name of the tenant - this is likely your username
	}

	repo.Register(&RackspaceRepo{rackConnection})
}

func (repo *RackspaceRepo) GetName() string {
	return RepoName
}

func (repo *RackspaceRepo) IsUpdated(s structs.Service) (bool, error) {
	if !repo.connection.Authenticated() {
		err := repo.connection.Authenticate()
		if err != nil {
			log.Fatal(err)
		}
	}

	info, _, err := repo.connection.Object(
		s.Artifact.Bucket,
		s.Artifact.Path,
	)

	if err != nil {
		return true, err
	}

	if info.Hash == s.Hash {
		return true, nil
	}

	return false, nil
}

func (r *RackspaceRepo) Download(s structs.Service) (io.Reader, repo.MetaData, error) {
	log.Println("Downloading", s.Artifact.Path, "from", s.Artifact.Bucket)
	file, hdrs, err := r.connection.ObjectOpen(
		s.Artifact.Bucket,
		s.Artifact.Path,
		true,
		swift.Headers{},
	)
	return file, repo.MetaData{hdrs["Etag"], hdrs["Content-Type"]}, err
}
