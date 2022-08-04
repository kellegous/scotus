package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	votePattern = regexp.MustCompile(` ([0-9])-([0-4]).? `)
	datePattern = regexp.MustCompile(` ([A-Za-z\.]+) (\d{1,2}), (\d{4})`)
)

var months = map[string]time.Month{
	"January":   time.January,
	"Jan.":      time.January,
	"February":  time.February,
	"Feb.":      time.February,
	"March":     time.March,
	"Mar.":      time.March,
	"April":     time.April,
	"Apr.":      time.April,
	"May":       time.May,
	"June":      time.June,
	"Jun.":      time.June,
	"July":      time.July,
	"Jul.":      time.July,
	"August":    time.August,
	"Aug.":      time.August,
	"September": time.September,
	"Sep.":      time.September,
	"October":   time.October,
	"Oct.":      time.October,
	"November":  time.November,
	"Nov.":      time.November,
	"December":  time.December,
	"Dec.":      time.December,
}

type Decision struct {
	Name   string
	Day    time.Time
	Maj    int
	Min    int
	Author string
}

func parseDateAndName(v []byte) (string, time.Time, error) {
	idx := datePattern.FindSubmatchIndex(v)
	if len(idx) == 0 {
		return "", time.Time{}, fmt.Errorf("could not find date in %s", v)
	}

	m := months[string(v[idx[2]:idx[3]])]
	if m == 0 {
		return "", time.Time{}, fmt.Errorf("invalid month: %s", v)
	}

	d, err := strconv.Atoi(string(v[idx[4]:idx[5]]))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("day: %w", err)
	}

	y, err := strconv.Atoi(string(v[idx[6]:idx[7]]))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("year: %w", err)
	}

	return string(v[:idx[0]]), time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}

func parseDecision(v []byte) (*Decision, error) {
	idx := votePattern.FindSubmatchIndex(v)
	if idx == nil {
		return nil, fmt.Errorf("could not find vote pattern in %s", v)
	}

	maj, err := strconv.Atoi(string(v[idx[2]:idx[3]]))
	if err != nil {
		return nil, fmt.Errorf("majority: %w", err)
	}

	min, err := strconv.Atoi(string(v[idx[4]:idx[5]]))
	if err != nil {
		return nil, fmt.Errorf("minority: %w", err)
	}

	prefix := bytes.TrimSpace(v[:idx[0]])
	name, day, err := parseDateAndName(prefix)
	if err != nil {
		return nil, err
	}

	suffix := bytes.TrimSpace(v[idx[1]:])
	return &Decision{
		Name:   name,
		Day:    day,
		Maj:    maj,
		Min:    min,
		Author: string(suffix),
	}, nil
}

type Flags struct {
	Src string
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.Src,
		"src",
		"ot21.txt",
		"the source file")
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

	var decisions []*Decision
	s := bufio.NewScanner(r)
	for s.Scan() {
		data := s.Bytes()
		decision, err := parseDecision(data)
		if err != nil {
			log.Panic(err)
		}
		decisions = append(decisions, decision)
	}

	if err := s.Err(); err != nil {
		log.Panic(err)
	}

	for _, decision := range decisions {
		row := []string{
			decision.Name,
			decision.Day.Format("2006-01-02"),
			strconv.Itoa(decision.Maj),
			strconv.Itoa(decision.Min),
			decision.Author,
		}
		fmt.Printf("%s\n", strings.Join(row, "\t"))
	}
}
