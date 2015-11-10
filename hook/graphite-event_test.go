package hook

import "testing"

func TestGraphiteEventHook(t *testing.T) {
	var hook GraphiteEventHook

	err := hook.Execute(map[string]string{
		"host": "http://graphite2.analytics.iad",
		"what": "deployed by dropship",
		"tags": "data-service deployment",
		"data": "dropship is awesome!",
	})

	if err != nil {
		t.Error(err)
	}
}
