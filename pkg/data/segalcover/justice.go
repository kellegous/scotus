package segalcover

type Justice struct {
	Name          string     `json:"name"`
	Chief         bool       `json:"as_chief"`
	Ideology      float64    `json:"ideology"`
	YearNominated int        `json:"year_nominated"`
	NominatedBy   *President `json:"nominated_by"`
}
