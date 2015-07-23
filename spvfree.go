package spvfree

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
	"regexp"
	"appengine"
	"appengine/datastore"
	//"appengine/user"
)

type Sheep struct {
	Name string
	Email string
	Addr string
	Phon string
	Area string 
	Note string
	JDte time.Time
}

var templates = template.Must(template.ParseFiles("reg.html", "ta.html", "how.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, s *Sheep) {
	err := templates.ExecuteTemplate(w, tmpl+".html", s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	c := appengine.NewContext(r)
	q0 := datastore.NewQuery("Sheep").Filter("Name=", "Guest")
	var s0 []*Sheep
	
	key, err := q0.GetAll(c, &s0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if key != nil {
		renderTemplate(w, "reg", s0[0])
		return
	}
	renderTemplate(w, "reg", &Sheep{Name: "Guest",Email: "guest@user.com",Addr: "123 JP Nagar", Phon: "123456", Area: "123", Note: "Rock on", JDte: time.Now()})
}

func handleRegPage(w http.ResponseWriter, r *http.Request ) {
	if r.Method != "GET" {
		http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
		return
	}
	c := appengine.NewContext(r)
	fmt.Sprintf("%s", "Creating sheep w name: " + r.FormValue("name"))
	u1 := Sheep { Name: r.FormValue("name"), Email: r.FormValue("email"), Addr: r.FormValue("addr"),
			Area: r.FormValue("area"), Note: r.FormValue("note"), JDte: time.Now() }
	if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"sheep", nil), &u1); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "ta", &u1)
}

func handleHowPage(w http.ResponseWriter, r *http.Request ) {
	if r.Method != "GET" {
		http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
		return
	}
        renderTemplate(w, "how", nil)
}

var validPath = regexp.MustCompile("^/(reg|ta|how)?/?$")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) { 
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w,r)
			return
		}
		fn(w, r) //, m[2])
	}
}

func init() {
	http.HandleFunc("/", makeHandler(handleMainPage))
	http.HandleFunc("/reg", makeHandler(handleRegPage))
	http.HandleFunc("/how", makeHandler(handleHowPage))
}
