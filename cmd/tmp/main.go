package main

import (
	"flag"
	"go2rss/mvc"
	"go2rss/util"
	"net/http"
)

var dir = flag.String("dir", "./config", "Input Your Name")

func main() {
	flag.Parse()
	feedMap, _ := util.ReadConfig(*dir)
	// http.Stat
	// fs := http.FileServer(http.Dir("./static"))
	// fmt.Printf("regexp.MustCompile(\"*.txt\").MatchString(\"/xxxx.txt\"): %v\n", regexp.MustCompile("/(.*)+.txt").MatchString("/xxxx.txt"))
	mvc.Route("/(.*)+(\\.)[txt|gif|ico|js|css]+(.*)", mvc.Files)
	mvc.Route("/", func(rw http.ResponseWriter, r *http.Request) {
		mvc.Html(rw, "templates/index.html", feedMap)
	})
	mvc.Run(8081)
}
