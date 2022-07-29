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
	"golang.org/x/net/html"
)

const dataFileName = "overrulings.html"

var yearPattern = regexp.MustCompile(`\((\d{4})\)`)

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

func isElementOf(name string) func(n *html.Node) bool {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode &&
			n.Data == name
	}
}

func hasParentOf(
	n *html.Node,
	fn func(n *html.Node) bool,
) bool {
	for c := n.Parent; c != nil; c = c.Parent {
		if fn(c) {
			return true
		}
	}
	return false
}

func findElementsByPath(
	root *html.Node,
	path ...string,
) []*html.Node {
	if len(path) == 0 {
		return nil
	}

	// find all elements matching the right-most selector
	results := findAll(
		root,
		isElementOf(path[len(path)-1]),
		nil)
	path = path[:len(path)-1]

	// now walk backwards along the path, removing elements that
	// do not have matching parents.
	for n := len(path); n > 0; n = len(path) {
		var filtered []*html.Node
		q := isElementOf(path[n-1])
		path = path[:n-1]
		for _, result := range results {
			if hasParentOf(result, q) {
				filtered = append(filtered, result)
			}
		}
		results = filtered
	}

	return results
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
	trs := findElementsByPath(
		doc,
		"table", "tbody", "tr")

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
