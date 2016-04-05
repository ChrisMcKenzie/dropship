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

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func serveInterface(port int) error {
	router := newRouter()
	log.Infof("HTTP Server Listening on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
