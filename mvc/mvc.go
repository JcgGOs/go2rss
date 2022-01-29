package mvc

import (
	"fmt"
	"net/http"
	"regexp"
	"text/template"
)

var routeMap = make(map[*regexp.Regexp]func(http.ResponseWriter, *http.Request), 4)

func Route(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	routeMap[regexp.MustCompile(pattern)] = handler
}

func Run(port int) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		for re, f := range routeMap {
			if re.MatchString(r.URL.RequestURI()) {
				fmt.Println("match", r.RequestURI)

				f(w, r)
			}
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func Html(rw http.ResponseWriter, path string, data interface{}) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	return t.Execute(rw, data)
}

var fileHandler = http.FileServer(http.Dir("./static"))

func Files(rw http.ResponseWriter, r *http.Request) {
	fileHandler.ServeHTTP(rw, r)
}
