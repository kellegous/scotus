package scotusdb

import (
	"encoding/json"
	"fmt"
)

type Decision byte

func (d Decision) MarshalJSON() ([]byte, error) {
	switch d {
	case Abstained:
		return []byte(`"x"`), nil
	case WithMajority:
		return []byte(`"+"`), nil
	case AgainstMajority:
		return []byte(`"-"`), nil
	}
	return nil, fmt.Errorf("invalid decision: %c", d)
}

func (d *Decision) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "x":
		*d = Abstained
		return nil
	case "+":
		*d = WithMajority
		return nil
	case "-":
		*d = AgainstMajority
		return nil
	}

	return fmt.Errorf("invalid decision: %s", s)
}

const (
	Abstained       Decision = 'x'
	AgainstMajority Decision = '-'
	WithMajority    Decision = '+'
)

func decisionFromString(v string) (Decision, error) {
	switch v {
	case "1":
		return AgainstMajority, nil
	case "2":
		return WithMajority, nil
	case "":
		return Abstained, nil
	}
	return Abstained, fmt.Errorf("invalid decision \"%s\"", v)
}
