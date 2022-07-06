package scotus

import "fmt"

type Decision byte

const (
	Abstained Decision = iota
	AgainstMajority
	WithMajority
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
