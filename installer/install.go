package installer

import "io"

type Installer interface {
	// Install Defines a Method that takes a destination path
	// and a io.Reader and untars and gzip decodes a tarball and
	// places the files inside on the FS with `dest` as their root
	// It returns the number of files written and an error
	Install(dest string, r io.Reader) (int, error)
}
