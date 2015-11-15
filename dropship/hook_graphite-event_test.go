package dropship

import (
	"os"
	"testing"
)

func TestGraphiteEventHook(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}

	var hook GraphiteEventHook

	err := hook.Execute(HookConfig{
		"host": "http://graphite2.analytics.iad",
		"what": "deployed by dropship",
		"tags": "data-service deployment",
		"data": "dropship is awesome!",
	}, Config{})

	if err != nil {
		t.Error(err)
	}
}
