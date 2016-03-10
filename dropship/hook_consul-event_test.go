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
	"os"
	"testing"

	"github.com/hashicorp/consul/api"
)

func TestConsulEventHook(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}
	hook := ConsulEventHook{api.DefaultConfig()}

	err := hook.Execute(HookConfig{
		"name":    "graphite",
		"tag":     "blue",
		"service": "data-service-api-v4",
		"node":    "api2.data-service-v4.iad",
	}, Config{})

	if err != nil {
		t.Error(err)
	}
}
