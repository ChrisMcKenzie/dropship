// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
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
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/hcl"
)

type Artifact map[string]string
type HookConfig map[string]string
type HookDefinition map[string]HookConfig

type Config struct {
	Name          string              `hcl:",key"`
	CheckInterval string              `hcl:"checkInterval"`
	UpdateTTL     string              `hcl:"updateTTL"`
	PostCommand   string              `hcl:"postCommand"`
	PreCommand    string              `hcl:"preCommand"`
	BeforeHooks   []HookDefinition    `hcl:"before"`
	AfterHooks    []HookDefinition    `hcl:"after"`
	Sequential    bool                `hcl:"sequentialUpdates"`
	RawArtifact   map[string]Artifact `hcl:"artifact,expand"`
	Artifact      Artifact            `hcl:"-"`
	Hash          string              `hcl:"hash"`
	Updater       Updater             `hcl:"-"`
	Locker        Locker              `hcl:"-"`
}

type ServiceFile struct {
	Services []Config `hcl:"service,expand"`
}

func LoadServices(root string) (d []Config, err error) {
	files, _ := filepath.Glob(root + "/*.hcl")
	for _, file := range files {
		data, err := readFile(file)
		if err != nil {
			return nil, err
		}

		var deploy ServiceFile
		err = hcl.Decode(&deploy, data)
		if err != nil {
			return nil, err
		}

		for i, service := range deploy.Services {
			for key, value := range service.RawArtifact {
				value["type"] = key
				deploy.Services[i].Artifact = value
			}
		}

		d = append(d, deploy.Services...)
	}
	return
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
