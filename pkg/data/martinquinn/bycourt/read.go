package bycourt

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kellegous/scotus/pkg/csv"
	"github.com/kellegous/scotus/pkg/data/internal"
	"github.com/kellegous/scotus/pkg/data/option"
)

const (
	DefaultURL = "https://mqscores.lsa.umich.edu/media/2020/court.csv"

	filename = "martinquinn-courts.csv"
)

func Read(
	ctx context.Context,
	opts ...option.DownloadOption,
) ([]*Court, error) {
	var o option.DownloadOptions
	o.ApplyOptions(opts, option.FromURL(DefaultURL))

	src := filepath.Join(o.DataDir, filename)

	if err := internal.EnsureDownload(
		ctx,
		o.Client,
		o.URL,
		src,
	); err != nil {
		return nil, err
	}

	r, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return read(r)
}

func read(r io.Reader) ([]*Court, error) {
	cr, err := csv.NewReader(r)
	if err != nil {
		return nil, err
	}

	byYear := map[int]*Court{}
	var courts []*Court

	for {
		row, err := cr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		year, err := row.GetInt(
			"term",
			func(s string) (int, error) {
				s = strings.TrimSpace(s)
				if len(s) < 4 {
					return 0, fmt.Errorf("invalid year: %s", s)
				}
				return strconv.Atoi(s[:4])
			})
		if err != nil {
			return nil, err
		}

		court := byYear[year]
		if court == nil {
			court = &Court{
				Year: year,
			}
			byYear[year] = court
			courts = append(courts, court)
		}

		stats, err := parseStats(row)
		if err != nil {
			return nil, err
		}

		court.Stats = append(court.Stats, stats)
	}

	return courts, nil
}
