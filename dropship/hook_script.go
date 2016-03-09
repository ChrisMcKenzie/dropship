package dropship

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ScriptHook struct{}

func (h ScriptHook) Execute(config HookConfig, service Config) error {
	if c := config["command"]; c != "" {

		// TODO(ChrisMcKenzie): Make this more secure by jailing it.
		var cwd string
		if len(service.Artifact) >= 1 {
			cwd = service.Artifact["destination"]
		}

		out, err := executeCommand(c, cwd)
		log.Printf("[INFO]: %s", out)
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
