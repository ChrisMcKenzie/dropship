package lock

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

var (
	locker Locker
	err    error
)

func TestMain(t *testing.T) {
	locker, err = NewConsulLocker("dropship/services", api.DefaultConfig())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestConsulLockerAcquire(t *testing.T) {
	lock, err := locker.Acquire("test")
	if err != nil {
		t.Fatal(err)
		return
	}
	if lock == nil {
		t.Fatalf("Acquire: expected channel signal got: %v", lock)
		return
	}

	select {
	case <-lock:
		t.Fatal("Acquire: should be held")
	default:
	}

	err = locker.Release("test")
	if err != nil {
		t.Fatal(err)
	}

	// Should lose resource
	select {
	case <-lock:
	case <-time.After(time.Second):
		t.Fatalf("Acquire: should not be held")
	}
}
