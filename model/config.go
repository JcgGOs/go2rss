package model

import (
	"encoding/json"
	"fmt"
	"net/url"

	"os"
)

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

		if f.FeedTitle.Expr == "" {
			f.FeedTitle.Expr = "/html/head/title"
		}
		if f.FeedDescription.Expr == "" {
			f.FeedDescription.Expr = "/html/head/media[@name=\"description\"]"
			f.FeedDescription.Attr = "content"
		}

		if f.Domain == "" {
			link, err := url.Parse(f.Feed)
			if err == nil {
				f.Domain = fmt.Sprintf("%s://%s", link.Scheme, link.Host)
			}
		}

	}
}

func (c *Config) GetFeed(name string) *Feed {
	for _, v := range c.Feeds {
		if v.Name == name {
			return v
		}
	}
	return nil
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
