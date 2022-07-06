package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/kellegous/scotus/pkg/scotus"
)

type Flags struct {
	DataFile string
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.DataFile,
		"data-file",
		"SCDB_2021_01_justiceCentered_Citation.csv",
		"The justic centered data file from SCOTUS database")
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	terms, err := scotus.ReadFile(flags.DataFile)
	if err != nil {
		log.Panic(err)
	}

	b, err := json.MarshalIndent(terms, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%s\n", b)
}
