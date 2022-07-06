package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

func toJSON(r io.Reader) ([]map[string]string, error) {
	cr := csv.NewReader(r)

	// read the fields from the header
	fields, err := cr.Read()
	if err != nil {
		return nil, err
	}

	rows := []map[string]string{}

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		obj := map[string]string{}

		for i, field := range fields {
			obj[field] = row[i]
		}

		rows = append(rows, obj)
	}

	return rows, nil
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Panic("wrong use")
	}

	r, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	data, err := toJSON(r)
	if err != nil {
		log.Panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
		log.Panic(err)
	}
}
