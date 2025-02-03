package handlers

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/database"
	"forum/internal/utils"
)

type Pages struct {
	All_Templates *template.Template
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
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found")
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	data := database.Fetch_Database(r)
	Pagess.All_Templates.ExecuteTemplate(w, "home.html", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "login.html", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "register.html", nil)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	data := database.Fetch_Database(r)
	Pagess.All_Templates.ExecuteTemplate(w, "createpost.html", data)
}

// todo : complete handeler for single post
func Post(w http.ResponseWriter, r *http.Request) {
}

// todo : complete handeler for created posts
func MyPosts(w http.ResponseWriter, r *http.Request) {
}

// todo : complete handeler for liked posts
func LikedPosts(w http.ResponseWriter, r *http.Request) {
}

// todo : need to prevent folders listing
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
