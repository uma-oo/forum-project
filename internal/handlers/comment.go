package handlers

import (
	"fmt"
	"net/http"

	"forum/internal/database"
)


func AddPostComment(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	
	user_id := 1
	post_id := 1
	comment := "hi there this is someone else"

	if comment == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}

	// lets insert this data to our database
	_, err := database.Database.Exec("INSERT INTO comments (user_id, post_id,content) VALUES (?,?,?)", user_id, post_id, comment)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	// nwita jdiiida
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	// lets extract the post id
	post_id := r.URL.Query().Get("id")

	result, err := database.Database.Exec("UPDATE comments SET total_likes = total_likes + 1 WHERE post_id = $1", post_id)
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

	fmt.Printf("Successfully updated totallikes for post ID %v\n", post_id)
	http.Redirect(w, r, "/", http.StatusFound)
}

// todo : complete dislikeComment Handler.
func DislikeComment(w http.ResponseWriter, r *http.Request) {
}