package updater

import "io"

// Updater is an interface that defines methods for checking a files
// freshness and downloading an updated version if needed.
type Updater interface {
	IsOutdated(hash string, opts *Options) (bool, error)
	Download(*Options) (io.ReadCloser, MetaData, error)
}
