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

	"github.com/kellegous/scotus/pkg/data/scotusdb"
)

var justices = []string{
	"SSotomayor",
	"EKagan",
	"SGBreyer",
	"JGRoberts",
	"BMKavanaugh",
	"ACBarrett",
	"NMGorsuch",
	"SAAlito",
	"CThomas",
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

func parseDecision(s string) (scotusdb.Decision, error) {
	switch s {
	case "":
		return scotusdb.AgainstMajority, nil
	case "1":
		return scotusdb.WithMajority, nil
	case "2":
		return scotusdb.Abstained, nil
	}
	return scotusdb.Abstained, fmt.Errorf("invalid decision: %s", s)
}

func parseCase(
	id string,
	row []string,
) (*scotusdb.Case, error) {
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

	var votes []*scotusdb.Vote

	cols := row[5:]
	for i, justice := range justices {
		d, err := parseDecision(cols[i])
		if err != nil {
			return nil, err
		}

		votes = append(votes, &scotusdb.Vote{
			ID:          fmt.Sprintf("%s-%s", id, justice),
			JusticeName: justice,
			Decision:    d,
		})
	}

	return &scotusdb.Case{
		ID:            id,
		Name:          name,
		DecisionDate:  date,
		MajorityVotes: maj,
		MinorityVotes: min,
		Votes:         votes,
	}, nil
}

func readCases(r io.Reader) ([]*scotusdb.Case, error) {
	cr := csv.NewReader(r)
	cr.Comma = '\t'

	if _, err := cr.Read(); err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, err
	}

	var cases []*scotusdb.Case
	for i := 1; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		c, err := parseCase(fmt.Sprintf("2021-%03d", i), row)
		if err != nil {
			return nil, err
		}

		cases = append(cases, c)
	}

	return cases, nil
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

	cases, err := readCases(r)
	if err != nil {
		log.Panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(cases); err != nil {
		log.Panic(err)
	}
}
