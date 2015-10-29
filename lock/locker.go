package lock

type Locker interface {
	Acquire(<-chan struct{}) (<-chan struct{}, error)
	Release() error
}
