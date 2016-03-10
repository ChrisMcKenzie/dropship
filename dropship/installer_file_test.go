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
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestInstallFile(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte("hello"))

	cases := []struct {
		file  io.Reader
		count int
		err   error
	}{
		{nil, 0, ErrNilReader},
		{&buf, 1, nil},
	}

	var fileInstaller FileInstaller
	for _, test := range cases {
		dir, err := ioutil.TempDir(".", "test")
		if err != nil {
			t.Error(err)
		}
		defer os.RemoveAll(dir)

		count, err := fileInstaller.Install(dir+"/test.txt", test.file)
		if err != test.err {
			t.Errorf("Install: Expected error to equal %v got: %v", test.err, err)
		}

		if count != test.count {
			t.Errorf("Install: Expected % files to be installed got %v", test.count, count)
		}
	}
}
