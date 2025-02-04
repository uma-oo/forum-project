package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/database"
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
	err = database.Database.QueryRow("SELECT last_insert_rowid()").Scan(&postId)
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func LikePost(w http.ResponseWriter, r *http.Request) {
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
	userid, erruser := strconv.Atoi(r.FormValue("user_id"))
	postid, errpost := strconv.Atoi(r.FormValue("post_id"))
	reaction, errreaction := strconv.Atoi(r.FormValue("reaction"))
	if erruser != nil || errpost != nil || errreaction != nil {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}
	// check if user has already reacted
	var reactionExist int
	err = database.Database.QueryRow("SELECT reaction FROM likes WHERE user_id = ? AND post_id = ?", userid, postid).Scan(&reactionExist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}

	if reactionExist == 1 {
		// update reaction instead of inserting a new one
		_, err = database.Database.Exec("UPDATE likes SET reaction = ? WHERE user_id = ? AND post_id = ?", -1, userid, postid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "internal server error")
			return
		}
		if _,err = database.Database.Exec("UPDATE posts SET  total_likes =  total_likes - 1 WHERE id = ?", postid);err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate( w, "error.html", "internal server error")
			return
		}
	} else {
		// insert reaction
		_, err = database.Database.Exec("UPDATE   likes SET (user_id, post_id, reaction) VALUES (?, ?, ?)", userid, postid, reaction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "internal server error")
			return
		}
		if _,err = database.Database.Exec("UPDATE posts SET  total_likes =  total_likes + 1 WHERE id = ?", postid);err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate( w, "error.html", "internal server error")
			return
		}
	}
	http.Redirect(w , r, "/", http.StatusFound)
}