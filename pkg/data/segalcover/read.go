package segalcover

import (
	"context"
	"path/filepath"

	"github.com/kellegous/scotus/pkg/data/internal"
	"github.com/kellegous/scotus/pkg/data/option"
)

const (
	dataFielName = "segal-cover.html"
	DefaultURL   = `https://en.wikipedia.org/wiki/Segal%E2%80%93Cover_score`
)

func Read(
	ctx context.Context,
	opts ...option.DownloadOption,
) ([]*Justice, error) {
	var o option.DownloadOptions
	o.ApplyOptions(opts, option.FromURL(DefaultURL))

	src := filepath.Join(o.DataDir, dataFielName)

	if err := internal.EnsureDownload(
		ctx,
		o.Client,
		o.URL,
		src,
	); err != nil {
		return nil, err
	}

	return nil, nil
}
