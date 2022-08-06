package scotusdb

import (
	"time"

	"github.com/kellegous/scotus/pkg/csv"
)

type Case struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	MajorityVotes int       `json:"majority-votes"`
	MinorityVotes int       `json:"minority-votes"`
	DecisionDate  time.Time `json:"decision-date"`
	Votes         []*Vote   `json:"votes"`
}

func readCase(
	cases map[string]*Case,
	row *csv.Row,
) (*Case, bool, error) {
	id, err := row.Get("caseId")
	if err != nil {
		return nil, false, err
	}

	if c := cases[id]; c != nil {
		return c, false, nil
	}

	name, err := row.Get("caseName")
	if err != nil {
		return nil, false, err
	}

	majVotes, err := row.GetInt("majVotes")
	if err != nil {
		return nil, false, err
	}

	minVotes, err := row.GetInt("minVotes")
	if err != nil {
		return nil, false, err
	}

	descDate, err := row.GetDate("dateDecision")
	if err != nil {
		return nil, false, err
	}

	c := &Case{
		ID:            id,
		Name:          name,
		MajorityVotes: majVotes,
		MinorityVotes: minVotes,
		DecisionDate:  descDate,
	}
	cases[id] = c

	return c, true, nil
}
