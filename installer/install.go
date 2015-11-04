package installer

import (
	"io"
	"os"
	"path/filepath"
)

// Installer is an interface that allows different methods of writing
// the given io.Reader to disk.
type Installer interface {
	Install(dest string, r io.Reader) (int, error)
}

func moveOld(dest string) error {
	return os.Rename(dest, filepath.Join(dest, ".old"))
}

func cleanup(dest string, err error) error {
	oldPath := filepath.Join(dest, ".old")
	if err != nil {
		return os.Rename(oldPath, dest)
	}

	return os.RemoveAll(oldPath)
}
