package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
)

const (
	ReactionLike    = 1
	ReactionDislike = -1
	Neutre          = 0
)

var (
	CreatePostFormData    = models.FormsData{}
	CreatePostFormErrors  = models.FormErrors{}
	InvalidCreatePostForm = false
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	db, err := database.NewDatabase()
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	CreatePostFormData.PostGategoriesInput = r.Form["post-categorie"]
	CreatePostFormData.PostContentInput = r.FormValue("postBody")
	CreatePostFormData.PostTitleInput = r.FormValue("postTitle")
	errNum, err := Gategoties_Checker(CreatePostFormData.PostGategoriesInput)
	if errNum == 500 {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	} else if errNum == 400 {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.BadRequest, http.StatusBadRequest)

		return
	}
	// lets check for emptyness
	IsValidCreatePostForm()
	if InvalidCreatePostForm {
		http.Redirect(w, r, "/create_post", http.StatusFound)
		return
	}
	// get the user ID from the session
	cookie, _ := r.Cookie("token")
	// get the user ID from the users table
	var userId int
	stm, err := db.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = stm.QueryRow(cookie.Value).Scan(&userId)
	if err != nil {
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	db, err = database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
	}
	// lets insert this data to our database
	stm, err = db.Prepare("INSERT INTO posts (user_id,title,content) VALUES ( ?,?,?)")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = stm.Exec(userId, CreatePostFormData.PostTitleInput, CreatePostFormData.PostContentInput)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	// get the last inserted post id
	var postId int
	stm, err = db.Prepare("SELECT last_insert_rowid()")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = stm.QueryRow().Scan(&postId)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	// insert categories
	db, err = database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	stm, err = db.Prepare("INSERT INTO post_categories (category, post_id) VALUES (?, ?)")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	for _, category := range CreatePostFormData.PostGategoriesInput {
		_, err = stm.Exec(category, postId)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PostReactions(res http.ResponseWriter, req *http.Request) {
	conn, err := database.NewDatabase()
	if err != nil {
		utils.RenderTemplate(res, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	var reaction_found string
	if req.Method != http.MethodPost {
		utils.RenderTemplate(res, "error.html", http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	reaction := req.FormValue("reaction")
	post_id := req.FormValue("post_id")
	user, err := UserData(req, "token", req.URL.Path)
	if err != nil {
		utils.RenderTemplate(res, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	query := `SELECT reaction_id FROM post_reaction WHERE  post_id= ? AND user_id =? `
	statement, err := conn.Prepare(query)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	err = statement.QueryRow(post_id, user.UserId).Scan(&reaction_found)
	if err == sql.ErrNoRows {
		query = `INSERT INTO post_reaction (user_id ,post_id,reaction_id) VALUES (?,?,?)`
		statement, err = conn.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		statement.Exec(user.UserId, post_id, reaction)

	} else if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	if reaction == reaction_found {
		query = `UPDATE post_reaction SET reaction_id = ? WHERE post_id = ? AND user_id = ?`
		statement, err = conn.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		statement.Exec("0", post_id, user.UserId)
	} else if reaction_found == "0" || reaction_found == "1" || reaction_found == "-1" {
		query = `UPDATE post_reaction SET reaction_id = ? WHERE post_id = ? AND user_id = ?`
		statement, err = conn.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		statement.Exec(reaction, post_id, user.UserId)
	}

	var PostData models.Post

	query = `SELECT * FROM posts WHERE id = ?`
	if err := conn.QueryRow(query, post_id).Scan(&PostData.PostId, &PostData.PostCreatedAt,
		&PostData.UserID, &PostData.PostTitle, &PostData.PostContent, &PostData.TotalLikes,
		&PostData.TotalDeslikes, &PostData.TotalComments); err != nil {
		if err != sql.ErrNoRows {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
	}

	PostData.PostCreator, err = FetchPostCreator(strconv.Itoa(PostData.UserID))
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	PostData.Categories, err = FetchCategories(post_id)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, req.Referer(), http.StatusSeeOther)
}

func IsValidCreatePostForm() {
	if CreatePostFormData.PostTitleInput == "" {
		CreatePostFormErrors.InvalidPostTitle = "Post title is required"
		InvalidCreatePostForm = true
	}
	if len(CreatePostFormData.PostTitleInput) > 50 {
		CreatePostFormErrors.InvalidPostTitle = "Exeeded post title length (50)"
		InvalidCreatePostForm = true
	}
	if CreatePostFormData.PostContentInput == "" {
		CreatePostFormErrors.InvalidPostContent = "Post content is required"
		InvalidCreatePostForm = true
	}
	if len(CreatePostFormData.PostContentInput) >= 10000 {
		CreatePostFormErrors.InvalidPostContent = "Exeeded post content length (10000)"
		InvalidCreatePostForm = true
	}
	if len(CreatePostFormData.PostGategoriesInput) == 0 {
		CreatePostFormErrors.InvalidPostCategories = "Post categories are required - pick at least one"
		InvalidCreatePostForm = true
	}
}
