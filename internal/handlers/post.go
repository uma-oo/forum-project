package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal"
	"forum/internal/database"
	"forum/internal/models"
	"forum/pkg/logger"
)

const (
	ReactionLike    = 1
	ReactionDislike = -1
	Neutre          = 0
)

var (
	CreatePostFormData   = models.FormsData{}
	CreatePostFormErrors = models.FormErrors{}
	InvalidCreatePostForm = false
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	pages := internal.Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "405 method not allowed")
		return
	}

	r.ParseForm() ///////////////////////////
	CreatePostFormData.PostGategoriesInput = r.Form["post-categorie"]
	CreatePostFormData.PostContentInput = r.FormValue("postBody")
	CreatePostFormData.PostTitleInput = r.FormValue("postTitle")
	fmt.Println(CreatePostFormData.PostGategoriesInput)
	errNum , err := Gategoties_Checker(CreatePostFormData.PostGategoriesInput)
	if errNum == 500 {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return 
	} else if errNum == 400 {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "400 Bad Request ")
		return 
	}

	IsValidCreatePostForm()
	if InvalidCreatePostForm {
		fmt.Println("redirecting ")
		http.Redirect(w, r, "/create_post", http.StatusFound)
		return
	}
	// lets check for emptyness

	// get the user ID from the session
	cookie, _ := r.Cookie("token")

	// get the user ID from the users table
	var userId int
	stm, err := database.Database.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		internal.Pagess.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error ")
		return
	}
	err = stm.QueryRow(cookie.Value).Scan(&userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}

	// lets insert this data to our database
	stm, err = database.Database.Prepare("INSERT INTO posts (user_id,title,content) VALUES ( ?,?,?)")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	_, err = stm.Exec(userId, CreatePostFormData.PostTitleInput, CreatePostFormData.PostContentInput)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	// get the last inserted post id
	var postId int
	stm, err = database.Database.Prepare("SELECT last_insert_rowid()")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}
	err = stm.QueryRow().Scan(&postId)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}
	// insert categories
	stm, err = database.Database.Prepare("INSERT INTO categories (category, post_id) VALUES (?, ?)")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}
	for _, category := range CreatePostFormData.PostGategoriesInput {
		_, err = stm.Exec(category, postId)
		if err != nil {
			logger.LogWithDetails(err)
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "500 internal server error")
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PostReactions(w http.ResponseWriter, r *http.Request) {
	pages := internal.Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "405 method not allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "400 bad request")
		return
	}

	postid, errpost := strconv.Atoi(r.FormValue("post_id"))
	reaction, errreaction := strconv.Atoi(r.FormValue("reaction"))
	Token, terr := r.Cookie("token")
	if terr != nil {
		logger.LogWithDetails(terr)
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "400 bad request")
		return
	}

	if errpost != nil || errreaction != nil {
		logger.LogWithDetails(fmt.Errorf("%s", " Itoi Errors"))
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "400 bad request")
		return
	}

	var reactionExist int
	var userid int
	stm, err := database.Database.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}
	err = stm.QueryRow(Token.Value).Scan(&userid)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}

	// Start a transaction
	tx, err := database.Database.Begin()
	if err != nil {
		tx.Rollback()
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}

	// Check if the user has already reacted
	stm, err = tx.Prepare("SELECT reaction FROM post_reaction WHERE user_id = ? AND post_id = ?")
	if err != nil {
		tx.Rollback()
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}
	err = stm.QueryRow(userid, postid).Scan(&reactionExist)
	if err == sql.ErrNoRows {

		// No existing like, insert new like
		stm, err = tx.Prepare("INSERT INTO post_reaction (user_id, post_id, reaction) VALUES (?, ?, ?)")
		if err != nil {
			tx.Rollback()
			logger.LogWithDetails(err)
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "500 internal server error")
			return
		}
		_, err = stm.Exec(userid, postid, reaction)
		if err != nil {
			logger.LogWithDetails(err)
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "500 internal server error")
			return
		}
		if reaction == 1 {
			stm, err = tx.Prepare("UPDATE posts SET total_likes =   total_likes  + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		} else {
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes = total_dislikes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		}

	} else if err != nil {
		logger.LogWithDetails(err)
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error3527")
		return
	} else {
		if reactionExist == Neutre && reaction == ReactionLike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(ReactionLike, userid, postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}

		} else if reactionExist == Neutre && reaction == ReactionDislike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(ReactionDislike, userid, postid)
			if err != nil {
				logger.LogWithDetails(err)
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes  = total_dislikes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		} else if reactionExist == ReactionLike && reaction == ReactionLike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(Neutre, userid, postid)
			if err != nil {
				logger.LogWithDetails(err)
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		} else if reactionExist == ReactionLike && reaction == ReactionDislike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(ReactionDislike, userid, postid)
			if err != nil {
				logger.LogWithDetails(err)
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes  = total_dislikes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}

			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		} else if reactionExist == ReactionDislike && reaction == ReactionDislike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(Neutre, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error25828588589")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes  = total_dislikes - 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		} else if reactionExist == ReactionDislike && reaction == ReactionLike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(ReactionLike, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error258285885825")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes  = total_dislikes - 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes +1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "500 internal server error")
				return
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error1")
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
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
