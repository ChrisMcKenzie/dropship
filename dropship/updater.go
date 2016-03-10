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

import "io"

type (
	MetaData struct {
		ContentType string
		Hash        string
	}

	// Updater is an interface that defines methods for checking a files
	// freshness and downloading an updated version if needed.
	Updater interface {
		IsOutdated(hash string, opts Artifact) (bool, error)
		Download(Artifact) (io.ReadCloser, MetaData, error)
	}
)
