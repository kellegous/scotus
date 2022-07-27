package overrulings

import "net/http"

const (
	DefaultURL     = "https://constitution.congress.gov/resources/decisions-overruled/"
	DefaultDataDir = "data"
)

type Options struct {
	url     string
	dataDir string
	client  *http.Client
}

func (o *Options) apply(opts []Option) {
	o.url = DefaultURL
	o.dataDir = DefaultDataDir
	o.client = http.DefaultClient
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(o *Options)

func FromURL(url string) Option {
	return func(o *Options) {
		o.url = url
	}
}

func WithDataDir(dir string) Option {
	return func(o *Options) {
		o.dataDir = dir
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(o *Options) {
		o.client = c
	}
}
