package dropship

import (
	"bytes"
	"io"
	"testing"
)

func TestInstallFile(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte("hello"))

	cases := []struct {
		file  io.Reader
		count int
		err   error
	}{
		{nil, 0, ErrNilReader},
		{&buf, 1, nil},
	}

	var fileInstaller FileInstaller
	for _, test := range cases {
		count, err := fileInstaller.Install("/tmp/test.txt", test.file)
		if err != test.err {
			t.Errorf("Install: Expected error to equal %v got: %v", test.err, err)
		}

		if count != test.count {
			t.Errorf("Install: Expected % files to be installed got %v", test.count, count)
		}
	}
}
