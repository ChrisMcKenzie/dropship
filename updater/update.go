package updater

import "io"

type Updater interface {
	IsOutdated(hash string, opts *Options) (bool, error)
	Download(*Options) (io.ReadCloser, MetaData, error)
}
