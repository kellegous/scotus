package byjustice

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/kellegous/scotus/pkg/csv"
	"github.com/kellegous/scotus/pkg/data/internal"
	"github.com/kellegous/scotus/pkg/data/option"
)

const (
	DefaultURL = "https://mqscores.lsa.umich.edu/media/2020/justices.csv"

	filename = "martinquinn-justices.csv"
)

func Read(
	ctx context.Context,
	opts ...option.DownloadOption,
) ([]*Term, error) {
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

func read(r io.Reader) ([]*Term, error) {
	cr, err := csv.NewReader(r)
	if err != nil {
		return nil, err
	}

	byYear := map[int]*Term{}
	var terms []*Term
	for {
		row, err := cr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		year, err := row.GetInt("term", strconv.Atoi)
		if err != nil {
			return nil, err
		}

		term := byYear[year]
		if term == nil {
			term = &Term{Year: year}
			byYear[year] = term
			terms = append(terms, term)
		}

		justice, err := parseJustice(row)
		if err != nil {
			return nil, err
		}

		term.Justices = append(term.Justices, justice)
	}

	sort.Slice(terms, func(i, j int) bool {
		return terms[i].Year < terms[j].Year
	})

	for _, term := range terms {
		justices := term.Justices
		sort.Slice(
			justices,
			func(i, j int) bool {
				return justices[i].Name < justices[j].Name
			})
	}

	return terms, nil
}
