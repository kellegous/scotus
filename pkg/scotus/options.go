package scotus

const DefaultCasesURL = "http://scdb.wustl.edu/_brickFiles/2021_01/SCDB_2021_01_justiceCentered_Citation.csv.zip"

type Options struct {
	casesURL string
	reset    bool
}

func (o *Options) applyDefaults() {
	if o.casesURL == "" {
		o.casesURL = DefaultCasesURL
	}
}

type Option func(o *Options)

func WithCasesFromURL(url string) Option {
	return func(o *Options) {
		o.casesURL = url
	}
}

func Reset(v bool) Option {
	return func(o *Options) {
		o.reset = v
	}
}
