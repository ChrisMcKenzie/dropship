package updater

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

func NewRackspaceUpdater(user, key, region string) *RackspaceUpdater {
	return &RackspaceUpdater{
		conn: &swift.Connection{
			// This should be your username
			UserName: user,
			// This should be your api key
			ApiKey: key,
			// This should be a v1 auth url, eg
			//  Rackspace US        https://auth.api.rackspacecloud.com/v1.0
			//  Rackspace UK        https://lon.auth.api.rackspacecloud.com/v1.0
			//  Memset Memstore UK  https://auth.storage.memset.com/v1.0
			AuthUrl: "https://auth.api.rackspacecloud.com/v1.0",
			// Region to use - default is use first region if unset
			Region: region,
			// Name of the tenant - this is likely your username
		},
	}
}

func (u *RackspaceUpdater) IsOutdated(hash string, opts *Options) (bool, error) {
	if u.conn == nil {
		return false, ErrUnableToConnect
	}

	info, _, err := u.conn.Object(opts.Bucket, opts.Path)
	if err != nil {
		return false, err
	}

	if info.Hash == hash {
		return false, nil
	}

	return true, nil
}

func (u *RackspaceUpdater) Download(opt *Options) (io.ReadCloser, MetaData, error) {
	var meta MetaData
	if u.conn == nil {
		return nil, meta, ErrUnableToConnect
	}

	r, hdrs, err := u.conn.ObjectOpen(opt.Bucket, opt.Path, true, swift.Headers{})
	if err != nil {
		return nil, meta, err
	}

	meta.Hash = hdrs["Etag"]
	meta.ContentType = hdrs["Content-Type"]

	return r, meta, nil
}
