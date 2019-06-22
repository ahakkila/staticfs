package staticfs

import (
	"fmt"
	"net/http"
	"strings"
)

type SFSOption int

const (
	DirectoryFilter SFSOption = 1 << iota // 0x00000001
	DotFileFilter                         // 0x00000010
)

type StaticFileSystem struct {
	fs      http.FileSystem
	options SFSOption
}

func NewStaticFileSystem(path string, options SFSOption) StaticFileSystem {
	sfs := StaticFileSystem{fs: http.Dir(path), options: options}
	return sfs
}

func (sfs StaticFileSystem) filterDirectory(f http.File, path string) error {
	s, err := f.Stat()
	if err != nil {
		return err
	}

	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := sfs.fs.Open(index); err != nil {
			return err
		}
	}
	return nil
}

func (sfs StaticFileSystem) filterDotFiles(path string) error {
	parts := strings.Split(path, "/")
	for ix, part := range parts {
		if strings.HasPrefix(part, ".") {
			return fmt.Errorf("part %d (%s) contains a dot prefix in path %s", ix, part, path)
		}
	}
	return nil
}

func (sfs StaticFileSystem) Open(path string) (http.File, error) {

	if (sfs.options & DotFileFilter) > 0 {
		if err := sfs.filterDotFiles(path); err != nil {
			return nil, err
		}
	}

	file, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	if (sfs.options & DirectoryFilter) > 0 {
		if err := sfs.filterDirectory(file, path); err != nil {
			return nil, err
		}
	}

	return file, nil
}
