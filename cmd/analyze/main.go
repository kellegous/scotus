package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/kellegous/scotus/pkg/scotus"
)

type Flags struct {
	StoreDir     string
	CasesFromURL string
	ResetStore   bool
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.StoreDir,
		"store.dir",
		"data",
		"the directory where downloaded data will be kept")
	fs.StringVar(
		&f.CasesFromURL,
		"store.cases-url",
		scotus.DefaultCasesURL,
		"the SCOTUS database file to use for case data")

	fs.BoolVar(
		&f.ResetStore,
		"store.reset",
		false,
		"whether to reset the store")
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	s, err := scotus.OpenStore(
		context.Background(),
		flags.StoreDir,
		scotus.WithCasesFromURL(flags.CasesFromURL),
		scotus.Reset(flags.ResetStore))
	if err != nil {
		log.Panic(err)
	}

	terms, err := s.Terms()
	if err != nil {
		log.Panic(err)
	}

	for _, term := range terms {
		fmt.Printf("%d (%d)\n", term.Year, len(term.Cases))
	}
}
