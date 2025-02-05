package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/database"
)

var (
	like    bool
	dislike bool
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	r.ParseForm() ///////////////////////////
	categories := r.Form["post-categorie"]
	postContent := r.FormValue("postBody")
	postTitle := r.FormValue("postTitle")
	// lets check for emptyness
	if postContent == "" || postTitle == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}
	// get the user ID from the session
	cookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		pages.ExecuteTemplate(w, "error.html", "unauthorized")
		return
	}
	// get the user ID from the users table
	var userId int
	err = database.Database.QueryRow("SELECT id FROM users WHERE token = ?", cookie.Value).Scan(&userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}

	// lets insert this data to our database
	_, err = database.Database.Exec("INSERT INTO posts (user_id,title,content) VALUES ( ?,?,?)", userId, postTitle, postContent)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	// get the last inserted post id
	var postId int
	err = database.Database.QueryRow("SELECT last_insert_rowid()").Scan(&postId) /////////////////////////// jjbcjkjnkjnkjnkjnkjnkj
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	// insert categories
	for _, category := range categories {
		_, err = database.Database.Exec("INSERT INTO categories (category, post_id) VALUES (?, ?)", category, postId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "internal server error")
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

const (
	ReactionLike    = 1
	ReactionDislike = -1
	Neutre          = 0
)

func PostReactions(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "invalid request")
		return
	}

	postid, errpost := strconv.Atoi(r.FormValue("post_id"))
	reaction, errreaction := strconv.Atoi(r.FormValue("reaction"))
	Token, terr := r.Cookie("token")
	if terr != nil {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}

	if errpost != nil || errreaction != nil {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}

	var reactionExist int
	var userid int
	err = database.Database.QueryRow("SELECT id FROM users WHERE token = ?", Token.Value).Scan(&userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}

	// Start a transaction
	tx, err := database.Database.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}

	// Check if the user has already reacted
	err = tx.QueryRow("SELECT reaction FROM post_reaction WHERE user_id = ? AND post_id = ?", userid, postid).Scan(&reactionExist)
	if err == sql.ErrNoRows {
		fmt.Println("haarani")
		// No existing like, insert new like
		_, err = tx.Exec("INSERT INTO post_reaction (user_id, post_id, reaction) VALUES (?, ?, ?)", userid, postid, reaction)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "internal server error4")
			return
		}
		if reaction == 1 {
			fmt.Println("waldi")
			_, err = tx.Exec("UPDATE posts SET total_likes =   total_likes  + 1 WHERE id = ?", postid)
		} else {
			_, err = tx.Exec("UPDATE posts SET total_dislikes = total_dislikes + 1 WHERE id = ?", postid)
		}

	} else if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error3527")
		return
	} else {
		if reactionExist == Neutre && reaction == ReactionLike {
			_, err = tx.Exec("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?", ReactionLike, userid, postid)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error2528552")
				return
			}
			_, err = tx.Exec("UPDATE posts SET total_likes = total_likes + 1 WHERE id = ?", postid)

		} else if reactionExist == Neutre && reaction == ReactionDislike {
			_, err = tx.Exec("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?", ReactionDislike, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error966")
				return
			}

			_, err = tx.Exec("UPDATE posts SET total_dislikes  = total_dislikes + 1 WHERE id = ?", postid)
		} else if reactionExist == ReactionLike && reaction == ReactionLike {
			_, err = tx.Exec("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?", Neutre, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error58/28/")
				return
			}
			_, err = tx.Exec("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?", postid)

		} else if reactionExist == ReactionLike && reaction == ReactionDislike {
			_, err = tx.Exec("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?", ReactionDislike, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error2582858858")
				return
			}

			_, err = tx.Exec("UPDATE posts SET total_dislikes  = total_dislikes + 1 WHERE id = ?", postid)

			_, err = tx.Exec("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?", postid)
		} else if reactionExist == ReactionDislike && reaction == ReactionDislike {
			_, err = tx.Exec("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?", Neutre, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error25828588589")
				return
			}

			_, err = tx.Exec("UPDATE posts SET total_dislikes  = total_dislikes - 1 WHERE id = ?", postid)
		} else if reactionExist == ReactionDislike && reaction == ReactionLike {
			_, err = tx.Exec("UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?", ReactionLike, userid, postid)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				pages.ExecuteTemplate(w, "error.html", "internal server error258285885825")
				return
			}

			_, err = tx.Exec("UPDATE posts SET total_dislikes  = total_dislikes - 1 WHERE id = ?", postid)
			_, err = tx.Exec("UPDATE posts SET total_likes = total_likes +1 WHERE id = ?", postid)
		}

		// if err != nil {
		// 	tx.Rollback()
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	pages.ExecuteTemplate(w, "error.html", "internal server error2")
		// 	return
		// }
	}

	err = tx.Commit()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error1")
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
