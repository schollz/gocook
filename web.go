package main

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var badTags = []atom.Atom{
	atom.Meta,
	atom.A,
	atom.Img,
	atom.Figcaption,
}

func parseURL(url string) (text string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err) // TODO: handle errors better?
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err) // TODO: handle errors better?
	}

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Meta {
			return
		}
		if n.Type == html.TextNode {
			for _, l := range strings.Split(n.Data, "\n\n") {
				if s := strings.TrimSpace(l); s != "" {
					if n.Parent != nil && n.Parent.Type == html.ElementNode && n.DataAtom == atom.Title {
						s = "__TITLE__=" + s + "__TITLE__"
					}
					text += strings.Replace(s, "\n", " ", -1) + "\n"
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)

	return
}
