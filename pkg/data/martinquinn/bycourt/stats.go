package bycourt

import "github.com/kellegous/scotus/pkg/csv"

type Stats struct {
	MedianJusticeScore    float64
	StdDevOfMedianJustice float64
	MinJusticeScore       float64
	MaxJusticeScore       float64
	MedianJustice         string
}

func parseStats(row *csv.Row) (*Stats, error) {
	med, err := row.GetFloat64("med", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	std, err := row.GetFloat64("med_sd", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	min, err := row.GetFloat64("min", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	max, err := row.GetFloat64("max", csv.ParseFloat64)
	if err != nil {
		return nil, err
	}

	jus, err := row.Get("justice")
	if err != nil {
		return nil, err
	}

	return &Stats{
		MedianJusticeScore:    med,
		StdDevOfMedianJustice: std,
		MinJusticeScore:       min,
		MaxJusticeScore:       max,
		MedianJustice:         jus,
	}, nil
}
