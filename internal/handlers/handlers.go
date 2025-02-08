package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"forum/internal/database"
	"forum/internal/utils"
	"forum/pkg/logger"
)

type Pages struct {
	All_Templates *template.Template
	buf           bytes.Buffer
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
	query := `
	SELECT 
		posts.id,posts.title, posts.content, posts.total_likes, posts.total_dislikes, posts.created_at,
		users.userName, users.id
	FROM 
		posts
	INNER JOIN 
		users
	ON 
		posts.user_id = users.id
	
`
	data, errr := database.Fetch_Database(r, query, -1, false)
	if errr != nil {
		log.Fatal(errr)
	}
	Pagess.buf.Reset()
	err := Pagess.All_Templates.ExecuteTemplate(&Pagess.buf, "home.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "home.html", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	if utils.IsCookieSet(r, "token") {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	Pagess.buf.Reset()
	err := Pagess.All_Templates.ExecuteTemplate(&Pagess.buf, "login.html", nil)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "login.html", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	if utils.IsCookieSet(r, "token") {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	Pagess.buf.Reset()
	err := Pagess.All_Templates.ExecuteTemplate(&Pagess.buf, "register.html", nil)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "register.html", nil)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	query := `
	SELECT 
		posts.id,posts.title, posts.content, posts.total_likes, posts.total_dislikes, posts.created_at,
		users.userName, users.id
	FROM 
		posts
	INNER JOIN 
		users
	ON 
		posts.user_id = users.id

`
	data, err := database.Fetch_Database(r, query, -1, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	Pagess.buf.Reset()
	err = Pagess.All_Templates.ExecuteTemplate(&Pagess.buf, "createpost.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "createpost.html", data)
}

// // todo : complete handeler for single post
func Post(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside single post")
}

func MyPosts(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	// Retrieve the user token from the cookie
	Token, _ := r.Cookie("token")

	var id int
	stm, err := database.Database.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	err = stm.QueryRow(Token.Value).Scan(&id)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	query := `
	SELECT 
		posts.id,posts.title, posts.content, posts.total_likes, posts.total_dislikes, posts.created_at,
		users.userName, users.id
	FROM 
		posts
	INNER JOIN 
		users
	ON 
		posts.user_id = users.id
	
`
	data, err := database.Fetch_Database(r, query, id, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	err = Pagess.All_Templates.ExecuteTemplate(&Pagess.buf, "myposts.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "myposts.html", data)
}

// todo : complete handeler for liked posts

func Serve_Files(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	path := r.URL.Path[1:]
	fileinfo, err := os.Stat(path)
	if err != nil || fileinfo.IsDir() {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusNotFound)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "404 page Not Found")
		return
	}
	http.ServeFile(w, r, path)
}

func LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	Token, errToken := r.Cookie("token")
	if errToken != nil {
		logger.LogWithDetails(errToken)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", " 500 Internal Server Error")
		return
	}
	var id int
	stm, err := database.Database.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	err = stm.QueryRow(Token.Value).Scan(&id)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	query := `
		SELECT 
			posts.id,
			posts.title,
			posts.content,
			posts.total_likes,
			posts.total_dislikes,
			posts.created_at,
			users.userName,
			users.id
			FROM 
				posts	
			JOIN users ON posts.user_id = users.id
			JOIN  post_reaction ON posts.id = post_reaction.post_id
	`
	data, err := database.Fetch_Database(r, query, id, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	Pagess.buf.Reset()
	err = Pagess.All_Templates.ExecuteTemplate(&Pagess.buf, "likedposts.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "likedposts.html", data)
}
