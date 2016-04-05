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
package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Service struct {
	Name        string   `json:"name"`
	Href        string   `json:"href"`
	Description string   `json:"description"`
	LastDeploy  string   `json:"last_deployed_on"`
	Hosts       []string `json:"hosts"`
}

func servicesIndex(w http.ResponseWriter, r *http.Request) {
	keys, err := kvstore.List(fmt.Sprintf("%s/services/", DefaultKeyPrefix))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	var services []Service
	for _, val := range keys {
		key := strings.Split(val.Key, "/")
		var svc Service
		err := json.Unmarshal(val.Value, &svc)
		if err != nil {
			log.Error(err)
		}
		svc.Name = key[len(key)-1]
		services = append(services, svc)
	}

	payload, err := json.Marshal(services)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(payload)
}

func servicesShow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("hello, world"))
}
