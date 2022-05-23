package config

// Option is config option.
type Option func(*options)

type options struct {
	sources []Source
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}
