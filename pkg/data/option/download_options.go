package option

import "net/http"

const (
	DefaultDataDir = "data"
)

type DownloadOptions struct {
	URL     string
	DataDir string
	Client  *http.Client
}

func (o *DownloadOptions) ApplyOptions(
	opts []DownloadOption,
	defs ...DownloadOption,
) {
	o.DataDir = DefaultDataDir
	o.Client = http.DefaultClient
	for _, opt := range defs {
		opt(o)
	}
	for _, opt := range opts {
		opt(o)
	}
}

type DownloadOption func(o *DownloadOptions)

func FromURL(url string) DownloadOption {
	return func(o *DownloadOptions) {
		o.URL = url
	}
}

func WithDataDir(dir string) DownloadOption {
	return func(o *DownloadOptions) {
		o.DataDir = dir
	}
}

func WithHTTPClient(client *http.Client) DownloadOption {
	return func(o *DownloadOptions) {
		o.Client = client
	}
}
