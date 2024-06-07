package main

import (
	"log"
	"strings"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"
)

type Component struct {
	Name    string
	Version string
	URL     string
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
	opts := mdhtml.RendererOptions{Flags: htmlFlags}
	renderer := mdhtml.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func parseHTMLTable(htmlStr, startElement, id string) []Component {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	var chartVersions []Component
	var parseTable bool

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == startElement {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					parseTable = true
				}
			}
		}

		if parseTable && n.Type == html.ElementNode && n.Data == "tr" {
			var chartVersion Component
			tdCount := 0
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					if tdCount == 0 {
						chartVersion.Name = c.FirstChild.Data
					} else if tdCount == 1 {
						chartVersion.Version = c.FirstChild.FirstChild.Data
						chartVersion.URL = c.FirstChild.Attr[0].Val
					}
					tdCount++
				}
			}
			if tdCount > 0 {
				chartVersions = append(chartVersions, chartVersion)
			}
		}

		if n.Type == html.ElementNode && n.Data == "table" && parseTable && len(chartVersions) > 0 {
			parseTable = false
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return chartVersions
}
