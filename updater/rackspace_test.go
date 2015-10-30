package updater

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ncw/swift"
)

var (
	user, key, region string
	content           string = "hello world\r\n"
	hash              string = "a0f2a3c1dcd5b1cac71bf0c03f2ff1bd"
	conn              *swift.Connection
	updater           Updater
)

const Container = "test-container"

func setup() error {
	user = os.Getenv("RACKSPACE_USER")
	key = os.Getenv("RACKSPACE_KEY")
	region = os.Getenv("RACKSPACE_REGION")
	if user == "" || key == "" || region == "" {
		return fmt.Errorf("user or key are required")
	}

	conn = &swift.Connection{
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
	}

	// setup test container
	err := conn.ContainerCreate(Container, nil)
	if err != nil {
		return err
	}

	obj, err := conn.ObjectCreate(Container, "test.txt", false, "", "text/plain", nil)
	if err != nil {
		return err
	}

	_, err = obj.Write([]byte(content))

	err = obj.Close()

	return err
}

func TestMain(t *testing.T) {
	err := setup()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	updater = NewRackspaceUpdater(user, key, region)
}

func TestRackspaceUpdaterIsOutdated(t *testing.T) {
	cases := []struct {
		hash     string
		expected bool
	}{
		{
			hash:     "",
			expected: true,
		},
		{
			hash:     hash,
			expected: false,
		},
		{
			hash:     "this hash is random!",
			expected: true,
		},
	}

	for _, test := range cases {
		result, err := updater.IsOutdated(test.hash, &Options{Container, "test.txt"})
		if err != nil {
			t.Error(err)
		}

		if result != test.expected {
			t.Error(fmt.Errorf("IsOutdated:Test hash %v Expected %t got %t", test.hash, test.expected, result))
		}
	}

	// negative non-existing file error
	_, err := updater.IsOutdated("", &Options{Container, "no-file.txt"})
	if err == nil {
		t.Error("IsOutdated: Expected an error when accessing non-existent file.")
	}

	// negative generic error testing.
	badUpdater := &RackspaceUpdater{}
	_, err = badUpdater.IsOutdated("", &Options{})
	if err == nil {
		t.Error("IsOutdated: Expected an error with invalid credentials.")
	}
}

func TestRackspaceUpdaterDownload(t *testing.T) {
	cases := []struct {
		Options  *Options
		Value    string
		MetaData MetaData
		Error    error
	}{
		{
			Options: &Options{
				Bucket: Container,
				Path:   "test.txt",
			},
			MetaData: MetaData{"text/plain", hash},
			Error:    nil,
			Value:    content,
		},
		{
			Options: &Options{
				Bucket: Container,
				Path:   "no-file.txt",
			},
			Error: swift.ObjectNotFound,
			Value: "",
		},
	}

	for _, test := range cases {
		r, meta, err := updater.Download(test.Options)
		if err != test.Error {
			t.Errorf("Download: Expected error to be %v got: %v", test.Error, err)
		}

		if meta != test.MetaData {
			t.Errorf("Download: Expected MetaData to be %v got: %v", test.MetaData, meta)
		}

		if r != nil {
			var data bytes.Buffer
			_, err = io.Copy(&data, r)
			if err != nil {
				t.Error(err)
			}

			if data.String() != test.Value {
				t.Errorf("Download: Expected io.Reader to contain %v got: %v", test.Value, data.String())
			}
		}
	}

	// negative generic error testing.
	badUpdater := &RackspaceUpdater{}
	_, _, err := badUpdater.Download(&Options{})
	if err == nil {
		t.Error("IsOutdated: Expected an error with invalid credentials.")
	}
}
