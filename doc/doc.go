package doc

import (
	"io"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func Text(node *html.Node, expr string) string {
	n := htmlquery.FindOne(node, expr)
	if n == nil {
		return ""
	}
	return htmlquery.InnerText(n)
}

func Attr(node *html.Node, expr string, attr string) string {
	n := htmlquery.FindOne(node, expr)
	if n == nil {
		return ""
	}
	return htmlquery.SelectAttr(n, attr)
}

func Href(node *html.Node, expr string) string {
	return Attr(node, expr, "href")
}

func FullHref(node *html.Node, expr string, domain string) string {
	href := Href(node, expr)
	if strings.HasPrefix(href, "http") {
		return href
	}
	return domain + href
}

func Nodes(doc *html.Node, expr string) []*html.Node {
	return htmlquery.Find(doc, expr)
}

func Doc(r io.Reader) (*html.Node, error) {
	return htmlquery.Parse(r)
}
