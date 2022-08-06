package scotusdb

import "github.com/kellegous/scotus/pkg/csv"

type Term struct {
	Year  int     `json:"year"`
	Cases []*Case `json:"cases"`
}

func readTerm(
	terms map[int]*Term,
	row *csv.Row,
) (*Term, bool, error) {
	year, err := row.GetInt("term")
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
