package main

import (
	"fmt"
	"go2rss/util"
	"net/http"
	"strings"
)

func main() {
	// PrintFetch()
	Boot()
	// Boot()
	// fileName := "hello.js.txt"
	// fmt.Println(fileName[:len(fileName)-len(filepath.Ext(fileName))])
}

func Boot() {
	feedMap, _ := util.ReadConfig("")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/" {
			var buf string
			for k, f := range feedMap {
				buf = buf + fmt.Sprintf("<tr><td><a href=\"%s\">%s</a></td><td>%s</td></tr>", k, f.Name, f.Desc)
			}

			var template = "<html><head> </head> <body style=\"text-align: center\"> <h2>aRss</h2> <table> <tr> <th>Name</th> <th>Desc</th> </tr> %s </table> </body> </html>"
			html := fmt.Sprintf(template, buf)
			w.Write([]byte(html))
			return
		}

		name := strings.Replace(r.RequestURI, "/", "", 1)
		feed := feedMap[name]
		if feed != nil {
			context, _ := util.Gen(feed)
			w.Write([]byte(context))
		}
	})
	fmt.Println(":8001")
	http.ListenAndServe(":8001", nil)
}

func PrintFetchNew() {
	feedMap, err := util.ReadConfig("")
	if err != nil {
		return
	}

	fmt.Println(feedMap)
	feed := feedMap["wangyin"]

	html, _ := util.Gen(feed)

	fmt.Println(html)

}
