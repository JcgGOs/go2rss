package model

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

type Feed struct {
	Domain          string
	Name            string     `json:"name"`
	Feed            string     `json:"feed"`
	Proxy           string     `json:"proxy"`
	UserAgent       string     `json:"user_agent"`
	Items           *ExprField `json:"items"`
	FeedTitle       *ExprField `json:"feed_title"`
	FeedDescription *ExprField `json:"feed_description"`
	Title           *ExprField `json:"title"`
	Description     *ExprField `json:"description"`
	Link            *ExprField `json:"link"`
	Author          *ExprField `json:"author"`
	Email           *ExprField `json:"email"`
	Created         *TimeField `json:"created"`
}

func (feed *Feed) Gen() (string, error) {
	resp, err := http.Get(feed.Feed)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	now := time.Now()
	doc, _ := Load(resp.Body)
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
