package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"forum/internal/utils"
)

type Server struct {
	Log bool
}
type Pages struct {
	All_Templates *template.Template
}
type Form struct {
	Title  string
	Button string
}

var Pagess Pages

func ParseTemplates() {
	var err error
	path, err := utils.GetFolderPath("..", "templates")
	fmt.Printf("path: %v\n", path)
	if err != nil {
		panic(err)
	}
	Pagess.All_Templates, err = template.ParseGlob("./web/templates" + "/*.html")
	if err != nil {
		log.Fatal(err)
	}
	Pagess.All_Templates, err = Pagess.All_Templates.ParseGlob("../forum/web/components" + "/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Page not found")
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "home.html", nil)
}

func Serve_Static(w http.ResponseWriter, r *http.Request) {
	path, _ := utils.GetFolderPath("..", "static")
	fs := http.FileServer(http.Dir(path))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
