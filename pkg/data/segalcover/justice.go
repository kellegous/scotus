package segalcover

type Justice struct {
	Name          string
	Chief         bool
	Ideology      float64
	YearNominated int
	NominatedBy   *President
}
