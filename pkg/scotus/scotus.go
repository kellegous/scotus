package scotus

import (
	"io"
	"os"
	"strconv"

	"github.com/kellegous/scotus/pkg/csv"
)

func readTerm(
	terms map[int]*Term,
	row *csv.Row,
) (*Term, bool, error) {
	term, err := row.Get("term")
	if err != nil {
		return nil, false, err
	}

	year, err := strconv.Atoi(term)
	if err != nil {
		return nil, false, err
	}

	if t := terms[year]; t != nil {
		return t, false, nil
	}

	t := &Term{
		Year: year,
	}
	terms[year] = t

	return t, true, nil
}

func ReadFile(src string) ([]*Term, error) {
	r, err := os.Open(src)
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
	}

	return terms, nil
}
