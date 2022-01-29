package model

import (
	"golang.org/x/net/html"
)

type Field interface {
	Value(node *html.Node) string
}

type ExprField struct {
	Expr string `json:"expr"`
	Attr string `json:"attr"`
}

type TimeField struct {
	*ExprField
	Fmt string `json:"fmt"`
}

type ContentField struct {
	*ExprField
	NThread int      `json:"n_thread"`
	NTop    int      `json:"n_top"`
	NDelay  int      `json:"n_delay"`
	Blocks  []string `json:"blocks"`
}
