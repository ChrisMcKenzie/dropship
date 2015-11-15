package dropship

import "io"

type (
	MetaData struct {
		ContentType string
		Hash        string
	}

	// Updater is an interface that defines methods for checking a files
	// freshness and downloading an updated version if needed.
	Updater interface {
		IsOutdated(hash string, opts Artifact) (bool, error)
		Download(Artifact) (io.ReadCloser, MetaData, error)
	}
)
