package handlers

import (
	"fmt"
	"net/http"

	"forum/internal/database"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	r.ParseForm()
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
	// lets extract the post id
	post_id := r.URL.Query().Get("id")

	result, err := database.Database.Exec("UPDATE posts SET total_likes = total_likes + 1 WHERE id = $1", post_id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("could not retrieve rows affected: %v", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Printf("no post found with id %v", post_id)
		return
	}
	fmt.Printf("Successfully updated total likes for post ID %v\n", post_id)
	http.Redirect(w, r, "/", http.StatusFound)
}

// todo : complete Dislike handelers
func DislikePost(w http.ResponseWriter, r *http.Request) {
}

// todo : complete the FilterPosts handeler
func FilterPosts(w http.ResponseWriter, r *http.Request) {
}