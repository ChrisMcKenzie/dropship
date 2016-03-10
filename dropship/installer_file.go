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
	"io"
	"os"
)

// FileInstaller defines an Installer that takes the reader and writes
// it to the dest directory.
type FileInstaller struct{}

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
