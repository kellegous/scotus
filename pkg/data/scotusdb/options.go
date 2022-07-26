package scotusdb

import "net/http"

const (
	DefaultCasesURL = "http://scdb.wustl.edu/_brickFiles/2021_01/SCDB_2021_01_justiceCentered_Citation.csv.zip"
	DefaultDataDir  = "data"
)

type Options struct {
	casesURL string
	dataDir  string
	client   *http.Client
}

func (o *Options) apply(opts []Option) {
	o.casesURL = DefaultCasesURL
	o.dataDir = DefaultDataDir
	o.client = http.DefaultClient
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(o *Options)

func WithCasesFromURL(url string) Option {
	return func(o *Options) {
		o.casesURL = url
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
