package hook

import (
	"os"
	"testing"

	"github.com/ChrisMcKenzie/dropship/service"
)

func TestGraphiteEventHook(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}

	var hook GraphiteEventHook

	err := hook.Execute(map[string]interface{}{
		"host": "http://graphite2.analytics.iad",
		"what": "deployed by dropship",
		"tags": "data-service deployment",
		"data": "dropship is awesome!",
	}, service.Config{})

	if err != nil {
		t.Error(err)
	}
}
