package agent

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func untar(fr io.Reader, dest string) error {
	gr, err := gzip.NewReader(fr)
	defer gr.Close()
	if err != nil && err != io.EOF {
		return err
	}

	tr := tar.NewReader(gr)
	if err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		path := filepath.Join(dest, hdr.Name)
		info := hdr.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				panic(err)
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(file, tr)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Finished Installing")
	return nil
}

func overwrite(mpath string) (*os.File, error) {
	f, err := os.Create(mpath)
	return f, err
}
