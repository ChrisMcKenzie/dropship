package dropship

import (
	"errors"
	"io"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

// S3Updater is an Updater for updating file that come from amazone s3
type S3Updater struct {
	conn *s3.S3
}

// NewS3Updater will return a configured Updater for check against amazon s3
func NewS3Updater(config map[string]string) *S3Updater {
	return &S3Updater{
		s3.New(
			aws.Auth{
				AccessKey: config["accessKey"],
				SecretKey: config["secret"],
				Token:     config["token"],
			},
			aws.Region{
				Name:       config["name"],
				S3Endpoint: config["endpoint"],
			},
		),
	}
}

// Download is an Updater method that will download the requested file as an
// io.ReadCloser and return MetaData to the use about the file.
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

// IsOutdated will check if a file is out of date by getting the ETag of the
// file and comparing it withe the current file.
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
