package main

import (
	"flag"
	"go2rss/mvc"
	"go2rss/util"
	"net/http"
	"path/filepath"
)

var dir = flag.String("dir", "./config", "Input Your Name")

func main() {
	flag.Parse()
	feedMap, _ := util.ReadConfig(*dir)
	mvc.Route("/(.*)+(\\.)[txt|gif|ico|js|css]+(.*)", mvc.Files)
	mvc.Route("^/$", func(rw http.ResponseWriter, r *http.Request) {
		mvc.Html(rw, "templates/index.html", feedMap)
	})
	mvc.Route("^/(.*)+", func(rw http.ResponseWriter, r *http.Request) {
		feed := feedMap[filepath.FromSlash(r.URL.Path)]
		if feed != nil {
			content, err := util.Gen(feed)
			if err != nil {
				rw.Write([]byte("err to gen"))
			}
			rw.Write([]byte(content))
		}
	})
	mvc.Run(8081)
}
