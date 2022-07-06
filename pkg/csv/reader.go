package csv

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Reader struct {
	r      *csv.Reader
	fields map[string]int
}

func NewReader(r io.Reader) (*Reader, error) {
	cr := csv.NewReader(r)

	hdrs, err := cr.Read()
	if err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, err
	}

	fields := map[string]int{}
	for i, hdr := range hdrs {
		fields[hdr] = i
	}

	return &Reader{
		r:      cr,
		fields: fields,
	}, nil
}

func (r *Reader) Next() (*Row, error) {
	vals, err := r.r.Read()
	if err != nil {
		return nil, err
	}

	if len(vals) != len(r.fields) {
		return nil, fmt.Errorf(
			"wrong number of columns in row, expected %d got %d",
			len(r.fields),
			len(vals))
	}

	return &Row{
		fields: r.fields,
		values: vals,
	}, nil
}
