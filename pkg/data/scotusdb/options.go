package scotusdb

import (
	"net/http"
)

const (
	DefaultModernCasesURL = "http://scdb.wustl.edu/_brickFiles/2021_01/SCDB_2021_01_justiceCentered_Citation.csv.zip"
	DefaultLegacyCasesURL = "http://scdb.wustl.edu/_brickFiles/Legacy_07/SCDB_Legacy_07_justiceCentered_Citation.csv.zip"
	DefaultOT21CasesURL   = "https://gist.githubusercontent.com/kellegous/cbb09234ed108700162ee80be52780c9/raw/d7bed134a15a82db901678d81e4a2b4a34237136/ot21.json"
	DefaultDataDir        = "data"
)

type Options struct {
	legacyCasesURL string
	modernCasesURL string
	ot21CasesURL   string
	dataDir        string
	client         *http.Client
}

func (o *Options) apply(opts []Option) {
	o.legacyCasesURL = DefaultLegacyCasesURL
	o.modernCasesURL = DefaultModernCasesURL
	o.ot21CasesURL = DefaultOT21CasesURL
	o.dataDir = DefaultDataDir
	o.client = http.DefaultClient
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(o *Options)

func WithCaseURLs(
	modernCasesURL string,
	legacyCasesURL string,
	ot21CasesURL string,
) Option {
	return func(o *Options) {
		o.modernCasesURL = modernCasesURL
		o.legacyCasesURL = legacyCasesURL
		o.ot21CasesURL = ot21CasesURL
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
