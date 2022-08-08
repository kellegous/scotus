package segalcover

import (
	"fmt"
	"strings"
)

type Party string

const (
	Republican Party = "R"
	Democrat   Party = "D"
	Unknown    Party = "?"
)

func partyFromString(s string) (Party, error) {
	switch strings.ToLower(s) {
	case "(republican)":
		return Republican, nil
	case "(democrat)":
		return Democrat, nil
	}
	return Unknown, fmt.Errorf("unknown party: %s", s)
}
