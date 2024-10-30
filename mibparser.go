package mibparser

// Option is a configuration function type
type Option func(*Opts)

// Opts holds configuration options
type Opts struct {
	Path string
}

// MIBParser stores configuration options
type MIBParser struct {
	opts Opts
}

// NewPath sets the path
func NewPath(path string) Option {
	return func(opts *Opts) {
		opts.Path = path
	}
}
