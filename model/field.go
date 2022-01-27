package model

import (
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Field interface {
	Value(node *html.Node) string
}

type ExprField struct {
	Expr string `json:"expr"`
	Attr string `json:"attr"`
}

func (f *ExprField) Value(node *html.Node) string {

	if !strings.HasPrefix(f.Expr, "/") {
		return f.Expr
	}

	if f.Attr == "" {
		return Text(node, f.Expr)
	}

	return Attr(node, f.Expr, f.Attr)
}

func (f *ExprField) Href(node *html.Node, domain string) string {
	v := f.Value(node)
	return domain + v
}

type TimeField struct {
	*ExprField
	Fmt string `json:"fmt"`
}

func (f *TimeField) Value(node *html.Node) time.Time {
	value := f.ExprField.Value(node)
	if value == "" {
		return time.Now()
	}

	t, err := time.Parse(f.Fmt, value)
	if err != nil {
		return time.Now()
	}

	return t
}
