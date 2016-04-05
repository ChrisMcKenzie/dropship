package manager

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
)

const (
	DefaultKeyPrefix = "dropship"
)

func init() {
	consul.Register()
	etcd.Register()
	zookeeper.Register()
	boltdb.Register()
}

func initStore(storeUrl *url.URL) (store.Store, error) {
	if storeUrl.Scheme == "" {
		return nil, errors.New("unable to determine store type please make sure you have added a scheme (ie. consul:// or etcd://)")
	}

	hosts := strings.Split(storeUrl.Host, ",")

	return libkv.NewStore(
		store.Backend(storeUrl.Scheme),
		hosts,
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
}
