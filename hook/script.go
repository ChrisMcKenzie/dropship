package hook

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/ChrisMcKenzie/dropship/service"
)

type ScriptHook struct{}

func (h ScriptHook) Execute(config map[string]interface{}, service service.Config) error {
	if c, ok := config["command"].(string); ok {
		_, err := executeCommand(c)
		return err
	}
	return errors.New("Script: exiting no command was given")
}

func executeCommand(c string) (string, error) {
	cmd := strings.Fields(c)
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}
