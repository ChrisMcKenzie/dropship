package work

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

func (r *Runner) Do(w Worker) {
	r.work <- w
}

func (r *Runner) Shutdown() {
	close(r.work)
	r.wg.Wait()
}
