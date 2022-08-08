package csv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Row struct {
	fields map[string]int
	values []string
}

func (r *Row) Get(name string) (string, error) {
	if i, ok := r.fields[name]; ok {
		return r.values[i], nil
	}

	return "", fmt.Errorf("unknown field: %s", name)
}

func (r *Row) GetInt(
	name string,
	emptyVal int,
) (int, error) {
	v, err := r.Get(name)
	if err != nil {
		return 0, err
	}

	if v == "" {
		return emptyVal, nil
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s is not a valid int (%s)", name, v)
	}

	return i, nil
}

func (r *Row) GetDate(name string) (time.Time, error) {
	v, err := r.Get(name)
	if err != nil {
		return time.Time{}, err
	}

	t, err := parseDate(v)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s is not a valid date", name)
	}

	return t, nil
}

func (r *Row) AsMap() map[string]string {
	row := map[string]string{}
	for f, ix := range r.fields {
		row[f] = r.values[ix]
	}
	return row
}

func parseDate(v string) (time.Time, error) {
	vals := strings.SplitN(v, "/", 3)
	if len(vals) != 3 {
		return time.Time{}, fmt.Errorf("invalid date: %s", v)
	}

	m, err := strconv.Atoi(vals[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %s", v)
	}

	d, err := strconv.Atoi(vals[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %s", v)
	}

	y, err := strconv.Atoi(vals[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %s", v)
	}

	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
}

func (r *Row) MustMarshal() []byte {
	m := map[string]string{}
	for field, idx := range r.fields {
		m[field] = r.values[idx]
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}

	return b
}
