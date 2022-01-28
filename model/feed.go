package model

import (
	"bytes"
	"go2rss/util"
	"math/rand"
	"sync"
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
	Content         *ExprField `json:"content"`
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

	for i, n := range Nodes(doc, feed.Items.Expr) {
		if i > 10 {
			break
		}
		_link := feed.Link.Href(n, feed.Domain)
		_item := &feeds.Item{
			Title:       feed.Title.Value(n),
			Link:        &feeds.Link{Href: _link},
			Description: feed.Description.Value(n),
			Created:     feed.Created.Value(n),
			Id:          _link,
		}
		nFeed.Items = append(nFeed.Items, _item)
	}

	if feed.Content != nil && len(nFeed.Items) > 0 {
		var wg sync.WaitGroup
		ch := make(chan *feeds.Item, len(nFeed.Items))
		for _, f := range nFeed.Items {
			ch <- f
			wg.Add(1)
			go func(_f *feeds.Item) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
				itemBody, err := util.GET(feed.Proxy, _f.Link.Href, feed.Headers)
				_itemDoc, _ := Load(bytes.NewReader(itemBody))
				if err == nil {
					_f.Content = feed.Content.Value(_itemDoc)
				}
				<-ch
			}(f)
		}
		wg.Wait()
	}
	return nFeed.ToRss()
}
