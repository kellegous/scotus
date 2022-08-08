package search

import (
	"regexp"

	"golang.org/x/net/html"
)

var whitespacePattern = regexp.MustCompile(`\s+`)

type Selector func(n *html.Node) bool

func IsElementOf(name string) Selector {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == name
	}
}

func HasClass(class string) Selector {
	return func(n *html.Node) bool {
		return NodeHasClass(n, class)
	}
}

func HasID(id string) Selector {
	return func(n *html.Node) bool {
		return NodeHasID(n, id)
	}
}

func HasParent(selector Selector) Selector {
	return func(n *html.Node) bool {
		for n := n.Parent; n != nil; n = n.Parent {
			if selector(n) {
				return true
			}
		}
		return false
	}
}

func HasDirectParent(selector Selector) Selector {
	return func(n *html.Node) bool {
		if p := n.Parent; p != nil {
			return selector(p)
		}
		return false
	}
}

func HasAll(first Selector, selectors ...Selector) Selector {
	return func(n *html.Node) bool {
		if !first(n) {
			return false
		}
		for _, sel := range selectors {
			if !sel(n) {
				return false
			}
		}
		return true
	}
}

func HasAny(first Selector, selectors ...Selector) Selector {
	return func(n *html.Node) bool {
		if first(n) {
			return true
		}
		for _, sel := range selectors {
			if sel(n) {
				return true
			}
		}
		return false
	}
}

func NodeHasClass(n *html.Node, c string) bool {
	attr := GetAttr(n, "class")
	if attr == nil {
		return false
	}
	for _, class := range whitespacePattern.Split(attr.Val, -1) {
		if class == c {
			return true
		}
	}
	return false
}

func NodeHasID(n *html.Node, id string) bool {
	if attr := GetAttr(n, "id"); attr != nil {
		return attr.Val == id
	}
	return false
}

func GetAttr(n *html.Node, name string) *html.Attribute {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return &attr
		}
	}
	return nil
}
