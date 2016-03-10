// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package dropship

import (
	"os"
	"testing"
	"time"
)

var (
	locker Locker
	err    error
)

func TestMain(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI")
	}
	locker, err = NewConsulLocker(map[string]string{
		"prefix": "dropship",
	})
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
