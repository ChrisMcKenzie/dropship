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

// TemplateData specifies some system configurations that can be used in hook
// options using go template syntax
type TemplateData struct {
	Config
	Hostname string
}

// Hook is an interface that defines a method for executing custom
// behavior at certain hook points
type Hook interface {
	Execute(config HookConfig, service Config) error
}

// GetHookByName returns a hook by the id name given.
//
// TODO(ChrisMcKenzie): this is kind of ugly and could probably be cleaned
// through the use of hooks being packages and registering with a hook manager
func GetHookByName(name string) Hook {
	switch name {
	case "script":
		return ScriptHook{}
	case "consul-event":
		return ConsulEventHook{}
	case "graphite-event":
		return GraphiteEventHook{}
	}

	return nil
}
