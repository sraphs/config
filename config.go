package config

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/sraphs/go/log"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported format")
)

var (
	SupportedFormats = []string{"env", "json", "xml", "yaml", "yml"}
)

// Observer is config observer.
type Observer func(Config)

// Config is a config interface.
type Config interface {
	Load() error
	Scan(v interface{}) error
	Watch(o Observer) error
	Close() error
}

type config struct {
	opts        options
	descriptors sync.Map
	observers   []Observer
	watchers    []Watcher
	mu          sync.Mutex
}

// New new a config with options.
func New(opts ...Option) Config {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &config{
		opts: o,
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		ds, err := src.Load()
		if err != nil {
			return err
		}

		for _, d := range ds {
			log.Debug("load config", "name", d.Name, "format", d.Format)
			c.descriptors.Store(d.Name, d)
		}

		w, err := src.Watch()

		if err != nil {
			log.Error("failed to watch config source", err)
			return err
		}

		if w != nil {
			c.watchers = append(c.watchers, w)
			go c.watch(w)
		}
	}

	return nil
}

func (c *config) Scan(v interface{}) error {
	c.descriptors.Range(func(key, value interface{}) bool {
		d := value.(*Descriptor)
		if err := d.GetCodec().Unmarshal(d.Data, v); err != nil {
			return true
		}
		return true
	})
	return nil
}

func (c *config) Watch(o Observer) error {
	c.observers = append(c.observers, o)
	return nil
}

func (c *config) Close() error {
	for _, w := range c.watchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) watch(w Watcher) {
	for {
		ds, err := w.Next()
		if errors.Is(err, context.Canceled) {
			log.Info("watcher's ctx cancel", err)
			return
		}
		if err != nil {
			time.Sleep(time.Second)
			log.Error("failed to watch next config", err)
			continue
		}
		for _, d := range ds {
			if v, ok := c.descriptors.Load(d.Name); ok {
				if !reflect.DeepEqual(v, d) {
					c.descriptors.Store(d.Name, d)
					for _, o := range c.observers {
						o(c)
					}
				}
			}
		}
	}
}
