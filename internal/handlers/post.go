package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal"
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

func PostReactions(w http.ResponseWriter, r *http.Request) {
	pages := internal.Templates
	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.BadRequest, http.StatusBadRequest)
		return
	}

	postid, errpost := strconv.Atoi(r.FormValue("post_id"))
	reaction, errreaction := strconv.Atoi(r.FormValue("reaction"))
	Token, terr := r.Cookie("token")
	if terr != nil {
		logger.LogWithDetails(terr)
		utils.RenderTemplate(w, "error.html", models.BadRequest, http.StatusBadRequest)
		return
	}
	if errpost != nil || errreaction != nil {
		logger.LogWithDetails(fmt.Errorf("%s", " Itoi Errors"))
		utils.RenderTemplate(w, "error.html", models.BadRequest, http.StatusBadRequest)
		return
	}

	_, postExists := utils.IsExist("posts", "id", "", strconv.Itoa(postid))
	if !postExists {
		logger.LogWithDetails(fmt.Errorf("%s", "Post id don't exists"))
		utils.RenderTemplate(w, "error.html", models.BadRequest, http.StatusBadRequest)
		return
	}

	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}

	var reactionExist int
	var userid int
	stm, err := db.Prepare("SELECT id FROM users WHERE token = ?")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = stm.QueryRow(Token.Value).Scan(&userid)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}

	// Check if the user has already reacted
	stm, err = tx.Prepare("SELECT reaction FROM post_reaction WHERE user_id = ? AND post_id = ?")
	if err != nil {
		tx.Rollback()
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = stm.QueryRow(userid, postid).Scan(&reactionExist)
	if err == sql.ErrNoRows {

		// No existing like, insert new like
		stm, err = tx.Prepare("INSERT INTO post_reaction (user_id, post_id, reaction) VALUES (?, ?, ?)")
		if err != nil {
			tx.Rollback()
			logger.LogWithDetails(err)
			utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
			return
		}
		_, err = stm.Exec(userid, postid, reaction)
		if err != nil {
			logger.LogWithDetails(err)
			tx.Rollback()
			utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
			return
		}
		if reaction == 1 {
			stm, err = tx.Prepare("UPDATE posts SET total_likes =   total_likes  + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		} else {
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes = total_dislikes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		}

	} else if err != nil {
		logger.LogWithDetails(err)
		tx.Rollback()
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	} else {
		if reactionExist == Neutre && reaction == ReactionLike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(ReactionLike, userid, postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}

		} else if reactionExist == Neutre && reaction == ReactionDislike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(ReactionDislike, userid, postid)
			if err != nil {
				logger.LogWithDetails(err)
				tx.Rollback()
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes  = total_dislikes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		} else if reactionExist == ReactionLike && reaction == ReactionLike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(Neutre, userid, postid)
			if err != nil {
				logger.LogWithDetails(err)
				tx.Rollback()
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		} else if reactionExist == ReactionLike && reaction == ReactionDislike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(ReactionDislike, userid, postid)
			if err != nil {
				logger.LogWithDetails(err)
				tx.Rollback()
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_dislikes  = total_dislikes + 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}

			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		} else if reactionExist == ReactionDislike && reaction == ReactionDislike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
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
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		} else if reactionExist == ReactionDislike && reaction == ReactionLike {
			stm, err = tx.Prepare("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
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
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			stm, err = tx.Prepare("UPDATE posts SET total_likes = total_likes +1 WHERE id = ?")
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
			_, err = stm.Exec(postid)
			if err != nil {
				tx.Rollback()
				logger.LogWithDetails(err)
				utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
				return
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	fmt.Println(r.Referer())
	http.Redirect(w, r, r.Referer(), http.StatusFound)
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
