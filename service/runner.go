package service

import "sync"

type Worker interface {
	Work()
}

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
