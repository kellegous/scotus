package scotusdb

import (
	"fmt"
	"strings"

	"github.com/kellegous/scotus/pkg/csv"
)

type Direction string

const (
	Unknown      Direction = "?"
	Liberal      Direction = "L"
	Conservative Direction = "C"
)

type Vote struct {
	ID          string    `json:"id"`
	JusticeName string    `json:"justice-name"`
	Decision    Decision  `json:"decision"`
	Direction   Direction `json:"direction"`
}

func directionFromString(s string) (Direction, error) {
	switch strings.TrimSpace(s) {
	case "1":
		return Conservative, nil
	case "2":
		return Liberal, nil
	case "":
		return Unknown, nil
	}
	return Unknown, fmt.Errorf("invalid direction: %s", s)
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

	dir, err := row.Get("direction")
	if err != nil {
		return nil, err
	}

	direction, err := directionFromString(dir)
	if err != nil {
		return nil, err
	}

	return &Vote{
		ID:          id,
		JusticeName: justiceName,
		Decision:    des,
		Direction:   direction,
	}, nil
}
