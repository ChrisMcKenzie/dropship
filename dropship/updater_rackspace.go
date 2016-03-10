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
	"errors"
	"io"

	"github.com/ncw/swift"
)

var (
	ErrUnableToConnect = errors.New("RackspaceUpdater: unable to connect to rackspace.")
)

type RackspaceUpdater struct {
	conn *swift.Connection
}

func NewRackspaceUpdater(config map[string]string) *RackspaceUpdater {
	return &RackspaceUpdater{
		conn: &swift.Connection{
			// This should be your username
			UserName: config["user"],
			// This should be your api key
			ApiKey: config["key"],
			// This should be a v1 auth url, eg
			//  Rackspace US        https://auth.api.rackspacecloud.com/v1.0
			//  Rackspace UK        https://lon.auth.api.rackspacecloud.com/v1.0
			//  Memset Memstore UK  https://auth.storage.memset.com/v1.0
			AuthUrl: "https://auth.api.rackspacecloud.com/v1.0",
			// Region to use - default is use first region if unset
			Region: config["region"],
			// Name of the tenant - this is likely your username
		},
	}
}

func (u *RackspaceUpdater) IsOutdated(hash string, opts Artifact) (bool, error) {
	if u.conn == nil {
		return false, ErrUnableToConnect
	}

	if _, ok := opts["bucket"]; !ok {
		return false, errors.New("Missing field: \"bucket\"")
	}

	if _, ok := opts["path"]; !ok {
		return false, errors.New("Missing field: \"path\"")
	}

	info, _, err := u.conn.Object(opts["bucket"], opts["path"])
	if err != nil {
		return false, err
	}

	if info.Hash == hash {
		return false, nil
	}

	return true, nil
}

func (u *RackspaceUpdater) Download(opt Artifact) (io.ReadCloser, MetaData, error) {
	var meta MetaData
	if u.conn == nil {
		return nil, meta, ErrUnableToConnect
	}

	r, hdrs, err := u.conn.ObjectOpen(opt["bucket"], opt["path"], true, swift.Headers{})
	if err != nil {
		return nil, meta, err
	}

	meta.Hash = hdrs["Etag"]
	meta.ContentType = hdrs["Content-Type"]

	return r, meta, nil
}
