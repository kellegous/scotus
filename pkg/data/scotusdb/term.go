package scotusdb

import (
	"fmt"
	"time"

	"github.com/kellegous/scotus/pkg/csv"
)

type Term struct {
	Year  int     `json:"year"`
	Cases []*Case `json:"cases"`
}

func readTerm(
	terms map[int]*Term,
	row *csv.Row,
) (*Term, bool, error) {
	year, err := row.GetInt("term", 0)
	if err != nil {
		return nil, false, err
	}

	if year < 1700 || year > time.Now().Year() {
		return nil, false, fmt.Errorf("invalid term year: %d", year)
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
