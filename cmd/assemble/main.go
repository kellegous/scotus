package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/kellegous/scotus/pkg/data"
	"github.com/kellegous/scotus/pkg/data/option"
	"github.com/kellegous/scotus/pkg/data/overrulings"
	"github.com/kellegous/scotus/pkg/data/scotusdb"
	"github.com/kellegous/scotus/pkg/data/segalcover"
)

type Flags struct {
	DataDir   string
	ResetData bool
	ScotusDB  struct {
		CasesURL string
	}
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.DataDir,
		"data-dir",
		"data",
		"the directory where the data will be kept")
	fs.BoolVar(
		&f.ResetData,
		"reset-data",
		false,
		"whether to reset the data dir")
}

func distinctStringsFromCases(
	cases []*scotusdb.Case,
	fn func(c *scotusdb.Case) string,
) []string {
	seen := map[string]bool{}
	var items []string
	for _, c := range cases {
		item := fn(c)
		if seen[item] {
			continue
		}

		seen[item] = true
		items = append(items, item)
	}
	return items
}

func distinctStringsFromVotes(
	cases []*scotusdb.Case,
	fn func(v *scotusdb.Vote) string,
) []string {
	seen := map[string]bool{}
	var items []string
	for _, c := range cases {
		for _, v := range c.Votes {
			item := fn(v)
			if seen[item] {
				continue
			}
			seen[item] = true
			items = append(items, item)
		}
	}
	sort.Strings(items)
	return items
}

type JusticeDirection struct {
	Justice    string
	Directions map[scotusdb.Direction]int
}

func getAllDirections(terms []*scotusdb.Term) []*JusticeDirection {
	byJustice := map[string]*JusticeDirection{}
	for _, t := range terms {
		for _, c := range t.Cases {
			for _, v := range c.Votes {
				forJustice := byJustice[v.JusticeName]
				if forJustice == nil {
					forJustice = &JusticeDirection{
						Justice:    v.JusticeName,
						Directions: map[scotusdb.Direction]int{},
					}
					byJustice[v.JusticeName] = forJustice
				}
				forJustice.Directions[v.Direction]++
			}
		}
	}

	var directions []*JusticeDirection
	for _, direction := range byJustice {
		directions = append(directions, direction)
	}

	sort.Slice(directions, func(i, j int) bool {
		return directions[i].Justice < directions[j].Justice
	})

	return directions
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	if err := data.EnsureDir(
		flags.DataDir,
		0755,
		flags.ResetData,
	); err != nil {
		log.Panic(err)
	}

	_, err := scotusdb.Read(
		context.Background(),
		scotusdb.WithDataDir(flags.DataDir))
	if err != nil {
		log.Panic(err)
	}

	overruled, err := overrulings.Read(
		context.Background(),
		option.WithDataDir(flags.DataDir),
	)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%d\n", len(overruled))

	if _, err := segalcover.Read(
		context.Background(),
		option.WithDataDir(flags.DataDir),
	); err != nil {
		log.Panic(err)
	}
}
