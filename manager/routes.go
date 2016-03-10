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
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name    string
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		// Wrap routes in log decorator
		handler = route.Handler
		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = []route{
	{
		Name:    "ServicesIndex",
		Path:    "/services/",
		Method:  "GET",
		Handler: servicesIndex,
	},
	{
		Name:    "ServicesShow",
		Path:    "/services/{service}",
		Method:  "GET",
		Handler: servicesShow,
	},
}
