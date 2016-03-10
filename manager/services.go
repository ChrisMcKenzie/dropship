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
package manager

import "net/http"

func servicesIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("hello, world"))
}

func servicesShow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("hello, world"))
}
