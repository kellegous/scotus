package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/kellegous/scotus/pkg/data"
	"github.com/kellegous/scotus/pkg/data/scotusdb"
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
	fs.StringVar(
		&f.ScotusDB.CasesURL,
		"scotusdb.cases-url",
		scotusdb.DefaultCasesURL,
		"the URL to download scotusdb case-centered data")
	fs.BoolVar(
		&f.ResetData,
		"reset-data",
		false,
		"whether to reset the data dir")
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

	terms, err := scotusdb.Read(
		context.Background(),
		scotusdb.WithCasesFromURL(flags.ScotusDB.CasesURL),
		scotusdb.WithDataDir(flags.DataDir))
	if err != nil {
		log.Panic(err)
	}

	for _, term := range terms {
		fmt.Printf("%d\n", term.Year)
	}
}
