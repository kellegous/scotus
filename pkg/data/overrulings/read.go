package overrulings

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/kellegous/scotus/pkg/data/internal"
	"github.com/kellegous/scotus/pkg/data/option"
	"github.com/kellegous/scotus/pkg/html/search"
	"golang.org/x/net/html"
)

const (
	dataFileName = "overrulings.html"
	DefaultURL   = "https://constitution.congress.gov/resources/decisions-overruled/"
)

var yearPattern = regexp.MustCompile(`\((\d{4})\)`)

func Read(
	ctx context.Context,
	opts ...option.DownloadOption,
) ([]*Decision, error) {
	var o option.DownloadOptions
	o.ApplyOptions(opts, option.FromURL(DefaultURL))

	src := filepath.Join(o.DataDir, dataFileName)

	if err := internal.EnsureDownload(
		ctx,
		o.Client,
		o.URL,
		src,
	); err != nil {
		return nil, err
	}

	return read(src)
}

func parseYear(td *html.Node) (int, error) {
	if c := td.FirstChild; c != nil && c.Type == html.TextNode {
		return strconv.Atoi(c.Data)
	}

	return 0, errors.New("cell doesn't include a single text node")
}

func innerTextOf(n *html.Node) string {
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			buf.WriteString(c.Data)
		} else if c.Type == html.ElementNode && c.Data == "b" {
			buf.WriteByte(' ')
		}
	}
	return buf.String()
}

func getAttribute(
	n *html.Node,
	name string,
) string {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func parseCaseText(n *html.Node) (string, *html.Node) {
	var buf bytes.Buffer
	for ; n != nil; n = n.NextSibling {
		switch n.Type {
		case html.TextNode:
			d := strings.TrimSpace(n.Data)
			if strings.HasSuffix(d, ";") {
				buf.WriteString(d[:len(d)-1])
				return strings.TrimSpace(buf.String()), n
			}
			buf.WriteString(d)
		case html.ElementNode:
			buf.WriteByte(' ')
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

func extractYear(name string) (int, error) {
	m := yearPattern.FindStringSubmatch(name)
	if len(m) != 2 {
		return 0, fmt.Errorf("year pattern not found in <%s>", name)
	}
	return strconv.Atoi(m[1])
}

func parseCases(td *html.Node) ([]*Case, error) {
	var cases []*Case
	for c := td.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			cases = append(cases, &Case{
				Name: innerTextOf(c),
				URL:  getAttribute(c, "href"),
			})
		} else {
			s, n := parseCaseText(c)
			if s != "" {
				cases = append(cases, &Case{
					Name: s,
				})
			}
			if n == nil {
				break
			}
			c = n
		}
	}
	return cases, nil
}

func parseDecision(tds []*html.Node) (*Decision, error) {
	c, err := parseCases(tds[1])
	if err != nil {
		return nil, err
	} else if len(c) != 1 {
		return nil, fmt.Errorf("expected a single case, found %d", len(c))
	}

	c[0].Year, err = parseYear(tds[2])
	if err != nil {
		return nil, err
	}

	cases, err := parseCases(tds[3])
	if err != nil {
		return nil, err
	}

	for _, c := range cases {
		year, err := extractYear(c.Name)
		if err != nil {
			return nil, err
		}
		c.Year = year
	}

	return &Decision{
		Case:      c[0],
		Overruled: cases,
	}, nil
}

func extract(doc *html.Node) ([]*Decision, error) {
	// selector "table tbody tr"
	trs := search.Query(
		doc,
		search.HasAll(
			search.IsElementOf("tr"),
			search.HasParent(
				search.HasAll(
					search.IsElementOf("tbody"),
					search.HasParent(
						search.IsElementOf("table"),
					),
				),
			),
		),
	)

	decisions := make([]*Decision, 0, len(trs))

	for _, tr := range trs {
		tds := search.Query(tr, search.IsElementOf("td"))
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
