package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kellegous/scotus/pkg/csv"
)

func main() {
	r, err := os.Open("SCDB_2021_01_justiceCentered_Vote.csv")
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	cr, err := csv.NewReader(r)
	if err != nil {
		log.Panic(err)
	}

	for {
		row, err := cr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}

		b, err := json.MarshalIndent(row.AsMap(), "", "  ")
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("%s\n", b)
	}
}
