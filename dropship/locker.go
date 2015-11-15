package dropship

// Locker is an interface that allows you to block execution of
// another process across a set of machines.
type Locker interface {
	// Acquire takes a shutdownCh and return a lock chan and error
	//
	// the lock chan can be used to block the process until the lock
	// has been acquired and the chan receives.
	Acquire(<-chan struct{}) (<-chan struct{}, error)
	// Release will release the lock allowing for other processes to
	// acquire.
	Release() error
}
