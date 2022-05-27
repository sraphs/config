package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/imdario/mergo"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported format")
	ErrScanNeedPtr       = errors.New("scan need ptr")
	ErrNotFound          = errors.New("key not found")
	ErrTypeAssert        = errors.New("type assert error")
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
	Get(key string) Value
}

var _ Config = (*config)(nil)

type config struct {
	opts        options
	descriptors sync.Map
	observers   []Observer
	watchers    []Watcher
	mu          sync.Mutex
	reader      Reader
	cached      sync.Map
}

// New new a config with options.
func New(opts ...Option) Config {
	o := options{
		decoder:  defaultDecoder,
		resolver: defaultResolver,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &config{
		opts:   o,
		reader: newReader(o),
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		descriptors, err := src.Load()
		if err != nil {
			return err
		}

		for _, d := range descriptors {
			if c.opts.enableLog {
				fmt.Printf("load config: name: %s format: %s\n", d.Name, d.Format)
			}
			c.descriptors.Store(d.Name, d)
		}

		if err = c.reader.Merge(descriptors...); err != nil {
			return fmt.Errorf("failed to watch config source: %v", err)
		}

		w, err := src.Watch()

		if err != nil {
			return fmt.Errorf("failed to watch config source: %v", err)
		}

		if w != nil {
			c.watchers = append(c.watchers, w)
			go c.watch(w)
		}
	}

	if err := c.reader.Resolve(); err != nil {
		return fmt.Errorf("failed to resolve config source: %v", err)
	}

	return nil
}

func (c *config) Scan(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return ErrScanNeedPtr
	}

	rvt := reflect.TypeOf(v).Elem()

	var envV = reflect.New(rvt).Interface()
	if ed, ok := c.descriptors.LoadAndDelete("environ"); ok {
		c.scanTo(ed.(*Descriptor), envV)
	}

	var flagV = reflect.New(rvt).Interface()
	if fd, ok := c.descriptors.LoadAndDelete("flag"); ok {
		c.scanTo(fd.(*Descriptor), flagV)
	}

	var fileV = reflect.New(rvt).Interface()
	data, err := c.reader.Source()
	if err != nil {
		return err
	}
	if err := unmarshalJSON(data, fileV); err != nil {
		return err
	}

	// order config source: flag > env > file
	if err := mergo.Merge(v, fileV, mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(v, envV, mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(v, flagV, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

func (c *config) scanTo(d *Descriptor, v interface{}) error {
	codec := d.GetCodec()
	if codec == nil {
		log.Printf("failed to get codec: name: %s format: %s", d.Name, d.Format)
		return fmt.Errorf("failed to get codec: name: %s format: %s", d.Name, d.Format)
	}
	nv := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	if err := codec.Unmarshal(d.Data, nv); err != nil {
		return fmt.Errorf("failed to unmarshal config: name: %s format: %s err: %v", d.Name, d.Format, err)

	}
	if err := mergo.Merge(v, nv); err != nil {
		return fmt.Errorf("failed to merge config: name: %s format: %s err: %v", d.Name, d.Format, err)
	}
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

func (c *config) Get(key string) Value {
	if v, ok := c.cached.Load(key); ok {
		return v.(Value)
	}
	if v, ok := c.reader.Value(key); ok {
		c.cached.Store(key, v)
		return v
	}
	return &errValue{err: ErrNotFound}
}

func (c *config) watch(w Watcher) {
	for {
		descriptors, err := w.Next()
		if errors.Is(err, context.Canceled) {
			fmt.Println("watcher's ctx cancel", err)
			return
		}
		if err != nil {
			time.Sleep(time.Second)
			fmt.Println("failed to watch next config", err)
			continue
		}
		if err := c.reader.Merge(descriptors...); err != nil {
			fmt.Println("failed to merge next config", err)
			continue
		}
		if err := c.reader.Resolve(); err != nil {
			fmt.Println("failed to resolve next config", err)
			continue
		}

		c.cached.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(Value)
			if n, ok := c.reader.Value(k); ok && reflect.TypeOf(n.Load()) == reflect.TypeOf(v.Load()) && !reflect.DeepEqual(n.Load(), v.Load()) {
				v.Store(n.Load())
				c.cached.Store(k, v)
			}
			return true
		})

		for _, d := range descriptors {
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
