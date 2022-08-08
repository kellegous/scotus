package search

import "golang.org/x/net/html"

func Query(
	root *html.Node,
	selector Selector,
) []*html.Node {
	return query(root, selector, nil)
}

func query(
	root *html.Node,
	selector Selector,
	results []*html.Node,
) []*html.Node {
	if selector(root) {
		results = append(results, root)
	}

	for n := root.FirstChild; n != nil; n = n.NextSibling {
		results = query(n, selector, results)
	}

	return results
}
