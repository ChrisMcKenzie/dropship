package dropship

import (
	"io"
	"os"
)

// FileInstaller defines an Installer that takes the reader and writes
// it to the dest directory.
type FileInstaller struct{}

func (i FileInstaller) Install(dest string, f io.Reader) (count int, err error) {
	moveOld(dest)
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
