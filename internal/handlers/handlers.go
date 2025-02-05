package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"forum/internal/database"
	"forum/internal/models"
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
	err := Pagess.All_Templates.ExecuteTemplate(w, "home.html", data)
	fmt.Printf("err: %v\n", err)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	data := models.Data{}
	data.User.CurrentPath = r.URL.Path
	err := Pagess.All_Templates.ExecuteTemplate(w, "login.html", data)
	fmt.Printf("err: %v\n", err)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	data := models.Data{}
	data.User.CurrentPath = r.URL.Path
	err := Pagess.All_Templates.ExecuteTemplate(w, "register.html", data)
	fmt.Printf("err: %v\n", err)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	data := database.Fetch_Database(r)
	Pagess.All_Templates.ExecuteTemplate(w, "createpost.html", data)
}

func Post(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside single post")
}

// todo : complete handeler for created posts
func MyPosts(w http.ResponseWriter, r *http.Request) {
}

// todo : complete handeler for liked posts
func LikedPosts(w http.ResponseWriter, r *http.Request) {
}

func Serve_Files(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}

	path := r.URL.Path[1:]
	fileinfo, err := os.Stat(path)
	if err != nil || fileinfo.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "File Not Found")
		return
	}
	http.ServeFile(w, r, path)
}
