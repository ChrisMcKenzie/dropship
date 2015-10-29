package installer

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrNilReader = errors.New("Install: must have a non-nil Reader")

type TarInstaller struct{}

func (i TarInstaller) Install(dest string, fr io.Reader) (count int, err error) {
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
			break
		}
		if err != nil {
			return
		}
		path := filepath.Join(dest, hdr.Name)
		info := hdr.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return
			}
			continue
		}

		var file *os.File
		file, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return
		}
		defer file.Close()

		_, err = io.Copy(file, tr)
		if err != nil {
			return
		}
		count++
	}

	return
}
