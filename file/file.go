package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sraphs/go/config"
	"github.com/sraphs/go/x/strslices"
)

var _ config.Source = (*file)(nil)

type file struct {
	path string
}

// NewSource new a file source.
func NewSource(path string) config.Source {
	return &file{path: path}
}

func (f *file) Load() (desc []*config.Descriptor, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}
	des, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	return []*config.Descriptor{des}, nil
}

func (f *file) Watch() (config.Watcher, error) {
	return newWatcher(f)
}

func (f *file) loadFile(path string) (*config.Descriptor, error) {
	if !isSupported(path) {
		return nil, config.ErrUnsupportedFormat
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &config.Descriptor{
		Name:   info.Name(),
		Format: format(info.Name()),
		Data:   data,
	}, nil
}

func (f *file) loadDir(path string) ([]*config.Descriptor, error) {
	files, err := os.ReadDir(f.path)
	if err != nil {
		return nil, err
	}

	var descs = make([]*config.Descriptor, 0, len(files))

	for _, file := range files {
		// ignore hidden files
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		// ignore files that are not supported formats
		if !isSupported(file.Name()) {
			continue
		}

		desc, err := f.loadFile(filepath.Join(f.path, file.Name()))
		if err != nil {
			return nil, err
		}
		descs = append(descs, desc)
	}

	return descs, nil
}

func format(name string) string {
	if p := strings.Split(name, "."); len(p) > 1 {
		return p[len(p)-1]
	}
	return ""
}

func isSupported(path string) bool {
	format := format(path)
	return strslices.Contains(config.SupportedFormats, format)
}
