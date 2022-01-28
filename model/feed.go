package model

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"go2rss/util"
	"time"

	"github.com/gorilla/feeds"
)

type Feed struct {
	Domain          string
	Name            string     `json:"name"`
	Feed            string     `json:"feed"`
	Proxy           string     `json:"proxy"`
	Items           *ExprField `json:"items"`
	FeedTitle       *ExprField `json:"feed_title"`
	FeedDescription *ExprField `json:"feed_description"`
	Title           *ExprField `json:"title"`
	Description     *ExprField `json:"description"`
	Link            *ExprField `json:"link"`
	Author          *ExprField `json:"author"`
	Email           *ExprField `json:"email"`
	Created         *TimeField `json:"created"`
	Headers         []string   `json:"headers"`
}

func (feed *Feed) Gen() (string, error) {
	// util.GET(feed, "")
	body, err := util.GET(feed.Proxy, feed.Feed, feed.Headers)
	if err != nil {
		return "", err
	}
	// resp, err := http.Get(feed.Feed)
	// if err != nil {
	// 	return "", err
	// }
	// defer resp.Body.Close()

	now := time.Now()
	doc, _ := Load(bytes.NewReader(body))
	nFeed := &feeds.Feed{
		Title:       feed.FeedTitle.Value(doc),
		Link:        &feeds.Link{Href: feed.Feed},
		Description: feed.FeedDescription.Value(doc),
		Author:      &feeds.Author{Name: feed.Author.Value(doc), Email: feed.Email.Value(doc)},
		Created:     now,
		Items:       []*feeds.Item{},
	}

	for _, n := range Nodes(doc, feed.Items.Expr) {
		title := feed.Title.Value(n)
		_link := feed.Link.Href(n, feed.Domain)
		nFeed.Items = append(nFeed.Items, &feeds.Item{
			Title:       feed.Title.Value(n),
			Link:        &feeds.Link{Href: _link},
			Description: feed.Description.Value(n),
			Created:     feed.Created.Value(n),
			Id:          fmt.Sprintf("%x", md5.Sum([]byte(title+_link))),
		})
	}
	return nFeed.ToRss()
}
