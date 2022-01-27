package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"os"

	d "go2rss/doc"

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

type Config struct {
	Proxy     string  `json:"proxy"`
	UserAgent string  `json:"user_agent"`
	Feeds     []*Feed `json:"feeds"`
}

func (c *Config) UpdateToFeeds() {
	for _, f := range c.Feeds {
		if f.UserAgent == "" {
			f.UserAgent = c.UserAgent
		}

		if f.Proxy == "" {
			f.Proxy = c.Proxy
		}
		if f.Proxy == "direct" {
			f.Proxy = ""
		}

		if f.FeedTitleExpr == "" {
			f.FeedTitleExpr = "/html/head/title"
		}
		if f.FeedDescriptionExpr == "" {
			f.FeedTitleExpr = "/html/head/media[@name=\"description\"]"
		}

		if f.Domain == "" {
			link, err := url.Parse(f.Feed)
			if err == nil {
				f.Domain = fmt.Sprintf("%s://%s", link.Scheme, link.Host)
			}
		}

	}
}

func (feed *Feed) FetchRss() (string, error) {
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

func Parse(config string) (*Config, error) {
	if config == "" {
		config = "config.json"
	}

	text, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	c := Config{}
	if err := json.Unmarshal(text, &c); err != nil {
		return nil, err
	}

	c.UpdateToFeeds()

	return &c, nil
}
