package handlers

import (
	"fmt"
	"net/http"

	"forum/internal/database"
)

// without authentication !!!!
// let's add the authentication

func CreateComment(w http.ResponseWriter, r *http.Request) {
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

