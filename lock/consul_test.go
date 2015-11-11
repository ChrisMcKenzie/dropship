package lock

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

var (
	locker Locker
	err    error
)

func TestMain(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}
	locker, err = NewConsulLocker("dropship", api.DefaultConfig())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestConsulLockerAcquire(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}
	lock, err := locker.Acquire(nil)
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

	err = locker.Release()
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
