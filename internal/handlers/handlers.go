package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	if utils.IsCookieSet(r, "token") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data, valid := auth.IsValidFormValues(auth.FormErrors)
	if !valid {
		data.User.CurrentPath = "/login"
		utils.RenderTemplate(w, "login.html", data, http.StatusBadRequest)
		auth.FormErrors = models.FormErrors{}
		auth.FormsData = models.FormsData{}
		return
	}
	data.User.CurrentPath = r.URL.Path
	utils.RenderTemplate(w, "login.html", data, http.StatusOK)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.RenderTemplate(w, "error.html", models.PageNotFound, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	query := `
	SELECT 
		posts.id,posts.title, posts.content, posts.total_likes, posts.total_dislikes, posts.total_comments,posts.created_at,
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
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	user_data, err := UserData(r, "token", "")
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	data.User = *user_data

	utils.RenderTemplate(w, "home.html", data, http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	if utils.IsCookieSet(r, "token") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data, valid := auth.IsValidFormValues(auth.FormErrors)
	if !valid {
		data.User.CurrentPath = "/register"
		utils.RenderTemplate(w, "register.html", data, http.StatusBadRequest)
		auth.FormErrors = models.FormErrors{}
		auth.FormsData = models.FormsData{}
		return
	}

	data.User.CurrentPath = r.URL.Path
	utils.RenderTemplate(w, "register.html", nil, http.StatusOK)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	data := models.Data{}
	var err error
	data.Categories, err = AllCategories()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	user, err := UserData(r, "token", "")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	data.User = *user
	if InvalidCreatePostForm {
		data.FormsData = CreatePostFormData
		data.FormsData.FormErrors = CreatePostFormErrors
		utils.RenderTemplate(w, "createpost.html", data, http.StatusBadRequest)
		CreatePostFormErrors = models.FormErrors{}
		CreatePostFormData = models.FormsData{}
		InvalidCreatePostForm = false
		return
	}
	utils.RenderTemplate(w, "createpost.html", data, http.StatusOK)
}

func Post(w http.ResponseWriter, r *http.Request) {
	var data models.Data
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	conn, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	var PostData models.Post
	PostId := r.URL.Query().Get("id")
	query := `SELECT * FROM posts WHERE id = ?`
	if err := conn.QueryRow(query, PostId).Scan(&PostData.PostId, &PostData.PostCreatedAt,
		&PostData.UserID, &PostData.PostTitle, &PostData.PostContent, &PostData.TotalLikes,
		&PostData.TotalDeslikes, &PostData.TotalComments); err != nil {
		if err == sql.ErrNoRows {
			utils.RenderTemplate(w, "error.html", models.PageNotFound, http.StatusNotFound)
			return
		}
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	PostData.PostCreator, err = FetchPostCreator(strconv.Itoa(PostData.UserID))
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	PostData.Categories, err = FetchCategories(PostId)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	PostData.Comments, err = FetchComments(PostId)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	data.Posts = append(data.Posts, PostData)
	data_user, err := UserData(r, "token", "/posts")
	data.User = *data_user
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	data.Categories, err = AllCategories()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	if InvalidComment != nil {
		data.InvalidComment = InvalidComment.Error()
		utils.RenderTemplate(w, "post.html", data, http.StatusBadRequest)
		InvalidComment = nil
		return
	}

	utils.RenderTemplate(w, "post.html", data, http.StatusOK)
}

func MyPosts(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the user token from the cookie
	Token, _ := r.Cookie("token")
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}

	var id int
	stmt, err := db.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = stmt.QueryRow(Token.Value).Scan(&id)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	query := `
	SELECT 
		posts.id,posts.title, posts.content, posts.total_likes, posts.total_dislikes, posts.total_comments,posts.created_at,
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
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	user_data, err := UserData(r, "token", "")
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	data.User = *user_data
	utils.RenderTemplate(w, "myposts.html", data, http.StatusOK)
}

func Serve_Files(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path[1:]
	fileinfo, err := os.Stat(path)
	if err != nil || fileinfo.IsDir() {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.PageNotFound, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}

func LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	Token, err := r.Cookie("token") // Is this is not handelled with the middleware ????????????
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	var id int
	stmt, err := db.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = stmt.QueryRow(Token.Value).Scan(&id)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	query := `
		SELECT 
			posts.id,
			posts.title,
			posts.content,
			posts.total_likes,
			posts.total_dislikes,
			posts.total_comments,
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
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	user_data, err := UserData(r, "token", "")
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	data.User = *user_data
	utils.RenderTemplate(w, "likedposts.html", data, http.StatusOK)
}

func FilterPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	Categories := r.Form["filter-category"]

	if len(Categories) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	errNum, err := Gategoties_Checker(Categories)
	if errNum == 500 {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	} else if errNum == 400 {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.BadRequest, http.StatusBadRequest)
		return
	}
	placeholders := strings.Repeat("?,", len(Categories)-1) + "?"
	query := fmt.Sprintf(`
		SELECT DISTINCT
			posts.id,
			posts.title,
			posts.content,
			posts.total_likes,
			posts.total_dislikes,
			posts.total_comments,
			posts.created_at,
			users.userName,
			users.id
		FROM posts
		JOIN users ON posts.user_id = users.id
		JOIN post_categories ON posts.id = post_categories.post_id
		WHERE post_categories.category IN (%s)
	`, placeholders)

	for _, val := range Categories {
		query = strings.Replace(query, "?", string('"')+val+string('"'), 1)
	}
	data, err := database.Fetch_Database(r, query, -1, true)
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	user_data, err := UserData(r, "token", "")
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	data.User = *user_data
	utils.RenderTemplate(w, "home.html", data, http.StatusOK)
}

func Gategoties_Checker(Gategories []string) (int64, error) {
	db, err := database.NewDatabase()
	if err != nil {
		return 500, err
	}
	defer db.Close()
	for _, val := range Gategories {
		stmt, Err := db.Prepare("SELECT EXISTS (SELECT 1 FROM  categories WHERE category = ?)")
		if Err != nil {
			return 500, Err
		}
		var exists bool
		stmt.QueryRow(val).Scan(&exists)
		if !exists {
			return 400, fmt.Errorf("%s", "category does not exist")
		}
	}
	return 200, nil
}
