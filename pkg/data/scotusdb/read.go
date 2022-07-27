package scotusdb

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/kellegous/scotus/pkg/csv"
	"github.com/kellegous/scotus/pkg/data/internal"
)

const casesFileName = "SCDB_justiceCentered_Citation.csv.zip"

func Read(
	ctx context.Context,
	opts ...Option,
) ([]*Term, error) {
	var o Options
	o.apply(opts)

	src := filepath.Join(o.dataDir, casesFileName)

	if err := internal.EnsureDownload(
		ctx,
		o.client,
		o.casesURL,
		src,
	); err != nil {
		return nil, err
	}

	return readTerms(src)
}

func readTerms(src string) ([]*Term, error) {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	if n := len(zr.File); n != 1 {
		return nil, fmt.Errorf(
			"expected a single file but there are %d",
			n)
	}

	r, err := zr.File[0].Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	cr, err := csv.NewReader(r)
	if err != nil {
		return nil, err
	}

	termsByYear := map[int]*Term{}
	casesByID := map[string]*Case{}
	var terms []*Term

	for {
		row, err := cr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		t, added, err := readTerm(termsByYear, row)
		if err != nil {
			return nil, err
		}

		if added {
			terms = append(terms, t)
		}

		c, added, err := readCase(casesByID, row)
		if err != nil {
			return nil, err
		}

		if added {
			t.Cases = append(t.Cases, c)
		}

		v, err := readVote(row)
		if err != nil {
			return nil, err
		}

		c.Votes = append(c.Votes, v)
	}

	return terms, nil
}
