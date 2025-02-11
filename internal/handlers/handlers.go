package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"forum/internal"
	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	data, valid := auth.IsValidFormValues(auth.FormErrors)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		data.User.CurrentPath = "/login"
		internal.Pagess.All_Templates.ExecuteTemplate(w, "login.html", data)
		auth.FormErrors = models.FormErrors{}
		auth.FormsData = models.FormsData{}
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
	internal.Pagess.All_Templates.ExecuteTemplate(w, "login.html", data)
}

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
	w.Write(internal.Pagess.Buf.Bytes())
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	data, valid := auth.IsValidFormValues(auth.FormErrors)
	if !valid {
		data.User.CurrentPath = "/register"
		w.WriteHeader(http.StatusBadRequest)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "register.html", data)
		auth.FormErrors = models.FormErrors{}
		auth.FormsData = models.FormsData{}
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
		logger.LogWithDetails(err)
		return
	}
	if InvalidCreatePostForm {
		w.WriteHeader(http.StatusBadRequest)
		data.FormsData = CreatePostFormData
		data.FormsData.FormErrors = CreatePostFormErrors
		internal.Pagess.All_Templates.ExecuteTemplate(w, "createpost.html", data)
		data.FormsData = models.FormsData{}
		data.FormsData.FormErrors = models.FormErrors{}
		InvalidCreatePostForm = false
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
		logger.LogWithDetails(err)
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

func FilterPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}
	r.ParseForm()
	Categories := r.Form["filter-category"]
	if len(Categories) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	placeholders := strings.Repeat("?,", len(Categories)-1) + "?"
	query := fmt.Sprintf(`
		SELECT 
			posts.id,
			posts.title,
			posts.content,
			posts.total_likes,
			posts.total_dislikes,
			posts.created_at,
			users.userName,
			users.id
		FROM posts
		JOIN users ON posts.user_id = users.id
		JOIN categories ON posts.id = categories.post_id
		WHERE categories.category IN (%s)
	`, placeholders)

	for _, val := range Categories {
		query = strings.Replace(query, "?", string('"')+val+string('"'), 1)
	}
	data, err := database.Fetch_Database(r, query, -1, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	internal.Pagess.Buf.Reset()
	err = internal.Pagess.All_Templates.ExecuteTemplate(&internal.Pagess.Buf, "home.html", data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	w.Write(internal.Pagess.Buf.Bytes())
	// internal.Pagess.All_Templates.ExecuteTemplate(w, "home.html", data)
}
func Gategoties_Checker( Gategories []string) (int64, error) {
	for _, val := range Gategories {
		stm, Err := database.Database.Prepare("SELECT EXISTS (SELECT 1 FROM  stoke_categories WHERE category = ?)")
		if Err != nil {

			return 500, Err
		}
		var exists bool
		stm.QueryRow(val).Scan(&exists)
		fmt.Println(exists)
		if !exists {
			return 400, fmt.Errorf("%s", "category does not exist")
		}
	}
	return 200, nil
}