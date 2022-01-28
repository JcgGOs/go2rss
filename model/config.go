package model

import (
	"encoding/json"
	"fmt"
	"net/url"

	"os"
)

type Config struct {
	Proxy   string   `json:"proxy"`
	Headers []string `json:"headers"`
	Feeds   []*Feed  `json:"feeds"`
}

func (c *Config) Init() {
	for _, f := range c.Feeds {

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

		if f.Headers == nil {
			f.Headers = []string{}
		}

		if len(c.Headers) > 0 {
			f.Headers = append(f.Headers, c.Headers...)
		}

		if len(f.Headers) < 1 {
			f.Headers = append(f.Headers, "User-Agent:Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.69")
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

	c.Init()

	return &c, nil
}
