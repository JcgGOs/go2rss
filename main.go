package main

import (
	"fmt"
	"go2rss/model"
)

func main() {
	fmt.Println("hello")
	PrintFetch()
}

func PrintFetch() {
	cfg, _ := model.Parse("")

	for _, feed := range cfg.Feeds {
		content, err := feed.FetchRss()
		if err == nil {
			fmt.Printf("content: %v\n", content)
		}
	}
}
