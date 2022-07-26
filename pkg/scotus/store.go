package scotus

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kellegous/scotus/pkg/csv"
)

const casesFileName = "SCDB_justiceCentered_Citation.csv.zip"

type Store struct {
	dir string
}

func (s *Store) Terms() ([]*Term, error) {
	zr, err := zip.OpenReader(filepath.Join(s.dir, casesFileName))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	if n := len(zr.File); n != 1 {
		return nil, fmt.Errorf("expected a single file but there are %d", n)
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

		term, added, err := readTerm(termsByYear, row)
		if err != nil {
			return nil, err
		}

		if added {
			terms = append(terms, term)
		}

		c, added, err := readCase(casesByID, row)
		if err != nil {
			return nil, err
		}

		if added {
			term.Cases = append(term.Cases, c)
		}

		vote, err := readVote(row)
		if err != nil {
			return nil, err
		}

		c.Votes = append(c.Votes, vote)
	}

	return terms, nil
}

func OpenStore(
	ctx context.Context,
	dir string,
	opts ...Option,
) (*Store, error) {
	o := Options{
		casesURL: DefaultCasesURL,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if err := prepareDir(ctx, dir, &o); err != nil {
		return nil, err
	}

	return &Store{
		dir: dir,
	}, nil
}

func prepareDir(
	ctx context.Context,
	dir string,
	opts *Options,
) error {
	if err := ensureDir(dir, opts.reset); err != nil {
		return err
	}

	if err := ensureCaseDownload(
		ctx,
		opts.casesURL,
		filepath.Join(dir, casesFileName),
	); err != nil {
		return err
	}

	return nil
}

func ensureDir(dir string, reset bool) error {
	if _, err := os.Stat(dir); err != nil {
		return os.MkdirAll(dir, 0755)
	} else if reset {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
		return os.MkdirAll(dir, 0755)
	}

	return nil
}

func ensureCaseDownload(
	ctx context.Context,
	src string,
	dst string,
) error {
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		src,
		nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if s := res.StatusCode; s != http.StatusOK {
		return fmt.Errorf("http status %d", s)
	}

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}

	return nil
}
