package dropship

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
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
		dir, err := ioutil.TempDir(".", "test")
		if err != nil {
			t.Error(err)
		}
		defer os.RemoveAll(dir)

		count, err := fileInstaller.Install(dir+"/test.txt", test.file)
		if err != test.err {
			t.Errorf("Install: Expected error to equal %v got: %v", test.err, err)
		}

		if count != test.count {
			t.Errorf("Install: Expected %d files to be installed got %v", test.count, count)
		}
	}
}
