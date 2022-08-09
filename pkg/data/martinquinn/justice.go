package martinquinn

import (
	"strconv"

	"github.com/kellegous/scotus/pkg/csv"
)

type Justice struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"stddev"`
	Median float64 `json:"median"`
}

func parseJustice(row *csv.Row) (*Justice, error) {
	id, err := row.GetInt("justice", strconv.Atoi)
	if err != nil {
		return nil, err
	}

	name, err := row.Get("justiceName")
	if err != nil {
		return nil, err
	}

	mean, err := row.GetFloat64("post_mn", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	stddev, err := row.GetFloat64("post_sd", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	median, err := row.GetFloat64("post_med", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	return &Justice{
		ID:     id,
		Name:   name,
		Mean:   mean,
		StdDev: stddev,
		Median: median,
	}, nil
}
