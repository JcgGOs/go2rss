package model

import (
	d "go2rss/doc"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

type Feed struct {
	Name                string `json:"name"`
	Feed                string `json:"feed"`
	Proxy               string `json:"proxy"`
	ItemsExpr           string `json:"items_expr"`
	FeedTitleExpr       string `json:"feed_title_expr"`
	FeedDescriptionExpr string `json:"feed_description_expr"`
	TitleExpr           string `json:"title_expr"`
	DescriptionExpr     string `json:"description_expr"`
	LinkExpr            string `json:"link_expr"`
	AuthorExpr          string `json:"author_expr"`
	CreatedExpr         string `json:"created_expr"`
	UserAgent           string `json:"user_agent"`
	Domain              string
}

func (feed *Feed) Gen() (string, error) {
	resp, err := http.Get(feed.Feed)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	now := time.Now()
	doc, _ := d.Doc(resp.Body)
	nFeed := &feeds.Feed{
		Title:       d.Text(doc, feed.FeedTitleExpr),
		Link:        &feeds.Link{Href: feed.Feed},
		Description: d.Attr(doc, feed.FeedDescriptionExpr, "content"),
		Author:      &feeds.Author{Name: "wangyin", Email: "jmoiron@jmoiron.net"},
		Created:     now,
		Items:       []*feeds.Item{},
	}
	nodes := d.Nodes(doc, feed.ItemsExpr)
	for _, n := range nodes {
		item := feeds.Item{
			Title:       d.Text(n, feed.TitleExpr),
			Link:        &feeds.Link{Href: d.FullHref(n, feed.LinkExpr, feed.Domain)},
			Description: d.Text(n, feed.DescriptionExpr),
			Created:     now,
		}
		nFeed.Items = append(nFeed.Items, &item)
	}
	return nFeed.ToRss()
}
