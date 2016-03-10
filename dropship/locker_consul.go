// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package dropship

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/consul/api"
)

// ConsulLocker is a Locker that will use consul as the coordinator for
// establish a lock amongst multiple machines
type ConsulLocker struct {
	semaphore *api.Semaphore
}

func NewConsulLocker(cfg map[string]string) (*ConsulLocker, error) {
	config := initializeConsulConfig(cfg)

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	name, _ := os.Hostname()

	var ttl = api.DefaultSemaphoreSessionTTL
	if cfg["ttl"] != "" {
		ttl = cfg["ttl"]
	}

	s, err := client.SemaphoreOpts(&api.SemaphoreOptions{
		Prefix: filepath.Join("dropship/services/", cfg["prefix"]),
		Limit:  1,

		SessionTTL: ttl,

		SessionName: name,
	})
	if err != nil {
		return nil, err
	}

	l := &ConsulLocker{s}

	return l, nil
}

func (l ConsulLocker) Acquire(shutdownCh <-chan struct{}) (<-chan struct{}, error) {
	return l.semaphore.Acquire(shutdownCh)
}

func (l ConsulLocker) Release() error {
	return l.semaphore.Release()
}

func initializeConsulConfig(cfg map[string]string) *api.Config {
	config := api.DefaultConfig()

	if addr, ok := cfg["address"]; ok {
		config.Address = addr
	}

	if token, ok := cfg["token"]; ok {
		config.Token = token
	}

	if user, ok := cfg["user"]; ok {
		var password string
		if pass, ok := cfg["password"]; ok {
			password = pass
		}

		config.HttpAuth = &api.HttpBasicAuth{
			Username: user,
			Password: password,
		}
	}

	if ssl, ok := cfg["useSSL"]; ok {
		enabled, err := strconv.ParseBool(ssl)
		if err != nil {
			log.Printf("[ERR]: Could not parse consul useSSL: %s", err)
		}

		if enabled {
			config.Scheme = "https"
		}
	}

	return config
}
