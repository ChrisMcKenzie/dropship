package agent

import (
	"path/filepath"

	"github.com/ChrisMcKenzie/dropship/structs"
	"github.com/hashicorp/consul/api"
)

const BasePrefix = "dropship/locks"

func AcquireLock(s structs.Service) (*api.Lock, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return client.LockKey(filepath.Join(BasePrefix, s.Name))
}
