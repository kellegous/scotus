package segalcover

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kellegous/scotus/pkg/data/internal"
	"github.com/kellegous/scotus/pkg/data/option"
	"github.com/kellegous/scotus/pkg/html/search"
	"golang.org/x/net/html"
)

const (
	dataFielName = "segal-cover.html"
	DefaultURL   = `https://en.wikipedia.org/wiki/Segal%E2%80%93Cover_score`
)

func Read(
	ctx context.Context,
	opts ...option.DownloadOption,
) ([]*Justice, error) {
	var o option.DownloadOptions
	o.ApplyOptions(opts, option.FromURL(DefaultURL))

	src := filepath.Join(o.DataDir, dataFielName)

	if err := internal.EnsureDownload(
		ctx,
		o.Client,
		o.URL,
		src,
	); err != nil {
		return nil, err
	}

	r, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return read(r)
}

func getInnerText(n *html.Node) string {
	var buf bytes.Buffer
	collectInnerText(n, &buf)
	return buf.String()
}

func collectInnerText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectInnerText(c, buf)
	}
}

func parseNominator(n *html.Node) (*President, error) {
	a := n.FirstChild
	if a == nil || a.Type != html.ElementNode || a.Data != "a" {
		return nil, fmt.Errorf("first child is not a element")
	}
	name := strings.TrimSpace(getInnerText(a))

	b := a.NextSibling
	if b == nil || b.Type != html.TextNode {
		return nil, fmt.Errorf("second child is not text")
	}

	party, err := partyFromString(
		strings.TrimSpace(getInnerText(b)))
	if err != nil {
		return nil, err
	}

	return &President{
		Name:  name,
		Party: party,
	}, nil
}

func parseJustice(tds []*html.Node) (*Justice, error) {
	if n := len(tds); n != 8 {
		return nil, fmt.Errorf("expected 8 columns, but found %d", n)
	}

	asChief := strings.TrimSpace(getInnerText(tds[2])) == "CJ"

	is := strings.TrimSpace(getInnerText(tds[4]))
	ideology, err := strconv.ParseFloat(is, 64)
	if err != nil {
		return nil, fmt.Errorf("ideology score: %w", err)
	}

	nominator, err := parseNominator(tds[6])
	if err != nil {
		return nil, fmt.Errorf("nominator: %w", err)
	}

	ys := strings.TrimSpace(getInnerText(tds[7]))
	year, err := strconv.Atoi(ys)
	if err != nil {
		return nil, fmt.Errorf("year: %w", err)
	}

	return &Justice{
		Name:          strings.TrimSpace(getInnerText(tds[1])),
		Chief:         asChief,
		Ideology:      ideology,
		YearNominated: year,
		NominatedBy:   nominator,
	}, nil
}

func read(r io.Reader) ([]*Justice, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	trs := search.Query(
		doc,
		search.HasAll(
			search.HasAll(
				search.IsElementOf("tr"),
				search.HasParent(
					search.HasAll(
						search.IsElementOf("tbody"),
						search.HasParent(
							search.HasID("segalcover"),
						),
					),
				),
			),
		),
	)

	var justices []*Justice
	for _, tr := range trs[1:] {
		justice, err := parseJustice(
			search.Query(tr, search.IsElementOf("td")))
		if err != nil {
			return nil, err
		}
		justices = append(justices, justice)
	}

	return justices, nil
}
