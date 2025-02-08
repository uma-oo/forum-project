package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/internal"
	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found")
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
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
	internal.Pagess.Buf.Reset()
	err := internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "home.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	internal.Pagess.All_Templates.ExecuteTemplate(w, "home.html", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside Login")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	data, valid := auth.IsValidFormValues(auth.FormErrors)
	if !valid {
		fmt.Printf("data: %v\n", data)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	if utils.IsCookieSet(r, "token") {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	internal.Pagess.Buf.Reset()
	err := internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "login.html", nil)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	data.User.CurrentPath = r.URL.Path
	err = internal.Pagess.All_Templates.ExecuteTemplate(w, "login.html", data)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	if utils.IsCookieSet(r, "token") {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	internal.Pagess.Buf.Reset()
	err := internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "register.html", nil)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	data := models.Data{}
	data.User.CurrentPath = r.URL.Path
	err = internal.Pagess.All_Templates.ExecuteTemplate(w, "register.html", data)
	fmt.Printf("err: %v\n", err)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
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
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	internal.Pagess.Buf.Reset()
	err = internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "createpost.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	internal.Pagess.All_Templates.ExecuteTemplate(w, "createpost.html", data)
}

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
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
	data, _ := database.Fetch_Database(r, query, -1, false)
	data.Posts = data.Posts[0:1]
	fmt.Println(internal.Pagess.All_Templates.ExecuteTemplate(w, "post.html", data))
}

func MyPosts(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	// Retrieve the user token from the cookie
	Token, _ := r.Cookie("token")

	var id int
	stm, err := database.Database.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	err = stm.QueryRow(Token.Value).Scan(&id)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
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
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	err = internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "myposts.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	internal.Pagess.All_Templates.ExecuteTemplate(w, "myposts.html", data)
}

func Serve_Files(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	path := r.URL.Path[1:]
	fileinfo, err := os.Stat(path)
	if err != nil || fileinfo.IsDir() {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusNotFound)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "404 page Not Found")
		return
	}
	http.ServeFile(w, r, path)
}

func LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	Token, errToken := r.Cookie("token")
	if errToken != nil {
		logger.LogWithDetails(errToken)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", " 500 Internal Server Error")
		return
	}
	var id int
	stm, err := database.Database.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	err = stm.QueryRow(Token.Value).Scan(&id)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
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
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	internal.Pagess.Buf.Reset()
	err = internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "likedposts.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	internal.Pagess.All_Templates.ExecuteTemplate(w, "likedposts.html", data)
}
