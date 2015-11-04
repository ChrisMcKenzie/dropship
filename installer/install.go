package installer

import "io"

// Installer is an interface that allows different methods of writing
// the given io.Reader to disk.
type Installer interface {
	Install(dest string, r io.Reader) (int, error)
}
