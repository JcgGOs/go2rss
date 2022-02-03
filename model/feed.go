package model

import "regexp"

var UA = "User-Agent:Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.69"

var Regexes = make(map[string]*regexp.Regexp, 2)

type Feed struct {
	Name            string        `json:"name"`
	Desc            string        `json:"desc"`
	Feed            string        `json:"feed"`
	Render          string        `json:"render"`
	Proxy           string        `json:"proxy"`
	Items           *ExprField    `json:"items"`
	FeedTitle       *ExprField    `json:"feed_title"`
	FeedDescription *ExprField    `json:"feed_description"`
	Title           *ExprField    `json:"title"`
	Description     *ExprField    `json:"description"`
	Link            *ExprField    `json:"link"`
	Author          *ExprField    `json:"author"`
	Email           *ExprField    `json:"email"`
	Created         *TimeField    `json:"created"`
	Content         *ContentField `json:"content"`
	Headers         []string      `json:"headers"`
}

func (f *Feed) Merge(other *Feed) {

	defer func() {
		if len(f.Headers) < 1 {
			f.Headers = append(f.Headers, UA)
		}

		if f.Render == "" {
			f.Render = "rss"
		}

		if f.Content.Blocks != nil {
			for _, v := range f.Content.Blocks {
				Regexes[v] = regexp.MustCompile(v)
			}
		}

	}()

	if other == nil {
		return
	}
	f.Headers = append(f.Headers, other.Headers...)
	if f.Proxy == "" {
		f.Proxy = other.Proxy
	}

	if f.Render == "" {
		f.Render = other.Render
	}

	if f.Proxy == "NO" {
		f.Proxy = ""
	}

	if f.FeedTitle == nil {
		f.FeedTitle = other.FeedTitle
		if f.FeedTitle == nil {
			f.FeedTitle = &ExprField{
				Expr: "/html/head/title",
			}
		}
	}

	if f.FeedDescription == nil {
		f.FeedDescription = other.FeedDescription
		f.FeedDescription = &ExprField{
			Expr: "/html/head/media[@name=\"description\"]",
			Attr: "content",
		}
	}

	if f.Headers == nil {
		f.Headers = []string{}
	} else {

		if other.Headers != nil {
			f.Headers = append(f.Headers, other.Headers...)
		}
	}

}
