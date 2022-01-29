package util

import (
	"fmt"
	"go2rss/model"
	"io"
	"net/url"
	"strings"
	"time"

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

func Nodes(doc *html.Node, expr string) []*html.Node {
	return htmlquery.Find(doc, expr)
}

func Load(r io.Reader) (*html.Node, error) {
	return htmlquery.Parse(r)
}

func Value(node *html.Node, f *model.ExprField) string {

	if !strings.HasPrefix(f.Expr, "/") {
		return f.Expr
	}

	if f.Attr == "" {
		return Text(node, f.Expr)
	}

	return Attr(node, f.Expr, f.Attr)
}

func AbsHref(v string, uri string) string {
	if strings.HasPrefix(v, "http") {
		return v
	}
	if strings.HasPrefix(v, "//") {
		idx := strings.Index(uri, "//")
		return uri[:idx] + v
	}

	idx := strings.LastIndex(uri, "/")
	if idx < len(uri)-1 {
		uri = uri[:idx+1]
	}

	if strings.HasPrefix(v, "/") {
		_url, _ := url.Parse(uri)
		return fmt.Sprintf("%s://%s%s", _url.Scheme, _url.Host, v)
	}

	return uri + v
}

func CDATA(node *html.Node, f *model.ExprField) string {
	v := Value(node, f)
	return fmt.Sprintf("<![CDATA[ %s ]]>", v)
}

func TimeValue(node *html.Node, f *model.TimeField) time.Time {
	value := Value(node, f.ExprField)
	if value == "" {
		return time.Now()
	}

	t, err := time.Parse(f.Fmt, value)
	if err != nil {
		return time.Now()
	}

	return t
}

func Html(node *html.Node, f *model.ContentField) string {
	n := htmlquery.FindOne(node, f.Expr)
	if n == nil {
		return ""
	}
	if f.Attr == "html" {
		return htmlquery.OutputHTML(n, false)
	}
	return Value(node, f.ExprField)
}
