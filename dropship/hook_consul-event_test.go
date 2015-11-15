package dropship

import (
	"os"
	"testing"
)

func TestConsulEventHook(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}
	hook := ConsulEventHook{}

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
