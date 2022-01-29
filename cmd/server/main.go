package main

import (
	"flag"
	"fmt"
	"go2rss/util"
	"net/http"
	"os"
	"strings"
)

var dir = flag.String("dir", "./config", "Input Your Name")

func main() {
	flag.Parse()

	feedMap, _ := util.ReadConfig(*dir)

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		feedMap, _ = util.ReadConfig(*dir)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/" {
			var buf string
			for k, f := range feedMap {
				buf = buf + fmt.Sprintf("<tr><td><a href=\"%s\">%s</a></td><td>%s</td></tr>", k, f.Name, f.Desc)
			}

			var template = "<html><head> </head> <body style=\"text-align: center\"> <h2>aRss</h2> <table> <tr> <th>Name</th> <th>Desc</th> </tr> %s </table> </body> </html>"
			html := fmt.Sprintf(template, buf)
			if size, err := w.Write([]byte(html)); err == nil {
				fmt.Println(size, err)
			}
			return
		}

		name := strings.Replace(r.URL.Path, "/", "", 1)
		feed := feedMap[name]
		if feed != nil {
			context, _ := util.Gen(feed)
			if size, err := w.Write([]byte(context)); err == nil {
				fmt.Println(size, err)
			}
		}
	})
	fmt.Println(":8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		os.Exit(1)
	}
}
