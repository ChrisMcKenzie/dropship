package lock

type Locker interface {
	Acquire() (<-chan struct{}, error)
	Release() error
}
