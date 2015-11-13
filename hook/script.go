package hook

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/ChrisMcKenzie/dropship/service"
)

type ScriptHook struct{}

func (h ScriptHook) Execute(config map[string]interface{}, service service.Config) error {
	if c, ok := config["command"].(string); ok {

		// TODO(ChrisMcKenzie): Make this more secure by jailing it.
		var cwd string
		if len(service.Artifact) >= 1 {
			cwd = service.Artifact[0].Destination
		}

		_, err := executeCommand(c, cwd)
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
