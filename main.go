// From this lovely https://www.chazzuka.com/2015/03/load-parse-json-file-golang/
package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"html/template"
	"regexp"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseGlob("templates/*.tmpl"))
// var templates *template.Template


var validPath = regexp.MustCompile("^(/$|/([a-zA-Z0-9]+))$")

type Stall struct {
	Id int	`json:"id"`
	StallName string `json:"stall-name"`
	StallDesc string `json:"stall-desc"`
	Contact string	`json:"contact"`
	Items []Item	`json:"items"`
}

type Item struct {
	ItemName string `json:"item-name"`
	ItemDesc string `json:"item-desc"`
	Price float64
}

func (s Stall) toString() string {
	return toJson(s)
}

// Convert the stall interface to JSON
func toJson(s interface{}) string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	return string(bytes)
}

// Read all the stalls from a JSON file
func getStalls() []Stall {
	raw, err := ioutil.ReadFile("stalls.json")
	if err != nil {
		fmt.Println(err)
	}

	var s []Stall
	json.Unmarshal(raw, &s)
	return s
}

func indexHandler(w http.ResponseWriter, r *http.Request, name string) {
	stalls := getStalls()
	templates.ExecuteTemplate(w, "base", nil)
	templates.ExecuteTemplate(w, "index", stalls)
	return
}

func stallHandler(w http.ResponseWriter, r *http.Request, name string) {
	stalls := getStalls()
	for _, s := range stalls {
		if strings.ToLower(s.StallName) == strings.ToLower(name) {
			templates.ExecuteTemplate(w, "base", nil)
			templates.ExecuteTemplate(w, "content", s)
			return
		}
	}
	http.NotFound(w, r)

}

// Make handlers, so that regex can check path
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract page title
		m := validPath.FindStringSubmatch(r.URL.Path)

		// Route to 404, or to stallHandler
		if m == nil {
			http.NotFound(w, r)
			return
		}
		if m[2] != "" {
			stallHandler(w, r, m[2])
			return
		}

		// If not a 404 or specific stall, generate index
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/", makeHandler(indexHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}