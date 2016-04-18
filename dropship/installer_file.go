package dropship

import (
	"io"
	"os"
)

// FileInstaller defines an Installer that takes the reader and writes
// it to the dest directory.
type FileInstaller struct{}

// Install will copy the given io.Reader to the destination path
func (i FileInstaller) Install(dest string, f io.Reader) (count int, err error) {
	// if file exists lets move it so we can recover on failure
	if _, err := os.Stat(dest); err == nil {
		err = moveOld(dest)
		if err != nil {
			return 0, err
		}
		defer cleanup(dest, err)
	}

	if f == nil {
		return count, ErrNilReader
	}

	file, err := os.Create(dest)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, f)
	if err != nil {
		return
	}

	return 1, nil
}
