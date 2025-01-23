package handlers

import (
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

	Pagess.All_Templates, err = template.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	Pagess.All_Templates, err = Pagess.All_Templates.ParseGlob("./web/components/*.html")
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
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed hassan")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "home.html", nil)
}

// this is not a good way to serve static files

func Serve_Static(w http.ResponseWriter, r *http.Request) {
	path, _ := utils.GetFolderPath("..", "static")
	fs := http.FileServer(http.Dir(path))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}

// MODIFY this one after finishing the task

// func StaticHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	if !strings.HasPrefix(r.URL.Path, "/static") {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	} else {
// 		file_info, err := os.Stat(r.URL.Path[1:])
// 		fmt.Println(file_info.Name())
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		} else if file_info.IsDir() {
// 			w.WriteHeader(http.StatusForbidden)
// 			return
// 		} else {
// 			http.ServeFile(w, r, r.URL.Path[1:])
// 		}
// 	}
// }
