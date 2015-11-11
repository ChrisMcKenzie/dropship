package hook

import (
	"os"
	"testing"

	"github.com/ChrisMcKenzie/dropship/service"
)

func TestConsulEventHook(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}
	hook := ConsulEventHook{}

	err := hook.Execute(map[string]interface{}{
		"name":    "graphite",
		"tag":     "blue",
		"service": "data-service-api-v4",
		"node":    "api2.data-service-v4.iad",
	}, service.Config{})

	if err != nil {
		t.Error(err)
	}
}
