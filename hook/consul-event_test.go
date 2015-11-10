package hook

import "testing"

func TestConsulEventHook(t *testing.T) {
	hook := ConsulEventHook{}

	err := hook.Execute(map[string]string{
		"name":    "graphite",
		"tag":     "blue",
		"service": "data-service-api-v4",
		"node":    "api2.data-service-v4.iad",
	})

	if err != nil {
		t.Error(err)
	}
}
