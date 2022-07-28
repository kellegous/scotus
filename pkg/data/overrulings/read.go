package overrulings

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kellegous/scotus/pkg/data/internal"
	"golang.org/x/net/html"
)

const dataFileName = "overrulings.html"

func Read(
	ctx context.Context,
	opts ...Option,
) ([]*Decision, error) {
	var o Options
	o.apply(opts)

	src := filepath.Join(o.dataDir, dataFileName)

	if err := internal.EnsureDownload(
		ctx,
		o.client,
		o.url,
		src,
	); err != nil {
		return nil, err
	}

	return read(src)
}

func findOne(
	root *html.Node,
	matches func(n *html.Node) bool,
) *html.Node {
	if matches(root) {
		return root
	}

	for n := root.FirstChild; n != nil; n = n.NextSibling {
		if n := findOne(n, matches); n != nil {
			return n
		}
	}

	return nil
}

func findAll(
	root *html.Node,
	matches func(n *html.Node) bool,
	existing []*html.Node,
) []*html.Node {
	if matches(root) {
		return append(existing, root)
	}

	for n := root.FirstChild; n != nil; n = n.NextSibling {
		existing = findAll(n, matches, existing)
	}

	return existing
}

func parseYear(td *html.Node) (int, error) {
	if c := td.FirstChild; c != nil && c.Type == html.TextNode {
		return strconv.Atoi(c.Data)
	}

	return 0, errors.New("cell doesn't include a single text node")
}

func parseYears(td *html.Node) ([]int, error) {
	var years []int
	for c := td.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.TextNode {
			continue
		}

		year, err := strconv.Atoi(c.Data)
		if err != nil {
			return nil, err
		}

		years = append(years, year)
	}

	return years, nil
}

func parseOverruledCases(td []*html.Node) ([]*Case, error) {
	years, err := parseYears(td[4])
	if err != nil {
		return nil, err
	}

	cases := make([]*Case, 0, len(years))
	for _, year := range years {
		cases = append(cases, &Case{
			Year: year,
		})
	}

	return cases, nil
}

func parseDecision(tds []*html.Node) (*Decision, error) {
	year, err := parseYear(tds[2])
	if err != nil {
		return nil, err
	}

	cases, err := parseOverruledCases(tds)
	if err != nil {
		return nil, err
	}

	return &Decision{
		Case: Case{
			Year: year,
		},
		Overruled: cases,
	}, nil
}

func extract(doc *html.Node) ([]*Decision, error) {
	table := findOne(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode &&
			n.Data == "table"
	})
	if table == nil {
		return nil, errors.New("no table found in html doc")
	}

	tbody := findOne(table, func(n *html.Node) bool {
		return n.Type == html.ElementNode &&
			n.Data == "tbody"
	})
	if tbody == nil {
		return nil, errors.New("no table body found in html doc")
	}

	trs := findAll(
		tbody,
		func(n *html.Node) bool {
			return n.Type == html.ElementNode &&
				n.Data == "tr"
		},
		nil)

	decisions := make([]*Decision, 0, len(trs))

	for _, tr := range trs {
		tds := findAll(
			tr,
			func(n *html.Node) bool {
				return n.Type == html.ElementNode &&
					n.Data == "td"
			},
			nil)
		decision, err := parseDecision(tds)
		if err != nil {
			return nil, err
		}

		decisions = append(decisions, decision)
	}

	return decisions, nil
}

func read(src string) ([]*Decision, error) {
	r, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return extract(doc)
}
