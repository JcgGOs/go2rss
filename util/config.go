package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go2rss/model"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

func WalkFiles(dir string) (map[string]string, error) {

	folder, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]string, 4)
	err = filepath.Walk(folder, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		// /full_name.json => /full_name
		relationPath := strings.ReplaceAll(path, folder, "")
		relationPath = relationPath[:len(relationPath)-len(filepath.Ext(relationPath))]

		configMap[relationPath] = path
		return nil

	})
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func NoExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func ReadConfig(dir string) (map[string]*model.Feed, error) {
	if dir == "" {
		dir = "./config/"
	}

	if !strings.HasPrefix(dir, "/") {
		dir = dir + "/"
	}

	pathMap, err := WalkFiles(dir)
	if err != nil {
		return nil, err
	}

	var _default *model.Feed
	feedsMap := make(map[string]*model.Feed, 4)
	for nameNoExt, path := range pathMap {

		feed, err := ToFeed(path)
		if err != nil {
			fmt.Println("Convert error:", path, err)
			continue
		}

		if nameNoExt == filepath.FromSlash("/default") {
			_default = feed
			continue
		}

		feedsMap[nameNoExt] = feed
	}

	for _, f := range feedsMap {
		f.Merge(_default)
	}

	return feedsMap, nil
}

func ToFeed(file string) (*model.Feed, error) {
	text, err := os.ReadFile(file)
	fmt.Printf("text: %v\n", string(text))

	if err != nil {
		return nil, err
	}

	c := &model.Feed{}
	if err := json.Unmarshal(text, c); err != nil {
		return nil, err
	}

	return c, nil
}

func Gen(feed *model.Feed) (string, error) {
	body, err := GET(feed.Proxy, feed.Feed, feed.Headers)
	if err != nil {
		return "", err
	}

	now := time.Now()
	doc, _ := Load(bytes.NewReader(body))
	nFeed := &feeds.Feed{
		Title:       Value(doc, feed.FeedTitle),
		Link:        &feeds.Link{Href: feed.Feed},
		Description: Value(doc, feed.FeedDescription),
		Author:      &feeds.Author{Name: Value(doc, feed.Author), Email: Value(doc, feed.Email)},
		Created:     now,
		Items:       []*feeds.Item{},
	}

	for _, n := range Nodes(doc, feed.Items.Expr) {
		_href := Value(n, feed.Link)
		_link := AbsHref(_href, feed.Feed)
		_item := &feeds.Item{
			Title:       Value(n, feed.Title),
			Link:        &feeds.Link{Href: _link},
			Description: Value(n, feed.Description),
			Created:     TimeValue(n, feed.Created),
			Id:          _link,
		}
		nFeed.Items = append(nFeed.Items, _item)
	}

	if feed.Content != nil && feed.Content.Expr != "" && len(nFeed.Items) > 0 {

		var wg sync.WaitGroup
		wg.Add(len(nFeed.Items))

		ch := make(chan *feeds.Item, feed.Content.NThread)
		for i, f := range nFeed.Items {
			ch <- f
			go func(_f *feeds.Item, idx int) {
				defer wg.Done()

				if feed.Content.NTop == 0 || idx < feed.Content.NTop {

					//random to delay
					if feed.Content.NDelay > 0 {
						time.Sleep(time.Duration(rand.Intn(feed.Content.NDelay)) * time.Millisecond)
					}

					href := AbsHref(_f.Link.Href, feed.Feed)
					itemBody, err := GET(feed.Proxy, href, feed.Headers)
					if err == nil {
						//Remove regex tag
						for _, v := range feed.Content.Blocks {
							re := model.Regexes[v]
							if re != nil {
								itemBody = []byte(re.ReplaceAllString(string(itemBody), ""))
							}
						}
						_itemDoc, _ := Load(bytes.NewReader(itemBody))
						_f.Content = Html(_itemDoc, feed.Content)
					}
				}
				<-ch

			}(f, i)
		}
		wg.Wait()
	}

	//
	switch feed.Render {
	case "atom":
		{
			return nFeed.ToAtom()
		}
	case "json":
		{
			return nFeed.ToJSON()
		}
	default:
		{
			return nFeed.ToRss()
		}
	}
}
