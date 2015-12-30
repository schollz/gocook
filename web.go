package main

import (
	"bytes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io/ioutil"
	"net/http"
	"strings"
)

var badTags = []atom.Atom{
	atom.Meta,
	atom.A,
	atom.Img,
	atom.Figcaption,
}

func parseURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	HTML := string(body)
	root, err := html.Parse(bytes.NewReader([]byte(HTML)))
	if err != nil {
		panic(err)
	}

	text := ""
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Meta {
			return
		}
		if n.Type == html.TextNode {
			for _, l := range strings.Split(n.Data, "\n\n") {
				if s := strings.TrimSpace(l); s != "" {
					text += strings.Replace(s, "\n", " ", -1) + "\n"
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)

	return text
}
