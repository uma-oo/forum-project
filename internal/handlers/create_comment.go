package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
)

var (
	InvalidComment  error
	InvalidReaction bool
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var user_id int
	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, "error.html", http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	conn, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	err = conn.QueryRow("SELECT id FROM users WHERE token = ?", cookie.Value).Scan(&user_id)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	comment := r.FormValue("comment")
	post_id := r.FormValue("post_id")

	err = IsValidComment(comment)
	if err != nil {
		logger.LogWithDetails(err)
		InvalidComment = err
		http.Redirect(w, r, "/posts?id="+post_id, http.StatusSeeOther)
		return
	}

	if comment == "" {
		utils.RenderTemplate(w, "error.html", http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	// insertiw data f blastha
	query := `INSERT INTO comments (user_id, post_id,content) VALUES (?,?,?)`
	statement, err := conn.Prepare(query)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(user_id, post_id, comment)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts?id="+post_id, http.StatusSeeOther)
}

func IsValidComment(comment string) error {
	if len(comment) > 10000 {
		return fmt.Errorf("exeeded max length allowed for comments (10000)")
	}
	if comment == "" {
		return fmt.Errorf("comment can't be empty")
	}
	return nil
}

// update comment reaction (like and dislike)
func ReactComment(res http.ResponseWriter, req *http.Request) {
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
	comment_id := req.FormValue("comment_id")
	user, err := UserData(req, "token", req.URL.Path)
	if err != nil {
		utils.RenderTemplate(res, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !utils.IsIdExist("comments", "id", comment_id) {
		fmt.Println("bad request")
		utils.RenderTemplate(res, "error.html", models.BadRequest, http.StatusBadRequest)
		return
	}

	query := `SELECT reaction_id FROM comment_reactions WHERE comment_id = ? AND user_id =? `
	statement, err := conn.Prepare(query)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	err = statement.QueryRow(comment_id, user.UserId).Scan(&reaction_found)
	if err == sql.ErrNoRows {
		query = `INSERT INTO comment_reactions (user_id ,comment_id,reaction_id) VALUES (?,?,?)`
		statement, err = conn.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		statement.Exec(user.UserId, comment_id, reaction)
	} else if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	if reaction == reaction_found {
		query = `UPDATE comment_reactions SET reaction_id = ? WHERE comment_id = ? AND user_id = ?`
		statement, err = conn.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		statement.Exec("0", comment_id, user.UserId)
	} else if reaction_found == "0" || reaction_found == "1" || reaction_found == "-1" {
		query = `UPDATE comment_reactions SET reaction_id = ? WHERE comment_id = ? AND user_id = ?`
		statement, err = conn.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(res, "error.html", http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		statement.Exec(reaction, comment_id, user.UserId)
	}
	http.Redirect(res, req, req.Referer(), http.StatusSeeOther)
	// utils.RenderTemplate(res, "post.html", data, http.StatusAccepted)
}

func FetchComments(post_id string) ([]models.Comment, error) {
	var comments []models.Comment
	query := `SELECT * FROM comments WHERE post_id = ? ORDER BY id DESC`
	conn, err := database.NewDatabase()
	if err != nil {
		return nil, err
	}
	statement, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(post_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.CommentId, &comment.PostId,
			&comment.UserId, &comment.TotalLikes, &comment.TotalDeslikes,
			&comment.CommentCreatedAt,
			&comment.CommentContent); err != nil {
			return nil, err
		}
		comment.CommentCreator, err = FetchCommentCreator(strconv.Itoa(comment.UserId))
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	defer rows.Close()
	return comments, nil
}

func UserData(r *http.Request, cookieName string, currentPath string) (*models.User, error) {
	var user models.User
	if !utils.IsCookieSet(r, cookieName) {
		return &models.User{
			CurrentPath: currentPath,
			IsLoged:     false,
			UserName:    "",
			UserEmail:   "",
			UserId:      "",
		}, nil
	} else {
		cookie, _ := r.Cookie(cookieName)
		db, err := database.NewDatabase()
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		query := `SELECT id, userEmail , userName FROM users where token= ? `
		statement, err := db.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		row := statement.QueryRow(cookie.Value)
		row.Scan(&user.UserId, &user.UserEmail, &user.UserName)
		return &models.User{
			IsLoged:     true,
			CurrentPath: currentPath,
			UserName:    user.UserName,
			UserEmail:   user.UserEmail,
			UserId:      user.UserId,
		}, nil
	}
}

func FetchCategories(post_id string) ([]models.Categorie, error) {
	var categories []models.Categorie
	db, err := database.NewDatabase()
	if err != nil {
		return nil, err
	}
	query := `SELECT  category FROM post_categories WHERE post_id = ?`
	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(post_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category models.Categorie
		err = rows.Scan(&category.CatergoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	defer rows.Close()
	return categories, nil
}

func FetchCommentCreator(user_id string) (string, error) {
	var username string
	conn, err := database.NewDatabase()
	if err != nil {
		return "", err
	}
	query := "SELECT userName FROM users INNER JOIN comments ON users.id = comments.user_id WHERE comments.user_id = ?"
	statement, err := conn.Prepare(query)
	if err != nil {
		return "", err
	}
	err = statement.QueryRow(user_id).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func FetchPostCreator(user_id string) (string, error) {
	var username string
	conn, err := database.NewDatabase()
	if err != nil {
		return "", err
	}
	query := "SELECT userName FROM users INNER JOIN posts ON users.id = posts.user_id WHERE posts.user_id = ?"
	statement, err := conn.Prepare(query)
	if err != nil {
		return "", err
	}
	err = statement.QueryRow(user_id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func AllCategories() ([]models.Categorie, error) {
	var all_categories []models.Categorie
	query := `SELECT category FROM categories`
	conn, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	statement, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var categorie models.Categorie
		err := rows.Scan(&categorie.CatergoryName)
		if err != nil {
			return nil, err
		}
		all_categories = append(all_categories, categorie)
	}
	defer rows.Close()
	return all_categories, nil
}
