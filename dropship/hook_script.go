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
	"errors"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// ScriptHook defines a Hook that will run a specified command on a the
// machine.
type ScriptHook struct{}

// Execute is a Hook method to execut the hook with standard options
func (h ScriptHook) Execute(config HookConfig, service Config) error {
	if c := config["command"]; c != "" {

		// TODO(ChrisMcKenzie): Make this more secure by jailing it.
		var cwd string
		if len(service.Artifact) >= 1 {
			cwd = service.Artifact["destination"]
		}

		out, err := executeCommand(c, cwd)
		log.Infof("%s", out)
		return err
	}
	return errors.New("Script: exiting no command was given")
}

func executeCommand(c, cwd string) (string, error) {
	sCmd := strings.Fields(c)
	cmd := exec.Command(sCmd[0], sCmd[1:]...)
	if _, err := os.Stat(cwd); !os.IsNotExist(err) {
		cmd.Dir = cwd
	}
	out, err := cmd.Output()
	return string(out), err
}
