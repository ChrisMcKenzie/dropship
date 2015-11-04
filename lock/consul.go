package lock

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/consul/api"
)

type ConsulLocker struct {
	semaphore *api.Semaphore
}

func NewConsulLocker(prefix string, config *api.Config) (*ConsulLocker, error) {
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	name, _ := os.Hostname()
	s, err := client.SemaphoreOpts(&api.SemaphoreOptions{
		Prefix: filepath.Join("dropship/services/", prefix),
		Limit:  1,

		SessionName: name,
	})
	if err != nil {
		return nil, err
	}

	l := &ConsulLocker{s}

	return l, nil
}

func (l *ConsulLocker) Acquire(shutdownCh <-chan struct{}) (<-chan struct{}, error) {
	return l.semaphore.Acquire(shutdownCh)
}

func (l *ConsulLocker) Release() error {
	return l.semaphore.Release()
}
