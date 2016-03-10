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
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestTarInstall(t *testing.T) {
	buf, err := CreateTar()
	if err != nil {
		t.Error(err)
	}

	var badGzip bytes.Buffer
	badGzip.Write([]byte("hello"))

	var badTar bytes.Buffer
	gw := gzip.NewWriter(&badTar)
	gw.Write([]byte("hello"))
	defer gw.Close()

	cases := []struct {
		tarball io.Reader
		count   int
		err     error
	}{
		{&buf, 3, nil},
		{nil, 0, ErrNilReader},
		{&badGzip, 0, io.ErrUnexpectedEOF},
		{&badTar, 0, io.ErrUnexpectedEOF},
	}

	var tarInstaller TarInstaller
	for _, test := range cases {
		count, err := tarInstaller.Install("/tmp/", test.tarball)
		if err != test.err {
			t.Errorf("Install: Expected error to equal %v got: %v", test.err, err)
		}

		if count != test.count {
			t.Errorf("Install: Expected % files to be installed got %v", test.count, count)
		}
	}
}

func CreateTar() (buf bytes.Buffer, err error) {
	// Create a new tar archive.
	gw := gzip.NewWriter(&buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence."},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err = tw.WriteHeader(hdr); err != nil {
			return
		}
		if _, err = tw.Write([]byte(file.Body)); err != nil {
			return
		}
	}

	return
}
