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
