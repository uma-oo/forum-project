package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"forum/internal/database"
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
	ORDER BY 
		posts.created_at DESC
`
	data := database.Fetch_Database(r, query, -1)
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
	ORDER BY 
		posts.created_at DESC
`
	data := database.Fetch_Database(r, query, -1)
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
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}

	// Retrieve the user token from the cookie
	Token, errtoken := r.Cookie("token")
	if errtoken != nil {
		fmt.Println(errtoken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var id int
	err := database.Database.QueryRow("SELECT id FROM users WHERE token = $1", Token.Value).Scan(&id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
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
			INNER JOIN 
				users
			ON 
				posts.user_id = users.id
			WHERE 
				users.id = ?
			ORDER BY 
				posts.created_at DESC;

	`
	data := database.Fetch_Database(r, query, id)
	Pagess.All_Templates.ExecuteTemplate(w, "myposts.html", data)
}

// todo : complete handeler for liked posts

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

func LikedPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Aloha")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	Token, errToken := r.Cookie("token")
	if errToken != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed hhh")
		return
	}
	var id int
	err := database.Database.QueryRow("SELECT id FROM users WHERE token = $1", Token.Value).Scan(&id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
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

			WHERE  post_reaction.user_id = ? AND  post_reaction.reaction = 1
			ORDER BY 
				posts.created_at DESC;

	`
	data := database.Fetch_Database(r, query, id)
	Pagess.All_Templates.ExecuteTemplate(w, "likedposts.html", data)
}
