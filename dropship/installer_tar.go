// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package dropship

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrNilReader = errors.New("Install: must have a non-nil Reader")

// TarInstaller Defines an install Method that takes a destination path
// and a io.Reader and untars and gzip decodes a tarball and
// places the files inside on the FS with `dest` as their root
// It returns the number of files written and an error
type TarInstaller struct{}

func (i TarInstaller) Install(dest string, fr io.Reader) (count int, err error) {
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

	defer cleanup(dest, err)
	return
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
