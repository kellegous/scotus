package scotusdb

import "github.com/kellegous/scotus/pkg/csv"

type Vote struct {
	ID          string   `json:"id"`
	JusticeName string   `json:"justice-name"`
	Decision    Decision `json:"decision"`
}

func readVote(row *csv.Row) (*Vote, error) {
	id, err := row.Get("voteId")
	if err != nil {
		return nil, err
	}

	justiceName, err := row.Get("justiceName")
	if err != nil {
		return nil, err
	}

	majority, err := row.Get("majority")
	if err != nil {
		return nil, err
	}

	des, err := decisionFromString(majority)
	if err != nil {
		return nil, err
	}

	return &Vote{
		ID:          id,
		JusticeName: justiceName,
		Decision:    des,
	}, nil
}
