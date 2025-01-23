package handlers

import (
	"fmt"
	"net/http"

	"forum/internal/database"
)

func Craete_post(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	user_Id := 3
	title := "oumayma"
	content := "post created by ouma to test!!!"
	total_likes := 1
	total_dislikes := 1

	// lets check for emptyness
	if title == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}

	// lets insert this data to our database
	_, err := database.Database.Exec("INSERT INTO posts (user_id,title,content,total_likes,total_dislikes) VALUES ( ?,?,?,?,?)", user_Id, title, content, total_likes, total_dislikes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
