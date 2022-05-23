package flag

import (
	"bytes"
	"os"

	"github.com/sraphs/config"
)

var _ config.Source = (*flag)(nil)

type flag struct {
	args []string
}

func NewSource() config.Source {
	return &flag{args: os.Args[1:]}
}

func (f *flag) Load() ([]*config.Descriptor, error) {
	var buf bytes.Buffer

	for _, arg := range f.args {
		buf.WriteString(arg)
		buf.WriteByte(' ')
	}

	d := &config.Descriptor{
		Name:   "flag",
		Data:   buf.Bytes(),
		Format: "flag",
	}

	return []*config.Descriptor{d}, nil
}

func (f *flag) Watch() (config.Watcher, error) {
	return nil, nil
}
