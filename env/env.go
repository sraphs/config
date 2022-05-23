package env

import (
	"bytes"
	"os"
	"strings"

	"github.com/sraphs/config"
)

var _ config.Source = (*env)(nil)

type env struct {
	prefix string
}

func NewSource(prefix string) config.Source {
	return &env{prefix: prefix}
}

func (s *env) Load() ([]*config.Descriptor, error) {
	var buf bytes.Buffer

	environ := os.Environ()

	for _, line := range environ {
		if strings.HasPrefix(line, s.prefix) {
			line = strings.TrimPrefix(line, s.prefix)
			line = strings.TrimPrefix(line, "_")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}

	return []*config.Descriptor{{
		Name:   "environ",
		Format: "env",
		Data:   buf.Bytes(),
	}}, nil
}

func (s *env) Watch() (config.Watcher, error) {
	return nil, nil
}
