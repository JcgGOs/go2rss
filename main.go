package main

import (
	"fmt"
	"go2rss/model"
	"net/http"
)

// func main() {
// 	fmt.Println("hello")
// 	PrintFetch()
// 	func main() {
// 		http.Handle("/", &indexHandler{content: "hello world!"})
// 		http.ListenAndServe(":8001", nil)
// 	}

// }

func main() {
	cfg, _ := model.Parse("")
	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		feed := cfg.GetFeed(name)
		if feed != nil {
			context, _ := feed.Gen()
			w.Write([]byte(context))
		}
	})
	fmt.Println(":8001")
	http.ListenAndServe(":8001", nil)
}

func PrintFetch() {
	cfg, _ := model.Parse("")

	for _, feed := range cfg.Feeds {
		content, err := feed.Gen()
		if err == nil {
			fmt.Printf("content: %v\n", content)
		}
	}
}
