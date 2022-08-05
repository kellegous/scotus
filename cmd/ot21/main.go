package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var justices = []string{
	"Sotomayor",
	"Kagan",
	"Breyer",
	"Roberts",
	"Kavanaugh",
	"Barrett",
	"Gorsuch",
	"Alito",
	"Thomas",
}

type Vote byte

func (v Vote) MarshalJSON() ([]byte, error) {
	switch v {
	case WithMajority:
		return []byte(`"+"`), nil
	case WithMinority:
		return []byte(`"-"`), nil
	case Recused:
		return []byte(`"x"`), nil
	}
	return nil, fmt.Errorf("invalid vote: %c", v)
}

const (
	WithMajority Vote = '+'
	WithMinority Vote = '-'
	Recused      Vote = 'x'
)

type Decision struct {
	Name     string          `json:"name"`
	Date     time.Time       `json:"day"`
	Majority int             `json:"majority"`
	Minority int             `json:"minority"`
	Author   string          `json:"author"`
	Votes    map[string]Vote `json:"votes,omitempty"`
}

func (d *Decision) isValid() bool {
	if d.Votes == nil {
		return true
	}

	var maj, min int
	for _, vote := range d.Votes {
		switch vote {
		case WithMajority:
			maj++
		case WithMinority:
			min++
		}
	}

	return maj == d.Majority && min == d.Minority
}

func parseDecision(row []string) (*Decision, error) {
	if len(row) != 14 {
		return nil, fmt.Errorf("row should have 14 columns but had %d instead", len(row))
	}

	name := strings.TrimSpace(row[0])

	date, err := time.ParseInLocation("2006-01-02", row[1], time.UTC)
	if err != nil {
		return nil, fmt.Errorf("date: %w", err)
	}

	maj, err := strconv.Atoi(row[2])
	if err != nil {
		return nil, fmt.Errorf("majority: %w", err)
	}

	min, err := strconv.Atoi(row[3])
	if err != nil {
		return nil, fmt.Errorf("minority: %w", err)
	}

	var votes map[string]Vote
	// so here's the thing comrades, LeDure was a per curium opinion where the votes were not
	// known. It is reported as 4-4 but we don't know how people voted ... except that Barrett
	// recused herself.
	if name != "LeDure" {
		votes, err = parseVotes(row[5:])
		if err != nil {
			return nil, fmt.Errorf("votes: %w", err)
		}
	}

	return &Decision{
		Name:     name,
		Date:     date,
		Majority: maj,
		Minority: min,
		Author:   strings.TrimSpace(row[4]),
		Votes:    votes,
	}, nil
}

func parseVote(s string) (Vote, error) {
	switch strings.ToLower(s) {
	case "":
		return WithMinority, nil
	case "1":
		return WithMajority, nil
	case "2":
		return Recused, nil
	}
	return Vote('?'), fmt.Errorf("unknown vote: %s", s)
}

func parseVotes(cols []string) (map[string]Vote, error) {
	votes := map[string]Vote{}
	for i, justice := range justices {
		vote, err := parseVote(cols[i])
		if err != nil {
			return nil, err
		}
		votes[justice] = vote
	}
	return votes, nil
}

type Flags struct {
	Src string
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.Src,
		"src",
		"ot21.tsv",
		"the source file")
}

func readDecisions(r io.Reader) ([]*Decision, error) {
	cr := csv.NewReader(r)
	cr.Comma = '\t'

	// discard the header
	if _, err := cr.Read(); err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, err
	}

	var decisions []*Decision
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		decision, err := parseDecision(row)
		if err != nil {
			return nil, err
		}
		decisions = append(decisions, decision)
	}

	return decisions, nil
}

func verifyAll(decisions []*Decision) error {
	for _, decision := range decisions {
		if !decision.isValid() {
			return fmt.Errorf("%s is invalid", decision.Name)
		}
	}
	return nil
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	r, err := os.Open(flags.Src)
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	decisions, err := readDecisions(r)
	if err != nil {
		log.Panic(err)
	}

	if err := verifyAll(decisions); err != nil {
		log.Panic(err)
	}

	b, err := json.MarshalIndent(decisions, "", "  ")
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%s\n", b)
}
