package hook

import (
	"os/exec"
	"strings"

	"github.com/ChrisMcKenzie/dropship/service"
)

type ScriptHook struct{}

func (h ScriptHook) Execute(config map[string]interface{}, service service.Config) error {
	_, err := executeCommand(config["command"].(string))
	return err
}

func executeCommand(c string) (string, error) {
	cmd := strings.Fields(c)
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}
