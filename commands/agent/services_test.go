package agent

import "testing"

func TestServices(t *testing.T) {
	services := loadServices("../..")
	if len(services) == 0 {
		t.Error("expected services to be loaded")
		t.Fail()
	}
}
