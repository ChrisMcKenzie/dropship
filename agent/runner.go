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
package agent

import "sync"

// Worker is a data type that can perform work.
type Worker interface {
	Work()
}

// Runner is a type of worker pool that takes its work over a channel
// this allows for a dispatcher to actually signal the work.
type Runner struct {
	work chan Worker
	wg   sync.WaitGroup
}

// NewRunner creates a new runner with the number go routines created and
// ready for work
func NewRunner(maxGoRoutines int) *Runner {
	r := Runner{
		work: make(chan Worker),
	}

	r.wg.Add(maxGoRoutines)
	for i := 0; i < maxGoRoutines; i++ {
		go func() {
			for w := range r.work {
				w.Work()
			}
			r.wg.Done()
		}()
	}

	return &r
}

// Do will schedule work on to the runner pull by placing a worker stuct
// on the channel
func (r *Runner) Do(w Worker) {
	r.work <- w
}

// Shutdown will signal the pool to stop accepting work and finish any
// current jobs
func (r *Runner) Shutdown() {
	close(r.work)
	r.wg.Wait()
}
