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
	"path/filepath"
	"strings"
)

// Installer is an interface that allows different methods of writing
// the given io.Reader to disk.
type Installer interface {
	Install(dest string, r io.Reader) (int, error)
}

func moveOld(dest string) error {
	return os.Rename(dest, strings.Join([]string{dest, "old"}, "."))
}

func cleanup(dest string, err error) error {
	oldPath := filepath.Join(dest, ".old")
	if err != nil {
		return os.Rename(oldPath, dest)
	}

	return os.RemoveAll(oldPath)
}
