package overrulings

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/kellegous/scotus/pkg/data/internal"
)

const dataFileName = "overrulings.html"

func Read(
	ctx context.Context,
	opts ...Option,
) (interface{}, error) {
	var o Options
	o.apply(opts)

	src := filepath.Join(o.dataDir, dataFileName)

	if err := internal.EnsureDownload(
		ctx,
		o.client,
		o.url,
		src,
	); err != nil {
		return nil, err
	}

	return read(src)
}

func read(src string) (interface{}, error) {
	r, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if _, err := io.Copy(io.Discard, r); err != nil {
		return nil, err
	}

	return nil, nil
}
