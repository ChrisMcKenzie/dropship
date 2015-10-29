package updater

import "io"

type Updater interface {
	CheckForUpdate(hash string, opts *Options) (bool, error)
	Update(*Options) (io.Reader, MetaData, error)
}
