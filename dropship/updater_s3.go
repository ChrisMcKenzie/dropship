package dropship

import (
	"io"

	"gopkg.in/amz.v1/s3"
)

type S3Updater struct {
	conn *s3.S3
}

func NewS3Updater(config map[string]string) *S3Updater {
	return nil
}

func (u S3Updater) Download(config Artifact) (io.ReadCloser, MetaData, error) {
	return nil, MetaData{}, nil
}

func (u S3Updater) IsOutdated(hash string, opts Artifact) (bool, error) {
	return false, nil
}
