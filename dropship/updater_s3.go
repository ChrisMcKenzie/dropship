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

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type S3Updater struct {
	conn *s3.S3
}

func NewS3Updater(config map[string]string) *S3Updater {
	return &S3Updater{
		s3.New(
			aws.Auth{config["accessKey"], config["secret"], config["token"]},
			aws.Region{
				Name:       config["name"],
				S3Endpoint: config["endpoint"],
			},
		),
	}
}

func (u S3Updater) Download(config Artifact) (io.ReadCloser, MetaData, error) {
	if _, ok := config["bucket"]; !ok {
		return nil, MetaData{}, errors.New("bucket name is required")
	}

	path := config["path"]

	bucket := u.conn.Bucket(config["bucket"])

	// Get file meta-data
	res, err := bucket.GetResponse(path)
	if err != nil {
		return nil, MetaData{}, err
	}

	meta := MetaData{
		ContentType: res.Header.Get("Content-Type"),
		Hash:        res.Header.Get("Etag"),
	}

	return res.Body, meta, nil
}

func (u S3Updater) IsOutdated(hash string, config Artifact) (bool, error) {
	if _, ok := config["bucket"]; !ok {
		return false, errors.New("bucket name is required")
	}

	path := config["path"]

	bucket := u.conn.Bucket(config["bucket"])

	// Get file meta-data
	res, err := bucket.GetKey(path)
	if err != nil {
		return false, err
	}

	if res.ETag == hash {
		return false, nil
	}

	return true, nil
}
