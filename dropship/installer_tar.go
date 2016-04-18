package dropship

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// ErrNilReader is a generic error for when the reader is nil
var ErrNilReader = errors.New("Install: must have a non-nil Reader")

// TarInstaller Defines an install Method that takes a destination path
// and a io.Reader and untars and gzip decodes a tarball and
// places the files inside on the FS with `dest` as their root
// It returns the number of files written and an error
type TarInstaller struct{}

// Install defines the Installer required function that will in this case
// untar and copy files to the destination directory
func (i TarInstaller) Install(dest string, fr io.Reader) (count int, err error) {
	defer cleanup(dest, err)
	moveOld(dest)
	if fr == nil {
		return count, ErrNilReader
	}

	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		var hdr *tar.Header
		hdr, err = tr.Next()
		if err == io.EOF {
			// end of tar archive
			err = nil
			return
		}
		if err != nil {
			return
		}

		if err = writePath(hdr, tr, dest); err != nil {
			return
		}
		count++
	}
}

func writePath(hdr *tar.Header, tr *tar.Reader, dest string) (err error) {
	path := filepath.Join(dest, hdr.Name)
	info := hdr.FileInfo()

	if info.IsDir() {
		if err = os.MkdirAll(path, info.Mode()); err != nil {
			return
		}
		return
	}

	var file *os.File
	dirPath := filepath.Dir(path)
	if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return
	}
	file, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, tr)
	if err != nil {
		return
	}

	return nil
}
