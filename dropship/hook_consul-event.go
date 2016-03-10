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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ConsulEventHook struct {
	config *api.Config
}

func NewConsulEventHook(cfg map[string]string) ConsulEventHook {
	config := initializeConsulConfig(cfg)
	return ConsulEventHook{config}
}

func (h ConsulEventHook) Execute(config HookConfig, service Config) error {
	client, err := api.NewClient(h.config)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"hash": service.Hash,
	}

	plBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	name, ok := config["name"]
	serv, ok := config["service"]
	tag, ok := config["tag"]
	node, ok := config["node"]

	if !ok {
		return errors.New("Consul Hook: invalid config")
	}

	id, meta, err := client.Event().Fire(&api.UserEvent{
		Name:          name,
		Payload:       plBytes,
		ServiceFilter: serv,
		TagFilter:     tag,
		NodeFilter:    node,
	}, nil)

	fmt.Println(id, meta, err)

	return err
}
