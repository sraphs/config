package config

import (
	"github.com/sraphs/go/encoding"
)

// Source is config source.
type Source interface {
	Load() ([]*Descriptor, error)
	Watch() (Watcher, error)
}

// Watcher watches a source for changes.
type Watcher interface {
	Next() ([]*Descriptor, error)
	Stop() error
}

// Descriptor is file or env or flag descriptor.
type Descriptor struct {
	Name   string
	Format string
	Data   []byte
}

func (d *Descriptor) GetCodec() encoding.Codec {
	return encoding.GetCodec(d.Format)
}
